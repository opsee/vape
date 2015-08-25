package store

import (
        _ "github.com/lib/pq"
        "time"
        "encoding/json"
        "golang.org/x/crypto/bcrypt"
        "github.com/dvsekhvalnov/jose2go"
        // "database/sql"
        // "github.com/jmoiron/sqlx"
)

const tokenExpHours = 72

type User struct {
        Id int `json:"id"`
        CustomerId string `json:"customer_id" db:"customer_id"`
        Email string `json:"email"`
        Name string `json:"name"`
        Verified bool `json:"verified"`
        Admin bool `json:"admin"`
        Active bool `json:"active"`
        Onboard bool `json:"onboard"`
        PasswordHash string `json:"-" db:"password_hash"`
        CreatedAt time.Time `json:"created_at" db:"created_at"`
        UpdatedAt time.Time `json:"created_at" db:"updated_at"`
}

func AuthenticateUser(email, password string) (*User, error) {
        user := &User{}
        err := db.Get(user, "select * from logins where email = $1 and active = true limit 1", email)
        if err != nil {
                return nil, err
        }

        err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
        if err != nil {
                return nil, err
        }

        return user, nil
}

func (user *User) MarshalJwe() (string, error) {
        payload := map[string]interface{}{
                "exp": time.Now().Add(time.Hour * tokenExpHours).Unix(),
                "id": user.Id,
                "email": user.Email,
                "customer_id": user.CustomerId,
                "admin": user.Admin,
                "verified": user.Verified,
                "active": user.Active,
                "name": user.Name,
                "onboard": user.Onboard,
                "created_at": user.CreatedAt,
                "updated_at": user.UpdatedAt,
        }

        json, err := json.Marshal(payload)
        if err != nil {
                return "", err
        }

        return jose.Encrypt(string(json), jose.A128GCMKW, jose.A128GCM, vapeKey)
}
