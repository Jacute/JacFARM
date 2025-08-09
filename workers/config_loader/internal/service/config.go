package service

import (
	"config_loader/internal/config"
	"config_loader/internal/models"
	"config_loader/internal/storage/postgres"
	"context"
	"errors"
	"log/slog"
	"strconv"

	"github.com/jacute/prettylogger"
)

const (
	ConfigFlagFormatKey           = "EXPLOIT_RUNNER_FLAG_FORMAT"
	ConfigExploitDuration         = "EXPLOIT_RUNNER_DURATION"
	ConfigExploitMaxWorkingTime   = "EXPLOIT_RUNNER_MAX_WORKING_TIME"
	ConfigMaxConcurrentExploits   = "EXPLOIT_RUNNER_MAX_CONCURRENT_EXPLOITS"
	ConfigFlagSenderPlugin        = "FLAG_SENDER_PLUGIN"
	ConfigFlagSenderSubmitTimeout = "FLAG_SENDER_SUBMIT_TIMEOUT"
	ConfigFlagSenderSubmitPeriod  = "FLAG_SENDER_SUBMIT_PERIOD"
	ConfigFlagSenderJuryFlagURL   = "FLAG_SENDER_JURY_FLAG_URL_OR_HOST"
	ConfigFlagSenderToken         = "FLAG_SENDER_TOKEN"
	ConfigFlagSenderSubmitLimit   = "FLAG_SENDER_SUBMIT_LIMIT"
	ConfigFlagSenderFlagTTL       = "FLAG_SENDER_FLAG_TTL"
)

func (s *Service) LoadConfigIntoDB(ctx context.Context, cfg *config.Config) {
	const op = "service.jacfarm.LoadConfigIntoDB"
	log := s.log.With(slog.String("op", op))

	var existTeams []string
	for _, ip := range cfg.ExploitRunner.TeamIPs {
		_, err := s.db.AddTeam(ctx, &models.Team{
			IP: ip,
		})
		if err != nil {
			if errors.Is(err, postgres.ErrTeamAlreadyExists) {
				existTeams = append(existTeams, ip)
				continue
			}
			log.Warn("team cannot be added", prettylogger.Err(err))
		}
	}
	if len(existTeams) > 0 {
		log.Info("some teams already exist", slog.Any("teams", existTeams))
	}

	var existParams []string
	configMap := map[string]string{
		ConfigFlagFormatKey:         cfg.ExploitRunner.FlagFormat,
		ConfigExploitDuration:       cfg.ExploitRunner.RunDuration.String(),
		ConfigExploitMaxWorkingTime: cfg.ExploitRunner.ExploitMaxWorkingTime.String(),
		ConfigMaxConcurrentExploits: strconv.Itoa(cfg.ExploitRunner.MaxConcurrentExploits),

		ConfigFlagSenderFlagTTL:       cfg.FlagSender.FlagTTL.String(),
		ConfigFlagSenderJuryFlagURL:   cfg.FlagSender.JuryFlagURL,
		ConfigFlagSenderPlugin:        cfg.FlagSender.Plugin,
		ConfigFlagSenderSubmitTimeout: cfg.FlagSender.SubmitTimeout.String(),
		ConfigFlagSenderSubmitLimit:   strconv.Itoa(cfg.FlagSender.SubmitLimit),
		ConfigFlagSenderSubmitPeriod:  cfg.FlagSender.SubmitPeriod.String(),
		ConfigFlagSenderToken:         cfg.FlagSender.Token,
	}

	for k, v := range configMap {
		_, err := s.db.AddConfigParameter(ctx, k, v)
		if err != nil {
			if errors.Is(err, postgres.ErrConfigParamAlreadyExists) {
				existParams = append(existParams, k)
				continue
			}
			log.Warn("config param cannot be added", slog.String("param", k), prettylogger.Err(err))
		}
	}

	if len(existParams) > 0 {
		log.Info("some config params already exist", slog.Any("params", existParams))
	}
}
