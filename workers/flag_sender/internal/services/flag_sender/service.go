package flag_sender

import (
	"context"
	"flag_sender/internal/models"
	"flag_sender/pkg/plugins"
	"fmt"
	"log/slog"
	"path"
	"plugin"
	"time"

	"github.com/jacute/prettylogger"
)

//go:generate mockgen -source=service.go -destination=./mocks/storage_mock.go -package=mocks -mock_names=storage=StorageMock storage
type storage interface {
	GetConfigParameter(ctx context.Context, key string) (string, error)
	GetFlagValuesByStatus(ctx context.Context, status models.FlagStatus) ([]string, error)
	UpdateStatusForOldFlags(ctx context.Context, flagTTL time.Duration) (int64, error)
	UpdateFlagByResult(ctx context.Context, flag string, result *plugins.FlagResult) error
}

type FlagSender struct {
	log          *slog.Logger
	db           storage
	cfg          *config
	shutdownChan chan struct{}
}

func New(log *slog.Logger, db storage, pluginDir string) (*FlagSender, error) {
	fs := &FlagSender{
		log: log,
		db:  db,
		cfg: &config{
			pluginDir: pluginDir,
		},
		shutdownChan: make(chan struct{}),
	}
	err := fs.loadConfig(context.Background(), true)
	if err != nil {
		return nil, fmt.Errorf("error loading exploit runner config: %s", err.Error())
	}
	return fs, nil
}

func (fs *FlagSender) Start() error {
	const op = "service.flag_sender.Start"
	log := fs.log.With(slog.String("op", op))
	log.Info("starting flag sender service")

	// load plugin
	pluginPath := path.Join(fs.cfg.pluginDir, fs.cfg.plugin+".so")
	sendPlugin, err := plugin.Open(pluginPath)
	if err != nil {
		log.Error("error opening send plugin",
			slog.String("plugin_path", pluginPath),
			prettylogger.Err(err),
		)
		return err
	}
	symbol, err := sendPlugin.Lookup("NewClient")
	if err != nil {
		log.Error("error looking up send plugin client", prettylogger.Err(err))
		return err
	}
	clientInit, ok := symbol.(*plugins.NewClientFunc)
	if !ok {
		log.Error("plugin client constructor not func(url, token string) Client")
		return fmt.Errorf("plugin client constructor not func(url, token string) Client")
	}
	client := (*clientInit)(fs.cfg.juryFlagURL, fs.cfg.token)

	for {
		timer := time.NewTimer(fs.cfg.submitPeriod)
		select {
		case <-timer.C:
			sendCtx, cancel := context.WithTimeout(context.Background(), fs.cfg.submitTimeout)

			// reload config
			err := fs.loadConfig(sendCtx, false)
			if err != nil {
				log.Error("error reloading flag sender config from db", prettylogger.Err(err))
				cancel()
				continue
			}

			flagsOld, err := fs.db.UpdateStatusForOldFlags(sendCtx, fs.cfg.flagTTL)
			if err != nil {
				log.Error("error setting flags to old", prettylogger.Err(err))
			} else {
				log.Info("flags set to old", slog.Int64("count", flagsOld))
			}

			flags, err := fs.db.GetFlagValuesByStatus(sendCtx, models.FlagStatusPending)
			if err != nil {
				log.Error("error getting flags from db", prettylogger.Err(err))
				cancel()
				continue
			}

			if len(flags) == 0 {
				cancel()
				continue
			}

			log.Info("sending flags", slog.Int("count", len(flags)))
			result, err := client.SendFlags(flags)
			if err != nil {
				log.Error("error sending flags", prettylogger.Err(err))
				cancel()
				continue
			}
			cancel()

			for flag, result := range result {
				err := fs.db.UpdateFlagByResult(context.Background(), flag, result)
				if err != nil {
					log.Error("error updating flag status", slog.String("flag", flag), prettylogger.Err(err))
				}
			}

		case <-fs.shutdownChan:
			timer.Stop()
			log.Info("flag sender service shut down")
			return nil
		}
	}
}

func (fs *FlagSender) Stop() {
	const op = "service.flag_sender.Stop"
	log := fs.log.With(slog.String("op", op))

	log.Debug("stopping flag sender service")
	close(fs.shutdownChan)
}
