package servicer

import (
	"database/sql"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/store"
)

func GetUser(id int) (*model.User, error) {
	user := new(model.User)
	err := store.Get(user, "user-by-id", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, RecordNotFound
		}

		return nil, err
	}

	return user, nil
}

func UpdateUser(user *model.User, newUserParams map[string]interface{}) error {
	err := user.Merge(newUserParams)
	if err != nil {
		return err
	}

	_, err = store.Exec("update-user", user)
	return err
}

func DeleteUser(id int) error {
	_, err := store.Exec("delete-user-by-id", id)
	return err
}
