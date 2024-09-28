package asset

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

func (h *Handler) GetAssets(c *fiber.Ctx) error {
	assets, err := h.service.GetAssets()
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Assets not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch assets"})
	}
	return c.JSON(assets)
}

func (h *Handler) GetAsset(c *fiber.Ctx) error {
	assetId, err := c.ParamsInt("assetId")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Asset ID"})
	}

	asset, err := h.service.GetAsset(int64(assetId))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Asset not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch asset"})
	}

	return c.JSON(asset)
}
