package user

import (
	"database/sql"

	"github.com/gofiber/fiber/v2/log"
	"github.com/karataydev/portfoliomanbackend/internal/database"
)

type Repository struct {
	db *database.DBConnection
}

func NewRepository(db *database.DBConnection) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(userId int64) (*User, error) {
	query := `
			SELECT *
			FROM users
			WHERE id = ($1)
		`
	var user User
	err := r.db.Get(&user, query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFoundErr
		}
		log.Errorf("Error fetching user: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetByEmail(email string) (*User, error) {
	query := `
			SELECT *
			FROM users
			WHERE email = ($1)
		`
	var user User
	err := r.db.Get(&user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFoundErr
		}
		log.Errorf("Error fetching user: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *Repository) Save(user *User) (*User, error) {
	query := `
		INSERT INTO users (first_name, last_name, email)
		VALUES (:first_name, :last_name, :email)
		RETURNING id, created_at, updated_at
	`
	rows, err := r.db.NamedQuery(query, user)
	if err != nil {
		log.Errorf("Error saving user: %v", err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Errorf("Error scanning returned user data: %v", err)
			return nil, err
		}
	}

	return user, nil
}
