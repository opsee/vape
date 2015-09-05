package servicer

import (
	"errors"
)

var (
	RecordNotFound       = errors.New("record not found")
	CustomerNotFound     = errors.New("customer not found")
	SignupAlreadyClaimed = errors.New("signup already claimed")
	SignupInvalidToken   = errors.New("invalid token for signup")
)
