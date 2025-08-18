package handlers

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"context"
)

type ServiceInterface interface {
	ListFlags(ctx context.Context, filter *dto.ListFlagsFilter) ([]*models.FlagEnrich, int, error)
	PutFlag(ctx context.Context, flag string) error

	ListExploits(ctx context.Context, filter *dto.ListExploitsFilter) ([]*models.Exploit, error)
	ListShortExploits(ctx context.Context) ([]*models.ExploitShort, error)
	ToggleExploit(ctx context.Context, id string) (bool, error)
	UploadExploit(ctx context.Context, req *dto.UploadExploitRequest) (string, error)
	DeleteExploit(ctx context.Context, id string) error

	ListShortTeams(ctx context.Context) ([]*models.ShortTeam, error)
}

type Handlers struct {
	service ServiceInterface
}

func New(service ServiceInterface) *Handlers {
	return &Handlers{
		service: service,
	}
}
