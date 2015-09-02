package servicer

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/store"
)

var (
	SignupAlreadyClaimed = errors.New("signup already claimed")
	SignupInvalidToken   = errors.New("invalid token for signup")
)

func GetSignup(id int) (*model.Signup, error) {
	signup := new(model.Signup)
	err := store.Get(signup, "signup-by-id", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, RecordNotFound
		}

		return nil, err
	}

	return signup, nil
}

func CreateSignup(email, name string) (*model.Signup, error) {
	signup := model.NewSignup(email, name)
	err := store.NamedInsert("insert-signup", signup)
	if err != nil {
		return nil, err
	}
	return signup, err
}

func ActivateSignup(id int, referer string) error {
	signup, err := GetSignup(id)
	if err != nil {
		return err
	}

	// send an email here!
	go func() {
		mergeVars := map[string]string{
			"id":      fmt.Sprint(signup.Id),
			"token":   signup.Token(),
			"name":    signup.Name,
			"referer": referer,
		}
		mailTemplatedMessage(signup.Email, signup.Name, "activation", mergeVars)
	}()

	return nil
}

func ListSignups(perPage int, page int) ([]*model.Signup, error) {
	if perPage < 1 {
		perPage = 20
	}

	if page < 1 {
		page = 1
	}

	limit := perPage
	offset := (perPage * page) - perPage

	signups := []*model.Signup{}
	err := store.Select(&signups, "list-signups", limit, offset)
	if err != nil {
		return nil, err
	}

	return signups, nil
}

func ClaimSignup(signup *model.Signup, token, password string) (*model.User, error) {
	if signup.Validate(token) == false {
		return nil, SignupInvalidToken
	}

	// ok, pop that stuff in the user, and make sure they're verified
	user, err := model.NewUser(signup.Name, signup.Email, password)
	if err != nil {
		return nil, err
	}
	user.Verified = true
	user.Active = true

	tx, err := store.Beginx()
	if err != nil {
		return nil, err
	}

	// need an customer id for the user
	var customerId string
	if err = tx.Get(&customerId, "insert-new-customer"); err != nil {
		tx.Rollback()
		return nil, err
	}
	user.CustomerId = customerId

	if _, err := tx.Exec("claim-signup", signup.Id); err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.NamedInsert("insert-user", user)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}
