package models

import "time"

type UserModel struct {
	ID        uint `gorm: "primaryKey"`
	Name      string
	Email     string
	Password  string
	Token     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
