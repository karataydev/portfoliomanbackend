package portfolio

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/karataydev/portfoliomanbackend/internal/database"
)

type Repository struct {
	db *database.DBConnection
}

func NewRepository(db *database.DBConnection) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetPortfolio(portfolioId int64) (*PortfolioDTO, error) {
	query := `
        SELECT id, user_id, name, description, created_at, updated_at
        FROM portfolio
        WHERE id = $1
    `
	var portfolio PortfolioDTO
	err := r.db.Get(&portfolio, query, portfolioId)
	if err != nil {
		log.Errorf("Error fetching portfolio: %v", err)
		return nil, err
	}
	return &portfolio, nil
}

func (r *Repository) GetAllocations(portfolioId int64) ([]AllocationDTO, error) {
	query := `
        SELECT
            a.id,
            a.target_percentage,
            ast.id AS "asset.id",
            ast.name AS "asset.name",
            ast.symbol AS "asset.symbol"
        FROM allocation a
        JOIN asset ast ON a.asset_id = ast.id
        WHERE a.portfolio_id = $1
    `

	var allocations []AllocationDTO
	err := r.db.Select(&allocations, query, portfolioId)
	if err != nil {
		log.Errorf("Error fetching allocations: %v", err)
		return nil, err
	}

	return allocations, nil
}

func (r *Repository) GetPortfolioWithAllocations(portfolioId int64) (*PortfolioDTO, error) {
	portfolio, err := r.GetPortfolio(portfolioId)
	if err != nil {
		return nil, err
	}

	allocations, err := r.GetAllocations(portfolioId)
	if err != nil {
		return nil, err
	}

	portfolio.Allocations = allocations
	return portfolio, nil
}
