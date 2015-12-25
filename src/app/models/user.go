package models

import (
	"time"
)

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	APIToken     string `json:"api_token"`
	UpdatedAt    time.Time
	CreatedAt    time.Time `json:"created_at"`
	Subs         []Sub     `json:"subs"`
}
