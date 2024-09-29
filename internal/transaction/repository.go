package transaction

import (
	"database/sql"

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

func (r *Repository) Get(allocationIds ...int64) ([]Transaction, error) {
	query := `
			SELECT *
			FROM transaction
			WHERE allocation_id = ANY($1)
		`
	var transactions []Transaction
	err := r.db.Select(&transactions, query, pq.Array(allocationIds))
	if err != nil {
		if err == sql.ErrNoRows {
			return []Transaction{}, nil
		}
		log.Errorf("Error fetching transactions: %v", err)
		return nil, err
	}
	return transactions, nil
}

func (r *Repository) Save(t *Transaction) (*Transaction, error) {
	query := `
			INSERT INTO transaction (allocation_id, side, quantity, price)
			VALUES (:allocation_id, :side, :quantity, :price)
			RETURNING id, created_at
		`
	rows, err := r.db.NamedQuery(query, t)
	if err != nil {
		log.Errorf("Error saving transaction: %v", err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&t.Id, &t.CreatedAt)
		if err != nil {
			log.Errorf("Error scanning returned transaction data: %v", err)
			return nil, err
		}
	}

	return t, nil
}
