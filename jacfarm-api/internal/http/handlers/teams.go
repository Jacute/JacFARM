package handlers

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/http/validator"
	"JacFARM/internal/models"
	"JacFARM/internal/storage"
	"net"

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
		filter, err := dto.MapQueryToListTeamsFilter(c.Queries())
		if err != nil {
			return c.JSON(dto.Error(err.Error()))
		}

		teams, count, err := h.service.ListTeams(c.RequestCtx(), filter)
		if err != nil {
			return c.JSON(dto.ErrInternal)
		}

		return c.JSON(dto.ListTeamsResponse{
			Response: dto.OK(),
			Teams:    teams,
			Count:    count,
		})
	}
}

func (h *Handlers) AddTeam() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		if c.Get("Content-Type") != "application/json" {
			return c.JSON(dto.ErrInvalidContentType)
		}

		var req dto.AddTeamRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.JSON(dto.ErrDecodingBody)
		}

		if err := validator.Validate(req); err != nil {
			return c.JSON(dto.Error(err.Error()))
		}

		id, err := h.service.AddTeam(c.RequestCtx(), &models.Team{
			Name: req.Name,
			IP:   net.ParseIP(req.IP),
		})
		if err != nil {
			if err == storage.ErrTeamAlreadyExists {
				return c.JSON(dto.Error("team already exists"))
			}
			return c.JSON(dto.ErrInternal)
		}

		return c.JSON(dto.AddTeamResponse{
			Response: dto.OK(),
			ID:       id,
		})
	}
}
