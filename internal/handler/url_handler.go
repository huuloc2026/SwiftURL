package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/huuloc2026/SwiftURL/internal/service"
)

type URLHandler struct {
	svc service.URLService
}

func NewURLHandler(s service.URLService) *URLHandler {
	return &URLHandler{svc: s}
}

func (h *URLHandler) ShortenURL(c *fiber.Ctx) error {
	type req struct {
		URL string `json:"url"`
	}
	var r req
	if err := c.BodyParser(&r); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	result, err := h.svc.Shorten(r.URL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"shortened": "http://localhost:8080/" + result.Code})
}

func (h *URLHandler) ResolveURL(c *fiber.Ctx) error {
	code := c.Params("code")
	result, err := h.svc.Resolve(code)
	if err != nil {
		return c.Status(404).SendString("Not found")
	}
	return c.Redirect(result.LongURL, fiber.StatusMovedPermanently)
}
