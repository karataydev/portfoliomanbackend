package param

import (
	"github.com/karataydev/portfoliomanbackend/internal/database"
)

type Repository struct {
	db *database.DBConnection
}

func NewRepository(db *database.DBConnection) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(key string) (*Param, error) {
	param := &Param{}
	err := r.db.Get(param, "SELECT * FROM param WHERE key = $1", key)
	if err != nil {
		return nil, err
	}
	return param, nil
}

func (r *Repository) Save(param *Param) error {
	query := `
		INSERT INTO param (key, value)
		VALUES (:key, :value)
		ON CONFLICT (key) DO UPDATE
		SET value = :value
		RETURNING key`
	rows, err := r.db.NamedQuery(query, param)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&param.Key)
	}
	return nil
}
