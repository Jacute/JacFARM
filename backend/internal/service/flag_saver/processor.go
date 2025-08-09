package flag_saver

import (
	"JacFARM/internal/models"
	"JacFARM/internal/rabbitmq"
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
)

func (fs *FlagSaver) processFlag(flagBytes []byte) error {
	const op = "service.flag_saver.processFlag"

	var flag *rabbitmq.Flag
	if err := sonic.Unmarshal(flagBytes, &flag); err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	_, err := fs.db.PutFlag(context.Background(), &models.Flag{
		Value:             flag.Value,
		Status:            models.FlagStatusPending,
		ExploitID:         flag.ExploitID,
		GetFrom:           flag.TeamID,
		MessageFromServer: "",
		CreatedAt:         time.Now().UTC().Unix(),
	})
	if err != nil {
		return err
	}

	return nil
}
