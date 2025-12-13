package handlers

import (
	"JacFARM/internal/http/dto"

	"github.com/gofiber/fiber/v3"
)

func (h *Handlers) ServicePutFlag() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		if c.Get("Content-Type") != "application/json" {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrInvalidContentType)
		}

		var req dto.ServicePutFlagRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrDecodingBody)
		}

		err := h.service.ServicePutFlag(c.RequestCtx(), &req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrInternal)
		}

		return c.JSON(dto.OK())
	}
}
