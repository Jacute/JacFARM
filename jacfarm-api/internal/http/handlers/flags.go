package handlers

import (
	"JacFARM/internal/http/dto"

	fiber "github.com/gofiber/fiber/v3"
)

func (h *Handlers) ListFlags() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		flagsFilter, err := dto.MapQueryToGetFlagsFilter(c.Queries())
		if err != nil {
			return c.JSON(dto.Error(err.Error()))
		}

		flags, count, err := h.service.ListFlags(c.RequestCtx(), flagsFilter)
		if err != nil {
			return c.JSON(dto.ErrInternal)
		}

		return c.JSON(dto.GetFlagsResponse{
			Response: dto.OK(),
			Flags:    flags,
			Count:    count,
		})
	}
}

func (h *Handlers) PutFlag() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		if c.Get("Content-Type") != "application/json" {
			return c.JSON(dto.ErrInvalidContentType)
		}

		var req dto.PutFlagRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.JSON(dto.ErrDecodingBody)
		}

		err := h.service.PutFlag(c.RequestCtx(), req.Flag)
		if err != nil {
			return c.JSON(dto.ErrInternal)
		}

		return c.JSON(dto.OK())
	}
}

func (h *Handlers) GetStatuses() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		statuses, err := h.service.GetStatuses(c.RequestCtx())
		if err != nil {
			return c.JSON(dto.ErrInternal)
		}

		return c.JSON(dto.GetStatusesResponse{
			Response: dto.OK(),
			Statuses: statuses,
		})
	}
}

func (h *Handlers) GetFlagsCount() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		count, err := h.service.GetFlagsCount()
		if err != nil {
			c = c.Status(fiber.StatusInternalServerError)
			return c.JSON(dto.ErrInternal)
		}

		return c.JSON(dto.GetFlagsCountResponse{
			Response: dto.OK(),
			Count:    count,
		})
	}
}
