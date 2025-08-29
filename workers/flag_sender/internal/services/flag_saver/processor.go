package flag_saver

import (
	"context"
	"errors"
	"flag_sender/internal/models"
	"flag_sender/internal/postgres"
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
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Debug("got flag", slog.Any("flag", flag))

	dbFlag := &models.Flag{
		Value:             flag.Value,
		Status:            models.FlagStatusPending,
		ExploitID:         &flag.ExploitID,
		GetFrom:           &flag.TeamID,
		MessageFromServer: "",
		CreatedAt:         time.Now().UTC(),
	}
	if flag.SourceType == rabbitmq_dto.ManualSendingSourceType {
		dbFlag.ExploitID = nil
		dbFlag.GetFrom = nil
	}

	flagID, err := fs.db.PutFlag(context.Background(), dbFlag)
	if err != nil {
		if errors.Is(err, postgres.ErrFlagAlreadyExists) {
			log.Info("flag already exists, skipping", slog.Any("flag", flag))
			return nil // skipping exists flag
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("flag send successfully", slog.Int64("flag_id", flagID))

	return nil
}
