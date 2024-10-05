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
        SELECT *
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

func (r *Repository) GetPortfolioByUser(userId int64) (*[]PortfolioDTO, error) {
	query := `
        SELECT *
        FROM portfolio
        WHERE user_id = $1
    `
	var portfolios []PortfolioDTO
	err := r.db.Select(&portfolios, query, userId)
	if err != nil {
		log.Errorf("Error fetching portfolio: %v", err)
		return nil, err
	}
	return &portfolios, nil
}

func (r *Repository) GetPortfolioBySymbol(symbol string) (*Portfolio, error) {
	query := `
        SELECT *
        FROM portfolio
        WHERE symbol = $1
    `
	var portfolio Portfolio
	err := r.db.Get(&portfolio, query, symbol)
	if err != nil {
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

func (r *Repository) GetPortfolioByUserIdWithAllocations(userId int64) ([]PortfolioDTO, error) {
	portfolios, err := r.GetPortfolioByUser(userId)
	if err != nil {
		return nil, err
	}

	var dtos []PortfolioDTO
	for _, val := range *portfolios {
		allocations, err := r.GetAllocations(val.Id)
		if err != nil {
			return nil, err
		}
		val.Allocations = allocations
		dtos = append(dtos, val)
	}

	return dtos, nil
}

func (r *Repository) GetFollowedPortfolios(userId int64) ([]PortfolioDTO, error) {
	query := `
        SELECT p.* FROM portfolio p
        JOIN portfolio_follow pf ON p.id = pf.portfolio_id
        WHERE pf.user_id = $1
    `
	var portfolios []PortfolioDTO
	err := r.db.Select(&portfolios, query, userId)
	if err != nil {
		log.Errorf("Error fetching portfolio: %v", err)
		return nil, err
	}
	return portfolios, nil
}

func (r *Repository) GetFollowedPortfoliosWithAllocations(userId int64) ([]PortfolioDTO, error) {
	portfolios, err := r.GetFollowedPortfolios(userId)
	if err != nil {
		return nil, err
	}

	var dtos []PortfolioDTO
	for _, val := range portfolios {
		allocations, err := r.GetAllocations(val.Id)
		if err != nil {
			return nil, err
		}
		val.Allocations = allocations
		dtos = append(dtos, val)
	}

	return dtos, nil
}

func (r *Repository) FollowPortfolio(userID, portfolioID int64) error {
	query := `
        INSERT INTO portfolio_follow (user_id, portfolio_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, portfolio_id) DO NOTHING
    `
	_, err := r.db.Exec(query, userID, portfolioID)
	return err
}

func (r *Repository) UnfollowPortfolio(userID, portfolioID int64) error {
	query := `
        DELETE FROM portfolio_follow
        WHERE user_id = $1 AND portfolio_id = $2
    `
	_, err := r.db.Exec(query, userID, portfolioID)
	return err
}

func (r *Repository) GetFollowerCount(portfolioID int64) (int, error) {
	query := `
        SELECT COUNT(*) FROM portfolio_follow
        WHERE portfolio_id = $1
    `
	var count int
	err := r.db.Get(&count, query, portfolioID)
	return count, err
}

func (r *Repository) IsFollowing(userID, portfolioID int64) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM portfolio_follow
            WHERE user_id = $1 AND portfolio_id = $2
        )
    `
	var isFollowing bool
	err := r.db.Get(&isFollowing, query, userID, portfolioID)
	return isFollowing, err
}
