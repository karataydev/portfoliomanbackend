package portfolio

import (
	"database/sql"
	"errors"
	"time"

	"github.com/karataydev/portfoliomanbackend/internal/asset"
	"github.com/karataydev/portfoliomanbackend/internal/transaction"
)

type Portfolio struct {
	Id          int64          `db:"id" json:"id"`
	Symbol      string         `db:"symbol" json:"symbol"`
	UserId      int64          `db:"user_id" json:"user_id"`
	Name        string         `db:"name" json:"name"`
	Description sql.NullString `db:"description" json:"description"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
}

type Allocation struct {
	Id               int64     `db:"id" json:"id"`
	PortfolioId      int64     `db:"portfolio_id" json:"portfolio_id"`
	AssetId          int64     `db:"asset_id" json:"asset_id"`
	TargetPercentage float64   `db:"target_percentage" json:"target_percentage"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type PortfolioFollow struct {
	Id          int64     `db:"id" json:"id"`
	UserId      int64     `db:"user_id" json:"user_id"`
	PortfolioId int64     `db:"portfolio_id" json:"portfolio_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type AllocationDTO struct {
	Id                int64                `db:"id" json:"id"`
	Asset             asset.SimpleAssetDTO `json:"asset"`
	TargetPercentage  float64              `db:"target_percentage" json:"target_percentage"`
	Amount            float64              `db:"-" json:"amount"`
	CurrentPercentage float64              `db:"-" json:"current_percentage"`
	UnrealizedPL      float64              `db:"-" json:"unrealized_pl"`
}

type PortfolioDTO struct {
	Portfolio
	Allocations []AllocationDTO `json:"allocations"`
}

type AddTransactionRequest struct {
	PortfolioId int64                 `json:"portfolio_id"`
	Symbol      string                `json:"symbol"`
	Quantity    float64               `json:"quantity"`
	AvgPrice    float64               `json:"avg_price"`
	Side        transaction.OrderSide `json:"side"`
}

func (r *AddTransactionRequest) validate() error {
	if r.PortfolioId <= 0 {
		return errors.New("invalid portfolio Id")
	}
	if r.Symbol == "" {
		return errors.New("symbol is required")
	}
	if r.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if r.AvgPrice <= 0 {
		return errors.New("average price must be positive")
	}
	return nil
}

type PortfolioListResponse struct {
	Id int64 `json:"id"`

	Symbol string  `json:"symbol"`
	Name   string  `json:"name"`
	Change float64 `json:"change"`
	Owner  string  `json:"owner"`
	Amount float64 `json:"amount"`
}

type CreatePortfolioRequest struct {
	UserId      int64               `json:"user_id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Allocations []AllocationRequest `json:"allocations"`
}

type AllocationRequest struct {
	AssetId          int64   `json:"asset_id"`
	TargetPercentage float64 `json:"target_percentage"`
}
