package shortlink_transport_http

import (
	"api/internal/core/domain"
	core_dto "api/internal/core/dto"
	"api/internal/core/validator"

	"context"
	"time"

	"github.com/gofiber/fiber/v3"
)

func (h *ShortlinkHTTPHandler) ShortLinkHandler(c fiber.Ctx) error {
	var s core_dto.ShortLinkDto
	if err := c.Bind().Body(&s); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := validator.Valide(&s); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID, ok := c.Locals("user_id").(int)
	if !ok {
      return c.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
			})
    }

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	shortlink, err := h.shortlinkService.ShortlinkPost(ctx, userID, s.Link)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"url": shortlink,
	})
}

func (h *ShortlinkHTTPHandler) ShortLinkGetHandler(c fiber.Ctx) error {
	link := c.Params("link")

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	device := domain.NewDevice(c.IP(), c.UserAgent())
	sl, err := h.shortlinkService.ShortlinkGet(ctx, link, device)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(303).JSON(fiber.Map{
		"link": sl.OriginalUrl,
	})
}
	
func (h *ShortlinkHTTPHandler) ShortLinkActiviteHandler(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
      return c.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
			})
    }

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	sli, err := h.shortlinkService.ShortLinkActivity(ctx, userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(sli)
}