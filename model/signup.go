package model

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	opsee_types "github.com/opsee/protobuf/opseeproto/types"
)

var signupKey = []byte{142, 80, 107, 188, 92, 20, 197, 218, 205, 136, 179, 124, 29, 252, 213, 190}

// 0111b -- TODO(dan) AllPerms(permset) (uint64, error) in basic
const AllUserPerms = uint64(0x7)

type Signup struct {
	Id         int                     `json:"id"`
	Email      string                  `json:"email"`
	Name       string                  `json:"name"`
	Claimed    bool                    `json:"claimed"`
	Activated  bool                    `json:"activated"`
	Referrer   string                  `json:"referrer"`
	CreatedAt  time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time               `json:"updated_at" db:"updated_at"`
	CustomerId string                  `json:"customer_id" db:"customer_id"`
	Perms      *opsee_types.Permission `json:"perms" db:"perms"`
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
