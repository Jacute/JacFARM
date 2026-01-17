package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/jacute/prettylogger"
)

func (w *Worker) runExecutor() {
	const op = "worker.startExecutor"
	log := w.log.With(slog.String("op", op))

	timer := time.NewTimer(w.attackPeriod)

	for {
		select {
		case <-w.stopCh:
			return
		case <-timer.C:
			log.Info("starting attack")
			ctx, cancel := context.WithTimeout(context.Background(), w.attackPeriod)
			defer cancel()

			teams, err := w.client.GetTeams(ctx)
			if err != nil {
				log.Error(
					"error getting teams",
					prettylogger.Err(err),
				)
				timer.Reset(w.attackPeriod)
				continue
			}

			err = w.attackAll(ctx, teams)
			if err != nil {
				log.Error(
					"error attacking",
					prettylogger.Err(err),
				)
				timer.Reset(w.attackPeriod)
				continue
			}
			timer.Reset(w.attackPeriod)
		}
	}
}
