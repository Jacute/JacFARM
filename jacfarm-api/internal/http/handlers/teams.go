package handlers

import (
	"JacFARM/internal/http/dto"

	"github.com/gofiber/fiber/v3"
)

func (h *Handlers) ListShortTeams() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		teams, err := h.service.ListShortTeams(c.RequestCtx())
		if err != nil {
			return c.JSON(dto.ErrInternal)
		}

		return c.JSON(dto.ListShortTeamsResponse{
			Response: dto.OK(),
			Teams:    teams,
		})
	}
}

func (h *Handlers) ListTeams() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		// TODO: implement me
		return c.JSON(dto.OK())
	}
}
