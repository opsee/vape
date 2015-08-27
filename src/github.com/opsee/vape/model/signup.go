package model

import (
	"time"
)

type Signup struct {
	Id           int       `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Claimed      bool      `json:"claimed"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
