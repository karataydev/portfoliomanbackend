package asset

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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

func (h *Handler) GetMarketOverview(c *fiber.Ctx) error {
	resp, err := h.service.GetMarketOverview()
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch market overview",
		})
	}
	return c.JSON(resp)
}

func (h *Handler) SearchAssets(c *fiber.Ctx) error {
    query := c.Query("q", "")
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    page, _ := strconv.Atoi(c.Query("page", "1"))

    if limit <= 0 {
        limit = 10
    }
    if page <= 0 {
        page = -1
    }

    offset := (page - 1) * limit

    assets, totalCount, err := h.service.SearchAssets(query, limit, offset)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to search assets",
        })
    }

    return c.JSON(fiber.Map{
        "assets":      assets,
        "total_count": totalCount,
        "page":        page,
        "limit":       limit,
        "total_pages": (totalCount + limit - 1) / limit,
    })
}
