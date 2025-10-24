package flag_sender

import (
	"context"
	"flag_sender/internal/models"
	"flag_sender/pkg/rabbitmq_dto"
	"fmt"
	"log/slog"
	"time"

	"github.com/bytedance/sonic"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (fs *FlagSender) processBatch(ctx context.Context, batch []amqp.Delivery) error {
	const op = "service.flag_sender.processBatch"
	log := fs.log.With(slog.String("op", op))

	flags := make([]*models.Flag, 0, len(batch))
	expiredFlags := make([]*models.Flag, 0, len(batch))
	for _, batchPart := range batch {
		var flag *rabbitmq_dto.Flag
		if err := sonic.Unmarshal(batchPart.Body, &flag); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		mappedFlag := mapQueueFlagIntoModel(flag)
		if time.Since(flag.CreatedAt) > fs.cfg.flagTTL {
			log.Debug("flag expired, skipping", slog.Any("flag", flag))
			expiredFlags = append(expiredFlags, mappedFlag)
			continue
		}
		flags = append(flags, mappedFlag)
	}
	if len(expiredFlags) > 0 {
		ids, err := fs.db.PutFlags(ctx, expiredFlags)
		if err != nil {
			return fmt.Errorf("%s: failed to add expired flags: %w", op, err)
		}
		log.Debug("added expired flags", slog.Any("ids", ids))
	}

	flagStrings := make([]string, 0, len(flags))
	for _, flag := range flags {
		flagStrings = append(flagStrings, flag.Value)
	}

	log.Info("sending flags", slog.Int("count", len(flags)))
	result, err := fs.pluginClient.SendFlags(ctx, flagStrings)
	if err != nil {
		return fmt.Errorf("error sending flags: %w", err)
	}
	log.Info("flags send successfully", slog.Int("count", len(result)))
	for flag, flagResult := range result {
		flagModel := findFlagByValue(flags, flag)
		if flagModel == nil {
			continue
		}
		flagModel.Status = flagResult.Status
		flagModel.MessageFromServer = flagResult.Msg
	}

	if len(flags) > 0 {
		ids, err := fs.db.PutFlags(context.Background(), flags)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		log.Debug("flags added in db", slog.Int("count", len(ids)))
	}

	return nil
}

func mapQueueFlagIntoModel(flag *rabbitmq_dto.Flag) *models.Flag {
	dbFlag := &models.Flag{
		Value:             flag.Value,
		Status:            models.FlagStatusOld,
		ExploitID:         &flag.ExploitID,
		GetFrom:           &flag.TeamID,
		MessageFromServer: "",
		CreatedAt:         time.Now().UTC(),
	}
	if flag.SourceType == rabbitmq_dto.ManualSendingSourceType {
		dbFlag.ExploitID = nil
		dbFlag.GetFrom = nil
	}
	return dbFlag
}

func findFlagByValue(flags []*models.Flag, value string) *models.Flag {
	for _, flag := range flags {
		if flag.Value == value {
			return flag
		}
	}
	return nil
}
