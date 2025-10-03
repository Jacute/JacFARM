package handlers

import (
	"JacFARM/internal/http/dto"

	fiber "github.com/gofiber/fiber/v3"
)

// ListFlags godoc
// @Summary Возврат полной информации о флагах
// @Param limit query int false "Количество флагов"
// @Param page query int false "Страница флагов"
// @Param status_id query int false "Фильтр по статусу"
// @Param exploit_id query int false "Фильтр по эксплойту"
// @Param team_id query int false "Фильтр по команде"
// @Produce json
// @Success 200 {object} dto.GetFlagsResponse
// @Failure 400 {object} dto.Response "Неверный запрос"
// @Failure 401 {object} dto.Response "Ошибка авторизации"
// @Failure 500 {object} dto.Response "Внутренняя ошибка"
// @Router /api/v1/flags [get]
// @Tags Flags
// @Security BasicAuth
func (h *Handlers) ListFlags() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		flagsFilter, err := dto.MapQueryToGetFlagsFilter(c.Queries())
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Error(err.Error()))
		}

		flags, count, err := h.service.ListFlags(c.RequestCtx(), flagsFilter)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternal)
		}

		return c.JSON(dto.GetFlagsResponse{
			Response: dto.OK(),
			Flags:    flags,
			Count:    count,
		})
	}
}

// PutFlag godoc
// @Summary Ручное добавление флага
// @Param request body dto.PutFlagRequest true "Флаг"
// @Produce json
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response "Неверный запрос"
// @Failure 401 {object} dto.Response "Ошибка авторизации"
// @Failure 500 {object} dto.Response "Внутренняя ошибка"
// @Router /api/v1/flags [post]
// @Tags Flags
// @Security BasicAuth
func (h *Handlers) PutFlag() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		if c.Get("Content-Type") != "application/json" {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrInvalidContentType)
		}

		var req dto.PutFlagRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrDecodingBody)
		}

		err := h.service.PutFlag(c.RequestCtx(), req.Flag)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternal)
		}

		return c.JSON(dto.OK())
	}
}

// GetStatuses godoc
// @Summary Получение статусов флагов
// @Produce json
// @Success 200 {object} dto.GetStatusesResponse
// @Failure 400 {object} dto.Response "Неверный запрос"
// @Failure 401 {object} dto.Response "Ошибка авторизации"
// @Failure 500 {object} dto.Response "Внутренняя ошибка"
// @Router /api/v1/flags/statuses [get]
// @Tags Flags
// @Security BasicAuth
func (h *Handlers) GetStatuses() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		statuses, err := h.service.GetStatuses(c.RequestCtx())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternal)
		}

		return c.JSON(dto.GetStatusesResponse{
			Response: dto.OK(),
			Statuses: statuses,
		})
	}
}

// GetFlagsCount godoc
// @Summary Возвращает количество флагов в очереди
// @Produce json
// @Success 200 {object} dto.GetFlagsCountResponse
// @Failure 400 {object} dto.Response "Неверный запрос"
// @Failure 401 {object} dto.Response "Ошибка авторизации"
// @Failure 500 {object} dto.Response "Внутренняя ошибка"
// @Router /api/v1/flags/count [get]
// @Tags Flags
// @Security BasicAuth
func (h *Handlers) GetFlagsCount() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		count, err := h.service.GetFlagsCount()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternal)
		}

		return c.JSON(dto.GetFlagsCountResponse{
			Response: dto.OK(),
			Count:    count,
		})
	}
}
