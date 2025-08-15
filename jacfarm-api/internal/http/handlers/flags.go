package handlers

import (
	"JacFARM/internal/http/dto"
	"slices"

	fiber "github.com/gofiber/fiber/v3"
)

func (h *Handlers) ListFlags() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		flagsFilter, err := dto.MapQueryToGetFlagsFilter(c.Queries())
		if err != nil {
			return c.JSON(dto.Error(err.Error()))
		}

		flags, err := h.service.ListFlags(c.RequestCtx(), flagsFilter)
		if err != nil {
			return c.JSON(dto.ErrInternal)
		}

		return c.JSON(dto.GetFlagsResponse{
			Response: dto.OK(),
			Flags:    flags,
		})
	}
}

func (h *Handlers) PutFlag() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		if !slices.Contains(c.GetReqHeaders()["Content-Type"], "application/json") {
			return c.JSON(dto.ErrInvalidContentType)
		}

		var req dto.PutFlagRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.JSON(dto.ErrDecodingBody)
		}

		err := h.service.PutFlag(c.RequestCtx(), req.Flag)
		if err != nil {
			c.JSON(dto.ErrInternal)
		}

		return c.JSON(dto.OK())
	}
}
