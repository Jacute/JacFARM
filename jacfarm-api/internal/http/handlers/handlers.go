package handlers

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/models"
	"context"
)

type ServiceInterface interface {
	GetFlags(ctx context.Context, filter *dto.GetFlagsFilter) ([]*models.FlagEnrich, error)
	PutFlag(ctx context.Context, flag string) error
}

type Handlers struct {
	service ServiceInterface
}

func New(service ServiceInterface) *Handlers {
	return &Handlers{
		service: service,
	}
}
