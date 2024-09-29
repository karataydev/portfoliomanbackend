package user

import (
	"errors"
	"time"
)

var UserNotFoundErr error = errors.New("User not found")
var UserExistsErr error = errors.New("User exists")

type User struct {
	Id                int64     `db:"id" json:"id"`
	FirstName         string    `db:"first_name" json:"first_name"`
	LastName          string    `db:"last_name" json:"last_name"`
	Email             string    `db:"email" json:"email"`
	GoogleId          string    `db:"google_id" json:"google_id"`
	ProfilePictureUrl string    `db:"profile_picture_url" json:"profile_picture_url"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

type SignInUpResponse struct {
	User        *User  `json:"user"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	UserExisted bool   `json:"user_existed"`
}
