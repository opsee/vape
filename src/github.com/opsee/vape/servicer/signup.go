package servicer

import (
	"database/sql"
	"errors"
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

func CreateSignup(signupParams map[string]interface{}) (*model.Signup, error) {
	signup := model.NewSignup(signupParams)
	_, err := store.NamedExec("insert-signup", signup)
	return signup, err
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
	err := store.Select(signups, "list-signups", limit, offset)
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
	user := model.NewUser(signup.Name, signup.Email, password)
	user.Verified = true
	user.Active = true

	tx, err := store.Beginx()
	if err != nil {
		return nil, err
	}

	// need an org id for the user
	var orgId int
	if err = tx.Get(&orgId, "insert-new-org"); err != nil {
		tx.Rollback()
		return nil, err
	}
	user.OrgId = orgId

	if _, err := tx.Exec("claim-signup", signup.Id); err != nil {
		tx.Rollback()
		return nil, err
	}

	// need to pull out the generated user id, so use a query instead
	rows, err := tx.NamedQuery("insert-user", user)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for rows.Next() {
		if err = rows.StructScan(user); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}
