package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SignUp(c *fiber.Ctx) error {
	var req struct {
		GoogleToken string `json:"google_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.service.SignUp(req.GoogleToken)
	if err != nil {
		log.Errorf("Failed to sign up: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to sign up",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Handler) SignIn(c *fiber.Ctx) error {
	var req struct {
		GoogleToken string `json:"google_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.service.SignIn(req.GoogleToken)
	if err != nil {
		log.Errorf("Failed to sign in: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to sign in",
		})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
