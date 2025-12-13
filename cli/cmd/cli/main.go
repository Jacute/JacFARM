package main

import (
	"cli_exploit_runner/internal/cli"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jacute/prettylogger"
)

func main() {
	args, err := cli.ParseArgs()
	if err != nil {
		fmt.Printf("usage: %s <jacfarm_host> <exploit>\n-help for more information\n", os.Args[0])
		os.Exit(1)
	}

	log := slog.New(prettylogger.NewColoredHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Info("running script")
	err = cli.Run(args, log)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	sign := <-sigCh

	log.Info("stopping script", slog.String("signal", sign.String()))
}
