package service

import (
	"config_loader/internal/config"
	"config_loader/internal/models"
	"config_loader/internal/postgres"
	"config_loader/internal/utils"
	"config_loader/pkg/common_config"
	"context"
	"errors"
	"log/slog"
	"strconv"

	"github.com/jacute/prettylogger"
)

func (s *Service) LoadConfigIntoDB(ctx context.Context, cfg *config.Config) error {
	const op = "service.jacfarm.LoadConfigIntoDB"
	log := s.log.With(slog.String("op", op))

	var existTeams []string

	// add single ips
	err := s.addIps(ctx, cfg.ExploitRunner.TeamIPs, existTeams)
	if err != nil {
		log.Warn("error adding ips", slog.Any("ips", cfg.ExploitRunner.TeamIPs), prettylogger.Err(err))
	}

	// add ip ranges
	for _, ipRange := range cfg.ExploitRunner.TeamIPRange {
		ips, err := utils.ExpandRange(ipRange)
		if err != nil {
			log.Warn("error adding ip range", slog.String("ip_range", ipRange), prettylogger.Err(err))
			continue
		}
		err = s.addIps(ctx, ips, existTeams)
		if err != nil {
			log.Warn("error adding ips", slog.Any("ips", ips), prettylogger.Err(err))
		}
	}

	// add ip cidrs
	for _, ipCIDR := range cfg.ExploitRunner.TeamIPCidr {
		ips, err := utils.ExpandCIDR(ipCIDR)
		if err != nil {
			log.Warn("error adding ip cidr", slog.String("ip_cidr", ipCIDR), prettylogger.Err(err))
			continue
		}
		err = s.addIps(ctx, ips, existTeams)
		if err != nil {
			log.Warn("error adding ips", slog.Any("ips", ips), prettylogger.Err(err))
		}
	}

	// add ip from N
	ips := utils.ExpandIpFromN(
		cfg.ExploitRunner.TeamIPFromN.NStart,
		cfg.ExploitRunner.TeamIPFromN.NEnd,
		cfg.ExploitRunner.TeamIPFromN.OffsetX,
		cfg.ExploitRunner.TeamIPFromN.OffsetY,
		cfg.ExploitRunner.TeamIPFromN.Block,
		cfg.ExploitRunner.TeamIPFromN.IPTemplate,
	)
	err = s.addIps(ctx, ips, existTeams)
	if err != nil {
		log.Warn("error adding ips", slog.Any("ips", ips), prettylogger.Err(err))
	}

	if len(existTeams) > 0 {
		log.Info("some teams already exist", slog.Any("teams", existTeams))
	}

	var existParams []string
	configMap := map[string]string{
		common_config.ConfigFlagFormatKey:         cfg.ExploitRunner.FlagFormat,
		common_config.ConfigExploitDuration:       cfg.ExploitRunner.RunDuration.String(),
		common_config.ConfigExploitMaxWorkingTime: cfg.ExploitRunner.ExploitMaxWorkingTime.String(),
		common_config.ConfigMaxConcurrentExploits: strconv.Itoa(cfg.ExploitRunner.MaxConcurrentExploits),
		common_config.ConfigVenvMaxInstallTime:    cfg.ExploitRunner.VenvMaxInstallTime.String(),

		common_config.ConfigFlagSenderFlagTTL:       cfg.FlagSender.FlagTTL.String(),
		common_config.ConfigFlagSenderJuryFlagURL:   cfg.FlagSender.JuryFlagURL,
		common_config.ConfigFlagSenderPlugin:        cfg.FlagSender.Plugin,
		common_config.ConfigFlagSenderSubmitTimeout: cfg.FlagSender.SubmitTimeout.String(),
		common_config.ConfigFlagSenderSubmitLimit:   strconv.Itoa(cfg.FlagSender.SubmitLimit),
		common_config.ConfigFlagSenderSubmitPeriod:  cfg.FlagSender.SubmitPeriod.String(),
		common_config.ConfigFlagSenderToken:         cfg.FlagSender.Token,
	}

	for k, v := range configMap {
		_, err := s.db.AddConfigParameter(ctx, k, v)
		if err != nil {
			if errors.Is(err, postgres.ErrConfigParamAlreadyExists) {
				existParams = append(existParams, k)
				continue
			}
			log.Warn("config param cannot be added", slog.String("param", k), prettylogger.Err(err))
			return err
		}
	}

	if len(existParams) > 0 {
		log.Info("some config params already exist", slog.Any("params", existParams))
	}

	return nil
}

func (s *Service) addIps(ctx context.Context, ips []string, existTeams []string) error {
	const op = "service.jacfarm.addIps"
	log := s.log.With(slog.String("op", op))

	for _, ip := range ips {
		_, err := s.db.AddTeam(ctx, &models.Team{
			IP: ip,
		})
		if err != nil {
			if errors.Is(err, postgres.ErrTeamAlreadyExists) {
				existTeams = append(existTeams, ip)
				continue
			}
			log.Warn("team cannot be added", prettylogger.Err(err))
			return err
		}
	}

	return nil
}
