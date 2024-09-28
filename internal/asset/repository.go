package asset

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
