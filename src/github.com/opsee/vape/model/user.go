package model

import (
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Id           int       `json:"id" token"id"`
	OrgId        int       `json:"org_id" token:"org_id" db:"org_id"`
	Email        string    `json:"email" token:"email"`
	Name         string    `json:"name" token:"name"`
	Verified     bool      `json:"verified" token:"verified"`
	Admin        bool      `json:"admin" token:"admin"`
	Active       bool      `json:"active" token:"active"`
	PasswordHash string    `json:"-" db:"password_hash"`       // not going in token
	CreatedAt    time.Time `json:"created_at" db:"created_at"` // not going in token
	UpdatedAt    time.Time `json:"created_at" db:"updated_at"` // not going in token
}

func (user *User) Authenticate(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}
