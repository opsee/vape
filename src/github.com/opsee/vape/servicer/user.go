package servicer

import (
	"database/sql"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/store"
	"github.com/opsee/vape/token"
	"time"
)

func GetUser(id int) (*model.User, error) {
	user := new(model.User)
	err := store.Get(user, "user-by-id", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}

		return nil, err
	}

	return user, nil
}

func UpdateUser(user *model.User, email, name, password string) error {
	err := user.Merge(email, name, password)
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

func TokenUser(user *model.User) (string, error) {
	token := token.New(user, user.Email, time.Now(), time.Now().Add(time.Hour*token.ExpHours))
	return token.Marshal()
}
