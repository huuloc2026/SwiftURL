package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/huuloc2026/SwiftURL/internal/shorturl/usecase"
)

type URLHandler struct {
	usecase usecase.ShortURLUsecase
}

func NewURLHandler(u usecase.ShortURLUsecase) *URLHandler {
	return &URLHandler{usecase: u}
}

type shortenRequest struct {
	URL string `json:"url"`
}

func (h *URLHandler) ShortenURL(c *fiber.Ctx) error {
	var req shortenRequest
	if err := c.BodyParser(&req); err != nil || req.URL == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid input",
		})
	}

	result, err := h.usecase.Shorten(c.Context(), req.URL)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"short_code": result.ShortCode,
		"long_url":   result.LongURL,
	})
}

func (h *URLHandler) ResolveURL(c *fiber.Ctx) error {
	code := c.Params("code")
	result, err := h.usecase.Resolve(c.Context(), code)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "short code not found",
		})
	}
	return c.Redirect(result.LongURL, http.StatusMovedPermanently)
}
