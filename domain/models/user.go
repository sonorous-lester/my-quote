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
