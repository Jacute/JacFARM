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
// @Produce json
// @Success 200 {object} dto.GetConfigResponse
// @Failure 400 {object} dto.Response "Неверный запрос"
// @Failure 401 {object} dto.Response "Ошибка авторизации"
// @Failure 500 {object} dto.Response "Внутренняя ошибка"
// @Router /api/v1/config [get]
// @Tags Auth
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

func (h *Handlers) UpdateConfig() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		if c.Get("Content-Type") != "application/json" {
			return c.JSON(dto.ErrInvalidContentType)
		}

		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.JSON(dto.Error("id should be uuid"))
		}

		var req dto.UpdateConfigRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.JSON(dto.ErrDecodingBody)
		}

		err = h.service.UpdateConfig(c.RequestCtx(), int64(id), req.Value)
		if err != nil {
			if errors.Is(err, storage.ErrConfigParamNotFound) {
				return c.JSON(dto.Error(err.Error()))
			}
			return c.JSON(dto.ErrInternal)
		}

		return c.JSON(dto.OK())
	}
}
