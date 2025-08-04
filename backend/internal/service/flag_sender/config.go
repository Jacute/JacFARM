package flag_sender

import (
	jacfarm "JacFARM/internal/service/JacFARM"
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/jacute/prettylogger"
)

type config struct {
	plugin        string
	pluginDir     string
	juryFlagURL   string
	token         string
	flagTTL       time.Duration
	submitTimeout time.Duration
	submitPeriod  time.Duration
	submitLimit   int
}

func (er *FlagSender) loadConfig(ctx context.Context, updatePlugin bool) error {
	const op = "service.flag_sender.LoadConfig"
	log := er.log.With(slog.String("op", op))

	if updatePlugin {
		plugin, err := er.db.GetConfigParameter(ctx, jacfarm.ConfigFlagSenderPlugin)
		if err != nil {
			log.Error(
				"error getting flag sender plugin from db config",
				slog.String("config_key", jacfarm.ConfigFlagSenderPlugin),
				prettylogger.Err(err),
			)
			return err
		}
		er.cfg.plugin = plugin
	}
	juryFlagURL, err := er.db.GetConfigParameter(ctx, jacfarm.ConfigFlagSenderJuryFlagURL)
	if err != nil {
		log.Error(
			"error getting jury flag url from db config",
			slog.String("config_key", jacfarm.ConfigFlagSenderJuryFlagURL),
			prettylogger.Err(err),
		)
		return err
	}

	token, err := er.db.GetConfigParameter(ctx, jacfarm.ConfigFlagSenderToken)
	if err != nil {
		log.Error(
			"error getting flag sender token from db config",
			slog.String("config_key", jacfarm.ConfigFlagSenderToken),
			prettylogger.Err(err),
		)
		return err
	}

	flagTTLstr, err := er.db.GetConfigParameter(ctx, jacfarm.ConfigFlagSenderFlagTTL)
	if err != nil {
		log.Error(
			"error getting flag ttl from db config",
			slog.String("config_key", jacfarm.ConfigFlagSenderFlagTTL),
			prettylogger.Err(err),
		)
		return err
	}
	flagTTL, err := time.ParseDuration(flagTTLstr)
	if err != nil {
		log.Error(
			"error parsing flag ttl from db config",
			slog.String("config_key", jacfarm.ConfigFlagSenderFlagTTL),
			prettylogger.Err(err),
		)
		return err
	}

	submitPeriodStr, err := er.db.GetConfigParameter(ctx, jacfarm.ConfigFlagSenderSubmitPeriod)
	if err != nil {
		log.Error(
			"error getting submit period from db config",
			slog.String("config_key", jacfarm.ConfigFlagSenderSubmitPeriod),
			prettylogger.Err(err),
		)
		return err
	}
	submitPeriod, err := time.ParseDuration(submitPeriodStr)
	if err != nil {
		log.Error(
			"error parsing submit period from db config",
			slog.String("config_key", jacfarm.ConfigFlagSenderSubmitPeriod),
			prettylogger.Err(err),
		)
		return err
	}

	submitTimeoutStr, err := er.db.GetConfigParameter(ctx, jacfarm.ConfigFlagSenderSubmitTimeout)
	if err != nil {
		log.Error(
			"error getting submit timeout from db config",
			slog.String("config_key", jacfarm.ConfigFlagSenderSubmitTimeout),
			prettylogger.Err(err),
		)
		return err
	}
	submitTimeout, err := time.ParseDuration(submitTimeoutStr)
	if err != nil {
		log.Error(
			"error parsing submit timeout from db config",
			slog.String("config_key", jacfarm.ConfigFlagSenderSubmitTimeout),
			prettylogger.Err(err),
		)
		return err
	}

	submitLimitStr, err := er.db.GetConfigParameter(ctx, jacfarm.ConfigFlagSenderSubmitLimit)
	if err != nil {
		log.Error(
			"error getting submit limit from db config",
			slog.String("config_key", jacfarm.ConfigFlagSenderSubmitLimit),
			prettylogger.Err(err),
		)
		return err
	}
	submitLimit, err := strconv.Atoi(submitLimitStr)
	if err != nil {
		log.Error(
			"error parsing submit limit from db config",
			slog.String("config_key", jacfarm.ConfigFlagSenderSubmitLimit),
			prettylogger.Err(err),
		)
		return err
	}

	er.cfg.flagTTL = flagTTL
	er.cfg.token = token
	er.cfg.juryFlagURL = juryFlagURL
	er.cfg.submitTimeout = submitTimeout
	er.cfg.submitPeriod = submitPeriod
	er.cfg.submitLimit = submitLimit

	return nil
}
