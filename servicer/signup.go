package servicer

import (
	"database/sql"
	"fmt"
	"github.com/opsee/basic/schema"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/store"
	"github.com/snorecone/closeio-go"
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

func DeleteSignup(id int) error {
	_, err := store.Exec("delete-signup-by-id", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return err
	}

	return nil
}

func CreateActiveSignup(email, name, referrer string) (*model.Signup, error) {
	signup, err := createSignup(email, name, referrer, true)
	if err != nil {
		return nil, err
	}

	// send an email, create a lead and notify slack here!
	go func() {
		mergeVars := map[string]interface{}{
			"signup_id":    fmt.Sprint(signup.Id),
			"signup_token": signup.Token(),
			"name":         signup.Name,
		}
		mailTemplatedMessage(signup.Email, signup.Name, "instant-approval", mergeVars)

		lead := &closeio.Lead{
			Name: signup.Email,
			Contacts: []*closeio.Contact{
				{
					Name: signup.Email,
					Emails: []*closeio.Email{
						{
							Type:  "work",
							Email: signup.Email,
						},
					},
				},
			},
		}

		if referrer != "" {
			lead.Custom = map[string]string{
				"referrer": referrer,
			}
		}

		createLead(lead)

		slackMap := map[string]interface{}{
			"user_name":  signup.Name,
			"user_email": signup.Email,
		}

		// work around template shortcomings
		if referrer != "" {
			slackMap["referrer"] = signup.Referrer
		}

		notifySlack("new-signup", slackMap)
	}()

	return signup, err
}

func createSignup(email, name, referrer string, activated bool) (*model.Signup, error) {
	existingSignup := new(model.Signup)
	err := store.Get(existingSignup, "signup-by-email", email)
	if err == nil {
		return nil, SignupExists
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	signup := &model.Signup{
		Email:     email,
		Name:      name,
		Referrer:  referrer,
		Activated: activated,
	}

	err = store.NamedInsert("insert-signup", signup)
	if err != nil {
		return nil, err
	}

	return signup, err
}

func ActivateSignup(id int) (*model.Signup, error) {
	signup, err := GetSignup(id)
	if err != nil {
		return nil, err
	}

	_, err = store.Exec("activate-signup", signup.Id)
	if err != nil {
		return nil, err
	}

	// send an email here!
	go func() {
		mergeVars := map[string]interface{}{
			"signup_id":    fmt.Sprint(signup.Id),
			"signup_token": signup.Token(),
			"name":         signup.Name,
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

func ClaimSignup(id int, token, name, password string, invite bool) (*schema.User, error) {
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
	user, err := NewUser(name, signup.Email, password)
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

	if invite {
		go inviteSlack(user.Name, user.Email)
	}

	return user, nil
}
