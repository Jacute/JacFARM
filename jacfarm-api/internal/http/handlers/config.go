package handlers

import (
	"JacFARM/internal/http/dto"

	"github.com/gofiber/fiber/v3"
)

func (h *Handlers) GetConfig() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		filter, err := dto.MapQueryToGetConfigFilter(c.Queries())
		if err != nil {
			return c.JSON(dto.Error(err.Error()))
		}

		config, count, err := h.service.GetConfig(c.RequestCtx(), filter)
		if err != nil {
			return c.JSON(dto.ErrInternal)
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

		// TODO: add update config
		// var req dto.PutConfigRequest
		// if err := c.Bind().Body(&req); err != nil {
		// 	return c.JSON(dto.ErrDecodingBody)
		// }

		// err := h.service.PutConfig(c.RequestCtx(), req.Value)
		// if err != nil {
		// 	return c.JSON(dto.ErrInternal)
		// }

		return c.JSON(dto.OK())
	}
}
