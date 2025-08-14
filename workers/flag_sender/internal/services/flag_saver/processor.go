package flag_saver

import (
	"context"
	"flag_sender/internal/models"
	"flag_sender/pkg/rabbitmq_dto"
	"fmt"
	"log/slog"
	"time"

	"github.com/bytedance/sonic"
)

func (fs *FlagSaver) processFlag(flagBytes []byte) error {
	const op = "service.flag_saver.processFlag"
	log := fs.log.With(slog.String("op", op))

	var flag *rabbitmq_dto.Flag
	if err := sonic.Unmarshal(flagBytes, &flag); err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	log.Debug("got flag", slog.Any("flag", flag))

	_, err := fs.db.PutFlag(context.Background(), &models.Flag{
		Value:             flag.Value,
		Status:            models.FlagStatusPending,
		ExploitID:         flag.ExploitID,
		GetFrom:           flag.TeamID,
		MessageFromServer: "",
		CreatedAt:         time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	return nil
}
