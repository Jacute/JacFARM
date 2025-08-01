package flag_saver

import (
	"JacFARM/internal/models"
	"JacFARM/internal/rabbitmq"
	"fmt"

	"github.com/bytedance/sonic"
)

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
