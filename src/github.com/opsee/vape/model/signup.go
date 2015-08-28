package model

import (
	"time"
        "fmt"
        "crypto/sha256"
        "crypto/hmac"
        "encoding/base64"
)

var signupKey = []byte{142, 80, 107, 188, 92, 20, 197, 218, 205, 136, 179, 124, 29, 252, 213, 190}

type Signup struct {
	Id           int       `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Claimed      bool      `json:"claimed"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func NewSignup(params map[string]interface{}) *Signup {
        signup := &Signup{}

        email, ok := params["email"]
        if ok {
                signup.Email = email.(string)
        }

        name, ok := params["name"]
        if ok {
                signup.Name = name.(string)
        }

        return signup
}


func (s *Signup) Token() string {
	return base64.StdEncoding.EncodeToString(s.token())
}

func (s *Signup) Validate(token string) bool {
        tok, err := base64.StdEncoding.DecodeString(token)
        if err != nil {
                return false
        }

	return hmac.Equal(s.token(), tok)
}

func (s *Signup) token() []byte {
        id := fmt.Sprintf("%s", s.Id)
        mac := hmac.New(sha256.New, signupKey)
	mac.Write([]byte(id))
        return mac.Sum(nil)
}
