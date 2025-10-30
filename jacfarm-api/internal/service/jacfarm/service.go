package jacfarm

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"JacFARM/pkg/rabbitmq_dto"
	"context"
	"log/slog"
)

//go:generate mockgen -source=service.go -destination=./mocks/storage_mock.go -package=mocks -mock_names=storage=StorageMock storage
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
	DeleteTeam(ctx context.Context, id int64) error

	GetStatuses(ctx context.Context) ([]*models.Status, error)

	UpdateConfigRow(ctx context.Context, id int64, value string) error
	GetConfig(ctx context.Context, filter *dto.GetConfigFilter) ([]*models.Config, int, error)

	ListLogs(ctx context.Context, filter *dto.ListLogsFilter) ([]*models.Log, int, error)
	ListModules(ctx context.Context) ([]*models.Module, int, error)
	ListLogLevel(ctx context.Context) ([]*models.LogLevel, int, error)
}

//go:generate mockgen -source=service.go -destination=./mocks/queue_mock.go -package=mocks -mock_names=queue=QueueMock queue
type queue interface {
	PublishFlag(flag *rabbitmq_dto.Flag) error
	GetFlagsCount() (int, error)
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
