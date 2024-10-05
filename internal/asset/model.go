package asset

import (
	"database/sql"
	"time"
)

type Asset struct {
	Id          int64          `db:"id" json:"id"`
	Name        string         `db:"name" json:"name"`
	Symbol      string         `db:"symbol" json:"symbol"`
	Description sql.NullString `db:"description" json:"description"`
}

type SimpleAssetDTO struct {
	Id     int64  `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Symbol string `db:"symbol" json:"symbol"`
}

type AssetQuote struct {
	Id        int64     `db:"id" json:"id"`
	AssetId   int64     `db:"asset_id" json:"asset_id"`
	Quote     float64   `db:"quote" json:"quote"`
	QuoteTime time.Time `db:"quote_time" json:"quote_time"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type AssetQuoteChanData struct {
	Symbol    string
	AssetId   int64
	Quote     float64
	QuoteTime time.Time
}

type MarketGrowthListResponse struct {
	Id int64 `json:"id"`

	Symbol string  `json:"symbol"`
	Name   string  `json:"name"`
	Change float64 `json:"change"`
	Amount float64 `json:"amount"`
}
