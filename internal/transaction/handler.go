package transaction

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Get(c *fiber.Ctx) error {
	portfolioId := c.QueryInt("allocationId")
	if portfolioId == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid portfolio ID"})
	}

	transactions, err := h.service.Get(int64(portfolioId))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Transactions not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch transactions"})
	}

	return c.JSON(fiber.Map{"transactions": transactions})
}

func (h *Handler) Save(c *fiber.Ctx) error {
	var transaction Transaction
	if err := c.BodyParser(&transaction); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	savedTransaction, err := h.service.Save(&transaction)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save transaction"})
	}

	return c.JSON(savedTransaction)
}
