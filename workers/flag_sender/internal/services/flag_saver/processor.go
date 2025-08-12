package flag_saver

import (
	"context"
	"flag_sender/internal/models"
	"flag_sender/internal/rabbitmq"
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
		CreatedAt:         time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	return nil
}
