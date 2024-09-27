package portfolio

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

func (h *Handler) GetPortfolio(c *fiber.Ctx) error {
    portfolioID, err := c.ParamsInt("portfolioID")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid portfolio ID"})
    }

    portfolio, err := h.service.GetPortfolio(int64(portfolioID))
    if err != nil {
        if err == sql.ErrNoRows {
            return c.Status(404).JSON(fiber.Map{"error": "Portfolio not found"})
        }
        return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch portfolio"})
    }

    return c.JSON(portfolio)
}

func (h *Handler) GetPortfolioWithAllocations(c *fiber.Ctx) error {
    portfolioID, err := c.ParamsInt("portfolioID")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid portfolio ID"})
    }

    portfolio, err := h.service.GetPortfolioWithAllocations(int64(portfolioID))
    if err != nil {
        if err == sql.ErrNoRows {
            return c.Status(404).JSON(fiber.Map{"error": "Portfolio not found"})
        }
        return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch portfolio"})
    }

    return c.JSON(portfolio)
}