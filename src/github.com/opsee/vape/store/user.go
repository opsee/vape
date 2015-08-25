package store

import (
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"

	// "database/sql"
	// "github.com/jmoiron/sqlx"
)

type User struct {
	Id           int       `json:"id" token"id"`
	CustomerId   string    `json:"customer_id" token:"customer_id" db:"customer_id"`
	Email        string    `json:"email" token:"email"`
	Name         string    `json:"name" token:"name"`
	Verified     bool      `json:"verified" token:"verified"`
	Admin        bool      `json:"admin" token:"admin"`
	Active       bool      `json:"active" token:"active"`
	Onboard      bool      `json:"onboard" token:"onboard"`
	PasswordHash string    `json:"-" db:"password_hash"`       // not going in token
	CreatedAt    time.Time `json:"created_at" db:"created_at"` // not going in token
	UpdatedAt    time.Time `json:"created_at" db:"updated_at"` // not going in token
}

var queries = map[string]string{
	"by-email-and-active": "select * from logins where email = $1 and active = $2 limit 1",
}

func GetUser(queryKey string, args ...interface{}) (*User, error) {
	user := &User{}
	err := db.Get(user, queries[queryKey], args...)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func AuthenticateUser(email, password string) (*User, error) {
	user, err := GetUser("by-email-and-active", email, true)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
