package servicer

import (
	"database/sql"
	"fmt"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/store"
)

func GetSignup(id int) (*model.Signup, error) {
	signup := new(model.Signup)
	err := store.Get(signup, "signup-by-id", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, SignupNotFound
		}

		return nil, err
	}

	return signup, nil
}

func CreateSignup(email, name string) (*model.Signup, error) {
	existingSignup := new(model.Signup)
	err := store.Get(existingSignup, "signup-by-email", email)
	if err == nil {
		return nil, SignupExists
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	signup := model.NewSignup(email, name)
	err = store.NamedInsert("insert-signup", signup)
	if err != nil {
		return nil, err
	}

	// send an email here!
	go func() {
		mergeVars := map[string]string{}
		mailTemplatedMessage(signup.Email, signup.Name, "signup-confirmation", mergeVars)
	}()

	return signup, err
}

func ActivateSignup(id int, referer string) (*model.Signup, error) {
	signup, err := GetSignup(id)
	if err != nil {
		return nil, err
	}

	// send an email here!
	go func() {
		mergeVars := map[string]string{
			"id":      fmt.Sprint(signup.Id),
			"token":   signup.Token(),
			"name":    signup.Name,
			"referer": referer,
		}
		mailTemplatedMessage(signup.Email, signup.Name, "beta-approval", mergeVars)
	}()

	return signup, nil
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

func ClaimSignup(id int, token, password string) (*model.User, error) {
	signup, err := GetSignup(id)
	if err != nil {
		return nil, err
	}

	if signup.Validate(token) == false {
		return nil, SignupInvalidToken
	}

	// make sure user hasn't been claimed
	if signup.Claimed {
		return nil, SignupAlreadyClaimed
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
