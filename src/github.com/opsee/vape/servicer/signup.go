package servicer

import (
        "database/sql"
        "github.com/opsee/vape/store"
        "github.com/opsee/vape/model"
)

func GetSignup(id int) (*model.Signup, error) {
        signup := new(model.Signup)
        err := store.Get(signup, "signup-by-id", id)
        if err != nil {
                if err == sql.ErrNoRows {
                        return nil, nil
                }

                return nil, err
        }

        return signup, nil
}

func CreateSignup(signupParams map[string]interface{}) (*model.Signup, error) {
        signup := model.NewSignup(signupParams)
        _, err = store.NamedExec("insert-signup", signup)
        return err
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
