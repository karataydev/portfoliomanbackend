package portfolio

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetPortfolio(c *fiber.Ctx) error {
	log.Info("in get by id")
	portfolioId, err := c.ParamsInt("portfolioId")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid portfolio ID"})
	}

	portfolio, err := h.service.GetPortfolio(int64(portfolioId))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Portfolio not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch portfolio"})
	}

	return c.JSON(portfolio)
}

func (h *Handler) GetPortfolioWithAllocations(c *fiber.Ctx) error {
	portfolioId, err := c.ParamsInt("portfolioId")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid portfolio ID"})
	}

	portfolio, err := h.service.GetPortfolioWithAllocations(int64(portfolioId))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Portfolio not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch portfolio"})
	}

	return c.JSON(portfolio)
}

func (h *Handler) AddTransactionToPortfolio(c *fiber.Ctx) error {
	var request AddTransactionRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := request.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	portfolio, err := h.service.AddTransactionToPortfolio(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(portfolio)
}

func (h *Handler) GetUserPortfolios(c *fiber.Ctx) error {
	log.Info("in get by user")
	userIdInterface := c.Locals("userId")
	userId, ok := userIdInterface.(int64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	portfolioListResponse, err := h.service.GetPortfolioListByUser(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch portfolios",
		})
	}

	return c.JSON(portfolioListResponse)
}


func (h *Handler) FollowPortfolio(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int64)
	portfolioID, err := c.ParamsInt("portfolioId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid portfolio ID"})
	}

	err = h.service.FollowPortfolio(userID, int64(portfolioID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to follow portfolio"})
	}

	return c.JSON(fiber.Map{"message": "Successfully followed portfolio"})
}

func (h *Handler) UnfollowPortfolio(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int64)
	portfolioID, err := c.ParamsInt("portfolioId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid portfolio ID"})
	}

	err = h.service.UnfollowPortfolio(userID, int64(portfolioID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to unfollow portfolio"})
	}

	return c.JSON(fiber.Map{"message": "Successfully unfollowed portfolio"})
}

func (h *Handler) GetFollowedPortfolios(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int64)

	portfolios, err := h.service.GetFollowedPortfolioList(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch followed portfolios"})
	}

	return c.JSON(portfolios)
}

func (h *Handler) GetFollowerCount(c *fiber.Ctx) error {
	portfolioID, err := c.ParamsInt("portfolioId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid portfolio ID"})
	}

	count, err := h.service.GetFollowerCount(int64(portfolioID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch follower count"})
	}

	return c.JSON(fiber.Map{"follower_count": count})
}

func (h *Handler) IsFollowing(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int64)
	portfolioID, err := c.ParamsInt("portfolioId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid portfolio ID"})
	}

	isFollowing, err := h.service.IsFollowing(userID, int64(portfolioID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check following status"})
	}

	return c.JSON(fiber.Map{"is_following": isFollowing})
}


func (h *Handler) CreatePortfolio(c *fiber.Ctx) error {
    var req CreatePortfolioRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    // Get the user ID from the context (set by the auth middleware)
    userId := c.Locals("userId").(int64)
    req.UserId = userId

    createdPortfolio, err := h.service.CreatePortfolio(req)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": fmt.Sprintf("Failed to create portfolio: %v", err),
        })
    }

    return c.JSON(createdPortfolio)
}