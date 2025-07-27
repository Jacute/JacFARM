package jacfarm

import (
	"JacFARM/internal/config"
	"JacFARM/internal/models"
	"context"
	"log/slog"
	"strconv"

	"github.com/jacute/prettylogger"
)

const (
	ConfigFlagFormatKey         = "EXPLOIT_RUNNER_FLAG_FORMAT"
	ConfigExploitDuration       = "EXPLOIT_RUNNER_DURATION"
	ConfigExploitMaxWorkingTime = "EXPLOIT_RUNNER_MAX_WORKING_TIME"
	ConfigMaxConcurrentExploits = "EXPLOIT_RUNNER_MAX_CONCURRENT_EXPLOITS"
)

func (s *Service) LoadConfigIntoDB(ctx context.Context, cfg *config.Config) {
	const op = "service.jacfarm.LoadConfigIntoDB"
	log := s.log.With(slog.String("op", op))

	for _, ip := range cfg.ExploitRunner.TeamIPs {
		err := s.db.AddTeam(&models.Team{
			IP: ip,
		})
		if err != nil {
			log.Warn("team cannot be added", prettylogger.Err(err))
			continue
		}
	}

	err := s.db.AddConfigParameter(ctx, ConfigFlagFormatKey, cfg.ExploitRunner.FlagFormat)
	if err != nil {
		log.Warn("error adding flag_format", prettylogger.Err(err))
	}

	err = s.db.AddConfigParameter(ctx, ConfigExploitDuration, cfg.ExploitRunner.RunDuration.String())
	if err != nil {
		log.Warn("error adding exploit_run_duration", prettylogger.Err(err))
	}

	err = s.db.AddConfigParameter(ctx, ConfigExploitMaxWorkingTime, cfg.ExploitRunner.ExploitMaxWorkingTime.String())
	if err != nil {
		log.Warn("error adding exploit_max_working_time", prettylogger.Err(err))
	}
	err = s.db.AddConfigParameter(ctx, ConfigMaxConcurrentExploits, strconv.Itoa(cfg.ExploitRunner.MaxConcurrentExploits))
	if err != nil {
		log.Warn("error adding max_concurrent_exploits", prettylogger.Err(err))
	}
}
