package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/huuloc2026/SwiftURL/internal/shorturl/usecase"
	"github.com/huuloc2026/SwiftURL/pkg/response"
	"github.com/huuloc2026/SwiftURL/pkg/utils"
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
		return response.Error(c, http.StatusBadRequest, err, "invalid input")
	}
	result, err := h.usecase.Shorten(c.Context(), req.URL)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err, "internal server error")
	}
	return response.Success(c, fiber.Map{
		"short_code": result.ShortCode,
		"long_url":   result.LongURL,
	})
}

// POST /api/shorten
func (h *URLHandler) CreateShortURL(ctx *fiber.Ctx) error {
	var req CreateShortURLRequest

	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, fiber.StatusBadRequest, err, "invalid request body")
	}

	if req.LongURL == "" {
		return response.Error(ctx, fiber.StatusBadRequest, nil, "long_url is required")
	}

	var expireTime *time.Time
	if req.ExpireAt != nil {
		t, err := time.Parse(time.RFC3339, *req.ExpireAt)
		if err != nil {
			return response.Error(ctx, fiber.StatusBadRequest, err, "invalid expire_at format (use RFC3339)")
		}
		expireTime = &t
	}

	code, err := h.usecase.Generate(ctx.Context(), req.LongURL, req.CustomCode, expireTime)
	if err != nil {
		return response.Error(ctx, fiber.StatusInternalServerError, err, "failed to generate short url")
	}

	return response.Success(ctx, CreateShortURLResponse{
		ShortCode: code,
	})
}

func (h *URLHandler) ResolveURL(c *fiber.Ctx) error {
	code := c.Params("code")
	result, err := h.usecase.Resolve(c.Context(), code)
	if err != nil {
		return response.Error(c, http.StatusNotFound, err, "short code not found")
	}
	// ✅ Đọc header/IP vào biến trước khi vào goroutine

	// ✅ Launch goroutine an toàn
	// 👇 Collect meta safely
	meta := utils.ExtractClickMetaFromCtx(c, code)
	fmt.Println(meta)

	// 👇 Launch safe async tracking
	go h.usecase.TrackClick(context.Background(), meta)
	return response.Success(c, fiber.Map{
		"short_code": code,
		"long_url":   result.LongURL,
	})
	// Uncomment the next line to redirect instead of returning JSON
	// return c.Redirect(result.LongURL, http.StatusMovedPermanently)
}

// DELETE /api/shorten/:shortCode
func (h *URLHandler) DeleteShortURL(ctx *fiber.Ctx) error {
	code := ctx.Params("shortCode")
	if err := h.usecase.Delete(ctx.Context(), code); err != nil {
		return response.Error(ctx, fiber.StatusInternalServerError, err, "failed to delete")
	}
	return response.Success(ctx, fiber.Map{"message": "deleted"})
}
