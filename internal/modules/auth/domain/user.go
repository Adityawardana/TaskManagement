package domain

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterRequestObject struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TeamName string `json:"team_name"`
}

type LoginRequestObject struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
