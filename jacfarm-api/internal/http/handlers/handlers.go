package handlers

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"context"
)

type ServiceInterface interface {
	ListFlags(ctx context.Context, filter *dto.ListFlagsFilter) ([]*models.FlagEnrich, error)
	PutFlag(ctx context.Context, flag string) error
	ListExploits(ctx context.Context, filter *dto.ListExploitsFilter) ([]*models.Exploit, error)
}

type Handlers struct {
	service ServiceInterface
}

func New(service ServiceInterface) *Handlers {
	return &Handlers{
		service: service,
	}
}
