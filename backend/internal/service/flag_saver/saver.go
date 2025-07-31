package flag_saver

import (
	"JacFARM/internal/models"
	"JacFARM/internal/rabbitmq"
	"fmt"
	"log/slog"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/jacute/prettylogger"
)

func (fs *FlagSaver) Start() error {
	const op = "service.flag_sender.Start"
	log := fs.log.With(slog.String("op", op))
	log.Info("Starting FlagSender service")

	flagChan, err := fs.queue.GetFlagChan()
	if err != nil {
		log.Error("Failed to get flag channel", prettylogger.Err(err))
		return err
	}

	var processFlagWg sync.WaitGroup
	for {
		select {
		case flag, ok := <-flagChan:
			if !ok {
				processFlagWg.Wait()
				log.Info("Flag channel closed")
				return nil
			}

			log.Info("Received flag", "body", string(flag.Body))

			processFlagWg.Add(1)
			go func() {
				defer processFlagWg.Done()
				if err := fs.processFlag(flag.Body); err != nil {
					log.Error("Failed to process flag", "error", err)
					flag.Nack(false, true) // if error, requeue the message
					return
				}

				flag.Ack(false) // if success, send ack
			}()
		case <-fs.stopChan:
			processFlagWg.Wait()
			return nil
		}
	}
}

func (fs *FlagSaver) processFlag(flagBytes []byte) error {
	const op = "service.flag_sender.processFlag"

	var flag *rabbitmq.Flag
	if err := sonic.Unmarshal(flagBytes, &flag); err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	err := fs.db.PutFlag(&models.Flag{
		Value:             flag.Value,
		Status:            models.FlagStatusPending,
		Exploit:           flag.Exploit,
		GetFrom:           flag.GetFrom,
		MessageFromServer: "",
	})
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (fs *FlagSaver) Stop() {
	const op = "service.flag_sender.Stop"

	close(fs.stopChan)
	fs.log.Info("FlagSender stopped gracefully", slog.String("op", op))
}
