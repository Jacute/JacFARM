package handlers

import (
	"JacFARM/internal/http/dto"

	"github.com/gofiber/fiber/v3"
)

// ListLogs godoc
// @Summary Возврат полной информации о логах
// @Param limit query int false "Количество логов"
// @Param page query int false "Страница логов"
// @Param exploit_id query int false "Идентификатор эксплойта"
// @Param module_id query int false "Идентификатор модуля"
// @Param log_level_id query int false "Идентификатор уровня логирования"
// @Produce json
// @Success 200 {object} dto.ListLogsResponse
// @Failure 400 {object} dto.Response "Неверный запрос"
// @Failure 401 {object} dto.Response "Ошибка авторизации"
// @Failure 500 {object} dto.Response "Внутренняя ошибка"
// @Router /api/v1/logs [get]
// @Tags Logs
// @Security BasicAuth
func (h *Handlers) ListLogs() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		filter, err := dto.MapQueryToListLogsFilter(c.Queries())
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Error(err.Error()))
		}

		logs, count, err := h.service.ListLogs(c.RequestCtx(), filter)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternal)
		}

		return c.JSON(dto.ListLogsResponse{
			Response: dto.OK(),
			Logs:     logs,
			Count:    count,
		})
	}
}

// ListModules godoc
// @Summary Возврат всех модулей фермы
// @Produce json
// @Success 200 {object} dto.ListModulesResponse
// @Failure 400 {object} dto.Response "Неверный запрос"
// @Failure 401 {object} dto.Response "Ошибка авторизации"
// @Failure 500 {object} dto.Response "Внутренняя ошибка"
// @Router /api/v1/logs/modules [get]
// @Tags Logs
// @Security BasicAuth
func (h *Handlers) ListModules() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		modules, count, err := h.service.ListModules(c.RequestCtx())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternal)
		}

		return c.JSON(dto.ListModulesResponse{
			Response: dto.OK(),
			Modules:  modules,
			Count:    count,
		})
	}
}

// ListLogLevels godoc
// @Summary Возврат всех уровней логирования
// @Produce json
// @Success 200 {object} dto.ListLogLevelsResponse
// @Failure 400 {object} dto.Response "Неверный запрос"
// @Failure 401 {object} dto.Response "Ошибка авторизации"
// @Failure 500 {object} dto.Response "Внутренняя ошибка"
// @Router /api/v1/logs/levels [get]
// @Tags Logs
// @Security BasicAuth
func (h *Handlers) ListLogLevels() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		logLevels, count, err := h.service.ListLogLevel(c.RequestCtx())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternal)
		}

		return c.JSON(dto.ListLogLevelsResponse{
			Response:  dto.OK(),
			LogLevels: logLevels,
			Count:     count,
		})
	}
}
