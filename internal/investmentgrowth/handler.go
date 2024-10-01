package investmentgrowth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CalculateInvestmentGrowth(c *fiber.Ctx) error {
	symbol := c.Params("symbol")
	if symbol == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Symbol is required"})
	}

	growth, err := h.service.CalculateInvestmentGrowth(symbol)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": fmt.Sprintf("Failed to calculate growth for symbol %s: %v", symbol, err)})
	}

	return c.JSON(growth)
}
