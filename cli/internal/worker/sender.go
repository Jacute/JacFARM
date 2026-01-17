package worker

import (
	jacfarm_client "cli_exploit_runner/internal/clients/jacfarm"
	"context"
	"log/slog"

	"github.com/jacute/prettylogger"
)

const senderSize = 50

func (w *Worker) runSender() {
	const op = "worker.startSender"
	log := w.log.With(slog.String("op", op))

	flagBuffer := make([]*jacfarm_client.ServiceFlag, 0, senderSize)

	for {
		select {
		case <-w.stopCh:
			for {
				select {
				case flags := <-w.flagQueue:
					flagBuffer = append(flagBuffer, flags...)
				default:
					err := w.client.SendFlags(context.Background(), flagBuffer)
					if err != nil {
						log.Error("error sending flags", prettylogger.Err(err))
					}
					return
				}
			}
		case flags := <-w.flagQueue:
			flagBuffer = append(flagBuffer, flags...)
			if len(flagBuffer) >= senderSize {
				err := w.client.SendFlags(context.Background(), flagBuffer)
				if err != nil {
					log.Error("error sending flags", prettylogger.Err(err))
				}
				flagBuffer = flagBuffer[:0]
			}
		}
	}
}
