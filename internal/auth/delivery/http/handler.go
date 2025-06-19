package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/huuloc2026/SwiftURL/internal/auth/usecase"
	"github.com/huuloc2026/SwiftURL/pkg/response"
)

type AuthHandler struct {
	uc usecase.AuthUsecase
}

func NewAuthHandler(uc usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err, "invalid request")
	}
	_, err := h.uc.Register(c.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err, err.Error())
	}
	return response.Success(c, fiber.Map{"message": "user registered successfully"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err, "invalid request")
	}
	token, err := h.uc.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, err, err.Error())
	}
	return response.Success(c, fiber.Map{"token": token})
}

func (h *AuthHandler) ForgetPassword(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err, "invalid request")
	}
	err := h.uc.ForgetPassword(c.Context(), req.Email)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err, err.Error())
	}
	return response.Success(c, fiber.Map{"message": "OTP sent successfully"})
}

func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err, "invalid request")
	}
	ok, err := h.uc.VerifyOTP(c.Context(), req.Email, req.OTP)
	if err != nil || !ok {
		return response.Error(c, fiber.StatusUnauthorized, err, "invalid OTP")
	}
	return response.Success(c, fiber.Map{"message": "OTP verified successfully"})
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	var req struct {
		Email       string `json:"email"`
		OTP         string `json:"otp"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err, "invalid request")
	}
	err := h.uc.ChangePassword(c.Context(), req.Email, req.OTP, req.NewPassword)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err, err.Error())
	}
	return response.Success(c, fiber.Map{"message": "password changed successfully"})
}
