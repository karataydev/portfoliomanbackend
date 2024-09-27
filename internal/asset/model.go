package asset

import "database/sql"

type Asset struct {
	Id          int64          `db:"id" json:"id"`
	Name        string         `db:"name" json:"name"`
	Symbol      string         `db:"symbol" json:"symbol"`
	Description sql.NullString `db:"description" json:"description"`
}


type SimpleAssetDTO struct {
	Id          int64          `db:"id" json:"id"`
	Name        string         `db:"name" json:"name"`
	Symbol      string         `db:"symbol" json:"symbol"`
}