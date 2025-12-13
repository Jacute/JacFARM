package worker

import (
	"log/slog"
	"time"
)

const (
	defaultAttackPeriod          = 5 * time.Second
	defaultMaxConcurrentExploits = 5
)

type JacFARMClient interface {
}

type Worker struct {
	client                JacFARMClient
	attackPeriod          time.Duration
	maxConcurrentExploits int

	log    *slog.Logger
	stopCh chan struct{}
}

type options struct {
	attackPeriod          *time.Duration
	maxConcurrentExploits *int
}

type Option func(opts *options) error

func New(client JacFARMClient, log *slog.Logger, opts ...Option) (*Worker, error) {
	workerOpts := &options{}
	for _, opt := range opts {
		err := opt(workerOpts)
		if err != nil {
			return nil, err
		}
	}

	w := &Worker{
		client:                client,
		attackPeriod:          defaultAttackPeriod,
		maxConcurrentExploits: defaultMaxConcurrentExploits,

		log:    log,
		stopCh: make(chan struct{}),
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

func (w *Worker) Run() {
	const op = "worker.Run"
	log := w.log.With(slog.String("op", op))

	log.Info(
		"starting worker",
		slog.Int("max_concurrent_exploits", w.maxConcurrentExploits),
		slog.Duration("attack_period", w.attackPeriod),
	)
	for {
		select {
		case <-w.stopCh:
			return
		default:

		}
	}
}

func (w *Worker) Stop() {
	const op = "worker.Stop"
	w.log.Info("stopping worker")
	w.stopCh <- struct{}{}
}
