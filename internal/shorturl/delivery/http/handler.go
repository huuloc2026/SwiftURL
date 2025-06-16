package handler

import (
	"net/http"
	"time"

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

type CreateShortURLRequest struct {
	LongURL    string  `json:"long_url" binding:"required,url"`
	CustomCode *string `json:"custom_code,omitempty"`
	ExpireAt   *string `json:"expire_at,omitempty"`
}

type CreateShortURLResponse struct {
	ShortCode string `json:"short_code"`
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

// POST /api/shorten
func (h *URLHandler) CreateShortURL(ctx *fiber.Ctx) error {
	var req CreateShortURLRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.LongURL == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "long_url is required",
		})
	}

	var expireTime *time.Time
	if req.ExpireAt != nil {
		t, err := time.Parse(time.RFC3339, *req.ExpireAt)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid expire_at format (use RFC3339)",
			})
		}
		expireTime = &t
	}

	code, err := h.usecase.Generate(ctx.Context(), req.LongURL, req.CustomCode, expireTime)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate short url: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(CreateShortURLResponse{
		ShortCode: code,
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

// DELETE /api/shorten/:shortCode
func (h *URLHandler) DeleteShortURL(ctx *fiber.Ctx) error {
	code := ctx.Params("shortCode")
	if err := h.usecase.Delete(ctx.Context(), code); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete: " + err.Error(),
		})
	}
	return ctx.JSON(fiber.Map{"message": "deleted"})
}
