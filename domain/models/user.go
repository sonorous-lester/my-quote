package models

import "time"

type UserModel struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	Token     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
