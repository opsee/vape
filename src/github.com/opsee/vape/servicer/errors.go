package servicer

import (
	"errors"
)

var (
	RecordNotFound       = errors.New("record not found")
	CustomerNotFound     = errors.New("customer not found")
	SignupNotFound       = errors.New("customer not found")
	SignupExists         = errors.New("signup exists with that email")
	SignupAlreadyClaimed = errors.New("signup already claimed")
	SignupInvalidToken   = errors.New("invalid token for signup")
)
