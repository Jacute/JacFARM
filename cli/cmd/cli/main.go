package main

import (
	"flag"
	"fmt"
	"jacfarmcli/internal/cli"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jacute/prettylogger"
)

func main() {
	args, err := cli.ParseArgs()
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}
	err = cli.ValidateArgs(args)
	if err != nil {
		fmt.Println(err.Error())
		flag.Usage()
		os.Exit(2)
	}

	log := slog.New(prettylogger.NewColoredHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Info("running script")
	err = cli.Run(args, log)
	if err != nil {
		os.Exit(2)
	}

	sign := <-sigCh

	log.Info("stopping script", slog.String("signal", sign.String()))
}
