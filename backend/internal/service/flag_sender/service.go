package flag_sender

import (
	"JacFARM/internal/models"
	"JacFARM/pkg/plugin_interfaces"
	"context"
	"fmt"
	"log/slog"
	"path"
	"plugin"
	"time"

	"github.com/jacute/prettylogger"
)

type storage interface {
	GetConfigParameter(ctx context.Context, key string) (string, error)
	GetFlagValuesByStatus(ctx context.Context, status models.FlagStatus) ([]string, error)
	UpdateStatusForOldFlags(ctx context.Context, flagTTL time.Duration) (int64, error)
	UpdateFlagStatus(ctx context.Context, flag string, status models.FlagStatus) error
}

type FlagSender struct {
	log          *slog.Logger
	db           storage
	cfg          *config
	shutdownChan chan struct{}
}

func New(log *slog.Logger, db storage, pluginDir string) *FlagSender {
	er := &FlagSender{
		log: log,
		db:  db,
		cfg: &config{
			pluginDir: pluginDir,
		},
		shutdownChan: make(chan struct{}),
	}
	err := er.loadConfig(context.Background(), true)
	if err != nil {
		panic("error loading exploit runner config: " + err.Error())
	}

	return er
}

func (fs *FlagSender) Start() error {
	const op = "service.flag_sender.Start"
	log := fs.log.With(slog.String("op", op))
	log.Info("starting flag sender service")

	// load plugin
	pluginPath := path.Join(fs.cfg.pluginDir, fs.cfg.plugin+".so")
	sendPlugin, err := plugin.Open(pluginPath)
	if err != nil {
		log.Error(
			"error opening send plugin",
			slog.String("plugin_path", pluginPath),
			prettylogger.Err(err),
		)
		return err
	}
	symbol, err := sendPlugin.Lookup("NewClient")
	if err != nil {
		log.Error(
			"error looking up send plugin client",
			prettylogger.Err(err),
		)
		return err
	}
	clientInit, ok := symbol.(*plugin_interfaces.NewClientFunc)
	if !ok {
		log.Error("plugin client constructor not func(url, token string) Client")
		return fmt.Errorf("plugin client constructor not func(url, token string) Client")
	}
	client := (*clientInit)(fs.cfg.juryFlagURL, fs.cfg.token)

	ticker := time.NewTicker(fs.cfg.submitPeriod)
	for {
		select {
		case <-ticker.C:
			sendCtx, cancel := context.WithTimeout(context.Background(), fs.cfg.submitTimeout)

			// every tick load new config from db
			err := fs.loadConfig(sendCtx, false)
			if err != nil {
				log.Error(
					"error reloading flag sender config from db",
					prettylogger.Err(err),
				)
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

			for flag, status := range result {
				err := fs.db.UpdateFlagStatus(context.Background(), flag, status)
				if err != nil {
					log.Error("error updating flag status", slog.String("flag", flag), prettylogger.Err(err))
				}
			}
			ticker.Reset(fs.cfg.submitPeriod)
		case <-fs.shutdownChan:
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
