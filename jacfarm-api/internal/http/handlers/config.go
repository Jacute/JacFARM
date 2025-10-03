package handlers

import (
	"JacFARM/internal/http/dto"
	"JacFARM/internal/storage"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// GetConfig godoc
// @Summary Возвращает конфигурацию JacFARM
// @Param limit query int false "Количество конфигураций"
// @Param page query int false "Страница конфигурации"
// @Produce json
// @Success 200 {object} dto.GetConfigResponse
// @Failure 400 {object} dto.Response "Неверный запрос"
// @Failure 401 {object} dto.Response "Ошибка авторизации"
// @Failure 500 {object} dto.Response "Внутренняя ошибка"
// @Router /api/v1/config [get]
// @Tags Config
// @Security BasicAuth
func (h *Handlers) GetConfig() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		filter, err := dto.MapQueryToGetConfigFilter(c.Queries())
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Error(err.Error()))
		}

		config, count, err := h.service.GetConfig(c.RequestCtx(), filter)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternal)
		}

		return c.JSON(dto.GetConfigResponse{
			Response: dto.OK(),
			Config:   config,
			Count:    count,
		})
	}
}

// UpdateConfig godoc
// @Summary Возвращает конфигурацию JacFARM
// @Produce json
// @Param id path int true "Config ID"
// @Param request body dto.UpdateConfigRequest true "Информация о конфигурации"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response "Неверный запрос"
// @Failure 401 {object} dto.Response "Ошибка авторизации"
// @Failure 404 {object} dto.Response "Конфигурация не найдена"
// @Failure 500 {object} dto.Response "Внутренняя ошибка"
// @Router /api/v1/config/{id} [patch]
// @Tags Config
// @Security BasicAuth
func (h *Handlers) UpdateConfig() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		if c.Get("Content-Type") != "application/json" {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrInvalidContentType)
		}

		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.Error("id should be uuid"))
		}

		var req dto.UpdateConfigRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrDecodingBody)
		}

		err = h.service.UpdateConfig(c.RequestCtx(), int64(id), req.Value)
		if err != nil {
			if errors.Is(err, storage.ErrConfigParamNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(dto.Error(err.Error()))
			}
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternal)
		}

		return c.JSON(dto.OK())
	}
}
