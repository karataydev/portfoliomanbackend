package asset

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/karataydev/portfoliomanbackend/internal/database"
	"github.com/lib/pq"
)

type Repository struct {
	db *database.DBConnection
}

func NewRepository(db *database.DBConnection) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAssets() ([]SimpleAssetDTO, error) {
	query := `
        SELECT id, name, symbol
        FROM asset
    `
	var assets []SimpleAssetDTO
	err := r.db.Select(&assets, query)
	if err != nil {
		log.Errorf("Error fetching assets: %v", err)
		return nil, err
	}

	return assets, nil
}

func (r *Repository) GetAsset(assetId int64) (*Asset, error) {
	query := `
        SELECT *
        FROM asset
        WHERE id = $1
    `
	var asset Asset
	err := r.db.Select(&asset, query, assetId)
	if err != nil {
		log.Errorf("Error fetching asset: %v", err)
		return nil, err
	}

	return &asset, nil
}

func (r *Repository) GetAssetBySymbol(symbol string) (*Asset, error) {
	query := `
        SELECT *
        FROM asset
        WHERE symbol = $1
    `
	var asset Asset
	err := r.db.Get(&asset, query, symbol)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *Repository) GetAssetBySymbolList(symbols []string) ([]Asset, error) {
	query := `
		SELECT *
		FROM asset
		WHERE symbol = ANY($1)
	`
	var assets []Asset
	err := r.db.Select(&assets, query, pq.Array(symbols))
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *Repository) GetAssetQuoteAtTime(assetId int64, t time.Time) (*AssetQuote, error) {
	query := `
        SELECT *
        FROM asset_quote
        WHERE asset_id = $1 AND quote_time <= $2
        ORDER BY quote_time DESC
        LIMIT 1
    `
	var quote AssetQuote
	err := r.db.Get(&quote, query, assetId, t)
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func (r *Repository) GetAssetQuotesForPeriod(assetId int64, startTime, endTime time.Time) ([]AssetQuote, error) {
	query := `
        SELECT *
        FROM asset_quote
        WHERE asset_id = $1 AND quote_time BETWEEN $2 AND $3
        ORDER BY quote_time ASC
    `
	var quotes []AssetQuote
	err := r.db.Select(&quotes, query, assetId, startTime, endTime)
	if err != nil {
		return nil, err
	}
	return quotes, nil
}

func (r *Repository) SaveAssetQuote(assetQuote AssetQuote) error {
	query := `
			INSERT INTO asset_quote (asset_id, quote, quote_time)
			VALUES (:asset_id, :quote, :quote_time)
			ON CONFLICT (asset_id, quote_time)
			DO UPDATE SET quote = :quote
	`

	// Begin a transaction
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Will be ignored if the tx has been committed later

	// Insert or update the assetQuote
	_, err = tx.NamedExec(query, assetQuote)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

func (r *Repository) SearchAssets(searchSymbol string, limit int, offset int) ([]SimpleAssetDTO, int, error) {
	isAll := limit == -1

	query := `
        SELECT id, name, symbol
        FROM asset
        WHERE LOWER(symbol) LIKE LOWER($1) OR LOWER(name) LIKE LOWER($1)
        `
	if !isAll {
		query += `
			ORDER BY symbol
			LIMIT $2 OFFSET $3
		`
	}

	countQuery := `
        SELECT COUNT(*)
        FROM asset
        WHERE LOWER(symbol) LIKE LOWER($1) OR LOWER(name) LIKE LOWER($1)
    `

	searchTerm := "%" + strings.ToLower(searchSymbol) + "%"

	var assets []SimpleAssetDTO
	err := r.db.Select(&assets, query, searchTerm, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching assets: %v", err)
	}
	if isAll {
		return assets, len(assets), nil

	}
	var totalCount int
	err = r.db.Get(&totalCount, countQuery, searchTerm)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting assets: %v", err)
	}

	return assets, totalCount, nil
}
