package portfolio

import (
	"database/sql"
	"time"

	"github.com/karataydev/portfoliomanbackend/internal/asset"
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

type AllocationDTO struct {
	Id                int64                `db:"id" json:"id"`
	Asset             asset.SimpleAssetDTO `json:"asset"`
	TargetPercentage  float64              `db:"target_percentage" json:"target_percentage"`
	Amount            float64              `db:"-" json:"amount"`
	CurrentPercentage float64              `db:"-" json:"current_percentage"`
}

type PortfolioDTO struct {
	Portfolio
	Allocations []AllocationDTO `json:"allocations"`
}
