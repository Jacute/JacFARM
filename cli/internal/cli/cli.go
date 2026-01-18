package cli

import (
	"errors"
	"flag"
	"fmt"
	jacfarm_client "jacfarmcli/internal/clients/jacfarm"
	"jacfarmcli/internal/worker"
	"log/slog"
	"os"
	"time"

	"github.com/jacute/prettylogger"
)

var (
	ErrUsage = errors.New("usage error")
)

type Args struct {
	Addr                  string
	FlagRe                string
	Token                 string
	Port                  int
	ExecutablePath        string
	Timeout               int
	AttackPeriod          int
	MaxConcurrentExploits int
}

func ParseArgs() (*Args, error) {
	var args Args

	flag.IntVar(&args.Timeout, "t", 5, "timeout for http client (in seconds)")
	flag.StringVar(&args.Token, "a", "", "JacFARM auth token")
	flag.IntVar(&args.AttackPeriod, "period", 5, "attack period (in seconds)")
	flag.IntVar(&args.MaxConcurrentExploits, "c", 50, "max concurrent exploits in one time")
	flag.IntVar(&args.Port, "p", 15050, "jacfarm port")
	flag.StringVar(&args.FlagRe, "f", "[A-Z0-9]{31}=", "flag regex")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s [flags] <addr> <executable>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example:\n")
		fmt.Fprintf(os.Stderr, "  %s -a TOKEN -p 15050 localhost ./exploit\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if args.Token == "" {
		return nil, flag.ErrHelp
	}

	rest := flag.Args()
	if len(rest) < 1 {
		return nil, flag.ErrHelp
	}

	args.Addr = rest[0]
	args.ExecutablePath = rest[1]

	return &args, nil
}

func ValidateArgs(args *Args) error {
	if _, err := os.Stat(args.ExecutablePath); os.IsNotExist(err) {
		return err
	}
	return nil
}

// Run starts the worker
// Non-block function
func Run(args *Args, log *slog.Logger) error {
	const op = "cli.Run"

	client, err := jacfarm_client.New(
		args.Addr,
		args.Token,
		jacfarm_client.WithCustomPort(args.Port),
		jacfarm_client.WithTimeout(time.Duration(args.Timeout)*time.Second),
	)
	if err != nil {
		log.Error("error creating jacfarm client", prettylogger.Err(err))
		return fmt.Errorf("%s: error creating jacfarm client %w", op, err)
	}
	w, err := worker.New(
		client, log,
		args.ExecutablePath,
		args.FlagRe,
		worker.WithAttackPeriod(time.Duration(args.AttackPeriod)*time.Second),
		worker.WithMaxConcurrentExploits(args.MaxConcurrentExploits),
	)
	if err != nil {
		log.Error("error creating worker", prettylogger.Err(err))
		return fmt.Errorf("%s: error creating worker %e", op, err)
	}

	w.Run()

	return nil
}
