package jacfarm

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"JacFARM/pkg/rabbitmq_dto"
	"context"
	"log/slog"
)

type storage interface {
	GetFlags(ctx context.Context, filter *dto.ListFlagsFilter) ([]*models.FlagEnrich, int, error)

	GetExploits(ctx context.Context, filter *dto.ListExploitsFilter) ([]*models.Exploit, int, error)
	GetShortExploits(ctx context.Context) ([]*models.ExploitShort, error)
	ToggleExploit(ctx context.Context, id string) (bool, error)
	CreateExploit(ctx context.Context, exploit *models.Exploit) error
	DeleteExploit(ctx context.Context, id string) error

	AddTeam(ctx context.Context, team *models.Team) (int64, error)
	GetShortTeams(ctx context.Context) ([]*models.ShortTeam, error)
	GetTeams(ctx context.Context, filter *dto.ListTeamsFilter) ([]*models.Team, int, error)

	GetStatuses(ctx context.Context) ([]*models.Status, error)
}

type queue interface {
	PublishFlag(flag *rabbitmq_dto.Flag) error
}

type Service struct {
	log        *slog.Logger
	db         storage
	que        queue
	exploitDir string
}

func New(log *slog.Logger, db storage, que queue, exploitDir string) *Service {
	return &Service{
		log:        log,
		db:         db,
		que:        que,
		exploitDir: exploitDir,
	}
}
