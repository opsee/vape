package model

import (
        "time"
)

type Bastion struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
