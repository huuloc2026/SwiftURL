package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/huuloc2026/SwiftURL/internal/user/usecase"
	"github.com/huuloc2026/SwiftURL/pkg/response"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(u usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: u}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err, "invalid input")
	}
	id, err := h.usecase.Register(c.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err, err.Error())
	}
	return response.Success(c, fiber.Map{"id": id})
}

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err, "invalid id")
	}
	user, err := h.usecase.GetByID(c.Context(), id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err, "user not found")
	}
	return response.Success(c, user)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err, "invalid id")
	}
	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err, "invalid input")
	}
	user, err := h.usecase.Update(c.Context(), id, req.Username, req.Email, req.Password)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err, "failed to update user")
	}
	return response.Success(c, user)
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err, "invalid id")
	}
	user, err := h.usecase.Delete(c.Context(), id)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err, "failed to delete user")
	}
	return response.Success(c, user)
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	users, err := h.usecase.List(c.Context(), limit, offset)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err, "failed to list users")
	}
	return response.Success(c, users)
}
