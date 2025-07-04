package handler

import (
	"Go-Exercise/pkg/model"
	"Go-Exercise/pkg/service"

	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.authService.Register(req.Name, req.Email, req.Password, req.Age)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := h.authService.Login(req.Email, req.Password, req.DeviceInfo)
	if err != nil {
		if err == service.ErrUserNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid credentials",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(token)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req model.RefreshRequest
	var refreshToken string
	var deviceInfo string

	authHeader := c.Get("Authorization")
	if authHeader != "" {
		refreshToken = strings.TrimPrefix(authHeader, "Bearer ")
	}

	if refreshToken == "" {
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		refreshToken = req.RefreshToken
		deviceInfo = req.DeviceInfo
	} else {
		deviceInfo = c.Get("Device-Info")
	}

	if refreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "refresh token is required",
		})
	}

	token, err := h.authService.RefreshToken(refreshToken, deviceInfo)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(token)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req model.LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.authService.Logout(req.RefreshToken, req.LogoutAll); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
