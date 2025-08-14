package handlers

import (
	"JacFARM/internal/http/dto"

	fiber "github.com/gofiber/fiber/v3"
)

func (h *Handlers) GetFlags() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		flagsFilter, err := dto.MapQueryToGetFlagsFilter(c.Queries())
		if err != nil {
			return c.JSON(dto.Error(err))
		}

		flags, err := h.service.GetFlags(c.RequestCtx(), flagsFilter)
		if err != nil {
			return c.JSON(dto.InternalError())
		}

		return c.JSON(flags)
	}
}
