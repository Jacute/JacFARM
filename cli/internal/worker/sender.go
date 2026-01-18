package worker

import (
	jacfarm_client "cli_exploit_runner/internal/clients/jacfarm"
	"context"
	"log/slog"
	"time"

	"github.com/jacute/prettylogger"
)

const senderSize = 50
const sendPeriod = 5 * time.Second

func (w *Worker) runSender() {
	const op = "worker.startSender"
	log := w.log.With(slog.String("op", op))

	timer := time.NewTimer(sendPeriod)
	defer timer.Stop()
	flagBuffer := make([]*jacfarm_client.ServiceFlag, 0, senderSize)

	for {
		select {
		case <-w.stopCh:
			for {
				select {
				case flags := <-w.flagQueue:
					flagBuffer = append(flagBuffer, flags...)
				default:
					if len(flagBuffer) > 0 {
						err := w.client.SendFlags(context.Background(), flagBuffer)
						if err != nil {
							log.Error("error sending flags", prettylogger.Err(err))
						} else {
							log.Info("flags send successfully", slog.Int("flags_count", len(flagBuffer)), slog.Any("first_flags", flagBuffer[:5]))
						}
					}
					return
				}
			}
		case <-timer.C:
			if len(flagBuffer) > 0 {
				err := w.client.SendFlags(context.Background(), flagBuffer)
				if err != nil {
					log.Error("error sending flags", prettylogger.Err(err))
				} else {
					log.Info("flags send successfully", slog.Int("flags_count", len(flagBuffer)), slog.Any("first_flags", flagBuffer[:5]))
				}
				flagBuffer = flagBuffer[:0]
			}
			timer.Reset(sendPeriod)
		case flags := <-w.flagQueue:
			flagBuffer = append(flagBuffer, flags...)
			if len(flagBuffer) >= senderSize {
				err := w.client.SendFlags(context.Background(), flagBuffer)
				if err != nil {
					log.Error("error sending flags", prettylogger.Err(err))
				} else {
					log.Info("flags send successfully", slog.Int("flags_count", len(flagBuffer)), slog.Any("first_flags", flagBuffer[:5]))
				}
				flagBuffer = flagBuffer[:0]
			}
		}
	}
}
