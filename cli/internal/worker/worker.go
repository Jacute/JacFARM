package worker

import (
	"context"
	"errors"
	"fmt"
	jacfarm_client "jacfarmcli/internal/clients/jacfarm"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"sync"
	"time"

	"github.com/jacute/prettylogger"
)

var (
	ErrExploitNotExecutable = errors.New("exploit is not executable")
	ErrExploitNotFile       = errors.New("exploit is not file")
)

const (
	defaultAttackPeriod          = 5 * time.Second
	defaultMaxConcurrentExploits = 5
)

//go:generate mockgen -source=worker.go -destination=./mocks/worker_mock.go -package=mocks -mock_names=storage=WorkerMock Service
type JacFARMClient interface {
	GetTeams(ctx context.Context) ([]*jacfarm_client.Team, error)
	SendFlags(ctx context.Context, flags []*jacfarm_client.ServiceFlag) error
}

type Worker struct {
	client                JacFARMClient
	attackPeriod          time.Duration
	maxConcurrentExploits int
	exploitPath           string
	flagRe                *regexp.Regexp

	log       *slog.Logger
	stopCh    chan struct{}
	flagQueue chan []*jacfarm_client.ServiceFlag
}

type options struct {
	attackPeriod          *time.Duration
	maxConcurrentExploits *int
}

type Option func(opts *options) error

func New(
	client JacFARMClient,
	log *slog.Logger,
	exploitPath string,
	flagRegexp string,
	opts ...Option,
) (*Worker, error) {
	workerOpts := &options{}
	for _, opt := range opts {
		err := opt(workerOpts)
		if err != nil {
			return nil, err
		}
	}

	flagRe, err := regexp.Compile(flagRegexp)
	if err != nil {
		return nil, err
	}

	w := &Worker{
		client:                client,
		attackPeriod:          defaultAttackPeriod,
		maxConcurrentExploits: defaultMaxConcurrentExploits,
		exploitPath:           exploitPath,
		flagRe:                flagRe,

		log:       log,
		stopCh:    make(chan struct{}),
		flagQueue: make(chan []*jacfarm_client.ServiceFlag, senderSize), // TODO: optimization of buf size
	}

	if workerOpts.attackPeriod != nil {
		w.attackPeriod = *workerOpts.attackPeriod
	}
	if workerOpts.maxConcurrentExploits != nil {
		w.maxConcurrentExploits = *workerOpts.maxConcurrentExploits
	}

	return w, nil
}

func WithAttackPeriod(attackPeriod time.Duration) Option {
	return func(opts *options) error {
		opts.attackPeriod = &attackPeriod
		return nil
	}
}

func WithMaxConcurrentExploits(count int) Option {
	return func(opts *options) error {
		opts.maxConcurrentExploits = &count
		return nil
	}
}

// Run starts the worker. It is a non-blocking function.
// It starts the receiver (sender) and the consumer (executor) in separate goroutines.
// The receiver (sender) is responsible for receiving flags from the JacFARM client and sending them to the consumer (executor).
// The consumer (executor) is responsible for executing the exploits and sending the result back to the JacFARM client.
// The function logs information about the worker when it starts, including the maximum number of concurrent exploits and the attack period.
func (w *Worker) Run() {
	const op = "worker.Run"
	log := w.log.With(slog.String("op", op))
	log.Info(
		"starting worker",
		slog.Int("max_concurrent_exploits", w.maxConcurrentExploits),
		slog.Duration("attack_period", w.attackPeriod),
	)

	// start receiver (sender)
	go w.runSender()
	// start consumer (executor)
	go w.runExecutor()
}

func (w *Worker) attackAll(
	ctx context.Context,
	teams []*jacfarm_client.Team,
) error {
	const op = "worker.attackAll"
	log := w.log.With(slog.String("op", op), slog.String("exploit_path", w.exploitPath))

	info, err := os.Stat(w.exploitPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("exploit %s does not exist", w.exploitPath)
		}
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("exploit %s is not file", w.exploitPath)
	}
	perms := info.Mode().Perm()
	isExec := perms&0111 != 0
	if !isExec {
		return fmt.Errorf("exploit %s is not executable", w.exploitPath)
	}

	concurrentCh := make(chan struct{}, w.maxConcurrentExploits)
	wg := &sync.WaitGroup{}
	wg.Add(len(teams))

	for _, t := range teams {
		concurrentCh <- struct{}{}
		go func() {
			defer func() {
				wg.Done()
				<-concurrentCh
			}()
			out, err := attack(ctx, w.exploitPath, t.IP.String())
			if err != nil {
				log.Error(
					"error attacking team",
					prettylogger.Err(err),
					slog.String("team_ip", t.IP.String()),
				)
				return
			}
			flags := parseFlags(out, w.flagRe)
			for _, f := range flags {
				w.flagQueue <- []*jacfarm_client.ServiceFlag{
					{
						Flag:   f,
						TeamID: t.ID,
					},
				}
			}
		}()
	}

	wg.Wait()
	close(concurrentCh)
	return nil
}

func attack(ctx context.Context, exploitPath, targetIP string) (exploitOut []byte, err error) {
	cmd := exec.CommandContext(ctx, exploitPath, targetIP)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return out, err
}

func parseFlags(text []byte, re *regexp.Regexp) []string {
	flags := make([]string, 0)
	bts := re.FindAll(text, -1)
	for _, b := range bts {
		flags = append(flags, string(b))
	}
	return flags
}

func (w *Worker) Stop() {
	const op = "worker.Stop"
	w.log.Info("stopping worker")
	close(w.stopCh)
}
