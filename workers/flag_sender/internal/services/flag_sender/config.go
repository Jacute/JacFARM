package flag_sender

import (
	"context"
	"flag_sender/pkg/common_config"
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
		plugin, err := er.db.GetConfigParameter(ctx, common_config.ConfigFlagSenderPlugin)
		if err != nil {
			log.Error(
				"error getting flag sender plugin from db config",
				slog.String("config_key", common_config.ConfigFlagSenderPlugin),
				prettylogger.Err(err),
			)
			return err
		}
		er.cfg.plugin = plugin
	}
	juryFlagURL, err := er.db.GetConfigParameter(ctx, common_config.ConfigFlagSenderJuryFlagURL)
	if err != nil {
		log.Error(
			"error getting jury flag url from db config",
			slog.String("config_key", common_config.ConfigFlagSenderJuryFlagURL),
			prettylogger.Err(err),
		)
		return err
	}

	token, err := er.db.GetConfigParameter(ctx, common_config.ConfigFlagSenderToken)
	if err != nil {
		log.Error(
			"error getting flag sender token from db config",
			slog.String("config_key", common_config.ConfigFlagSenderToken),
			prettylogger.Err(err),
		)
		return err
	}

	flagTTLstr, err := er.db.GetConfigParameter(ctx, common_config.ConfigFlagSenderFlagTTL)
	if err != nil {
		log.Error(
			"error getting flag ttl from db config",
			slog.String("config_key", common_config.ConfigFlagSenderFlagTTL),
			prettylogger.Err(err),
		)
		return err
	}
	flagTTL, err := time.ParseDuration(flagTTLstr)
	if err != nil {
		log.Error(
			"error parsing flag ttl from db config",
			slog.String("config_key", common_config.ConfigFlagSenderFlagTTL),
			prettylogger.Err(err),
		)
		return err
	}

	submitPeriodStr, err := er.db.GetConfigParameter(ctx, common_config.ConfigFlagSenderSubmitPeriod)
	if err != nil {
		log.Error(
			"error getting submit period from db config",
			slog.String("config_key", common_config.ConfigFlagSenderSubmitPeriod),
			prettylogger.Err(err),
		)
		return err
	}
	submitPeriod, err := time.ParseDuration(submitPeriodStr)
	if err != nil {
		log.Error(
			"error parsing submit period from db config",
			slog.String("config_key", common_config.ConfigFlagSenderSubmitPeriod),
			prettylogger.Err(err),
		)
		return err
	}

	submitTimeoutStr, err := er.db.GetConfigParameter(ctx, common_config.ConfigFlagSenderSubmitTimeout)
	if err != nil {
		log.Error(
			"error getting submit timeout from db config",
			slog.String("config_key", common_config.ConfigFlagSenderSubmitTimeout),
			prettylogger.Err(err),
		)
		return err
	}
	submitTimeout, err := time.ParseDuration(submitTimeoutStr)
	if err != nil {
		log.Error(
			"error parsing submit timeout from db config",
			slog.String("config_key", common_config.ConfigFlagSenderSubmitTimeout),
			prettylogger.Err(err),
		)
		return err
	}

	submitLimitStr, err := er.db.GetConfigParameter(ctx, common_config.ConfigFlagSenderSubmitLimit)
	if err != nil {
		log.Error(
			"error getting submit limit from db config",
			slog.String("config_key", common_config.ConfigFlagSenderSubmitLimit),
			prettylogger.Err(err),
		)
		return err
	}
	submitLimit, err := strconv.Atoi(submitLimitStr)
	if err != nil {
		log.Error(
			"error parsing submit limit from db config",
			slog.String("config_key", common_config.ConfigFlagSenderSubmitLimit),
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
