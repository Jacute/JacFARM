package handlers

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/http/validator"
	"JacFARM/internal/models"
	"JacFARM/internal/storage"
	"errors"
	"net"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// ListShortTeams godoc
// @Summary Возврат короткой информации о командах
// @Produce json
// @Success 200 {object} dto.ListShortTeamsResponse
// @Failure 401 {object} dto.Response "Ошибка авторизации"
// @Failure 500 {object} dto.Response "Внутренняя ошибка"
// @Router /api/v1/flags/short [get]
// @Tags Flags
// @Security BasicAuth
func (h *Handlers) ListShortTeams() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		teams, err := h.service.ListShortTeams(c.RequestCtx())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternal)
		}

		return c.JSON(dto.ListShortTeamsResponse{
			Response: dto.OK(),
			Teams:    teams,
		})
	}
}

// ListTeams godoc
// @Summary Возврат полной информации о командах
// @Param limit query int false "Количество команд"
// @Param page query int false "Страница команд"
// @Param query query string false "Поиск по названию/ip команды"
// @Produce json
// @Success 200 {object} dto.ListTeamsResponse
// @Failure 400 {object} dto.Response "Неверный запрос"
// @Failure 401 {object} dto.Response "Ошибка авторизации"
// @Failure 500 {object} dto.Response "Внутренняя ошибка"
// @Router /api/v1/teams [get]
// @Tags Teams
// @Security BasicAuth
func (h *Handlers) ListTeams() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		filter, err := dto.MapQueryToListTeamsFilter(c.Queries())
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Error(err.Error()))
		}

		teams, count, err := h.service.ListTeams(c.RequestCtx(), filter)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternal)
		}

		return c.JSON(dto.ListTeamsResponse{
			Response: dto.OK(),
			Teams:    teams,
			Count:    count,
		})
	}
}

// AddTeam godoc
// @Summary Добавление команды
// @Param request body dto.AddTeamRequest true "Команда"
// @Produce json
// @Success 200 {object} dto.AddTeamResponse
// @Failure 400 {object} dto.Response "Неверный запрос"
// @Failure 401 {object} dto.Response "Ошибка авторизации"
// @Failure 500 {object} dto.Response "Внутренняя ошибка"
// @Router /api/v1/teams [post]
// @Tags Teams
// @Security BasicAuth
func (h *Handlers) AddTeam() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		if c.Get("Content-Type") != "application/json" {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrInvalidContentType)
		}

		var req dto.AddTeamRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrDecodingBody)
		}

		if err := validator.Validate(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Error(err.Error()))
		}

		id, err := h.service.AddTeam(c.RequestCtx(), &models.Team{
			Name: req.Name,
			IP:   net.ParseIP(req.IP),
		})
		if err != nil {
			if err == storage.ErrTeamAlreadyExists {
				return c.Status(fiber.StatusBadRequest).JSON(dto.Error("team already exists"))
			}
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternal)
		}

		return c.JSON(dto.AddTeamResponse{
			Response: dto.OK(),
			ID:       id,
		})
	}
}

// DeleteTeam godoc
// @Summary Удаление команды
// @Param id path int true "Team ID"
// @Produce json
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response "Неверный запрос"
// @Failure 401 {object} dto.Response "Ошибка авторизации"
// @Failure 404 {object} dto.Response "Команда не найдена"
// @Failure 500 {object} dto.Response "Внутренняя ошибка"
// @Router /api/v1/teams/{id} [delete]
// @Tags Teams
// @Security BasicAuth
func (h *Handlers) DeleteTeam() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Error("id should be int"))
		}

		err = h.service.DeleteTeam(c.RequestCtx(), int64(id))
		if err != nil {
			if errors.Is(err, storage.ErrTeamNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(dto.Error(err.Error()))
			}
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternal)
		}

		return c.JSON(dto.OK())
	}
}
