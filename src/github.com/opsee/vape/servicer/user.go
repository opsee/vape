package servicer

import (
	"database/sql"
	"fmt"
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

func GetUserEmail(email string) (*model.User, error) {
	user := new(model.User)
	err := store.Get(user, "user-by-email", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}

		return nil, err
	}

	return user, nil
}

func UpdateUser(user *model.User, email, name, password string, duration time.Duration) (string, error) {
	err := user.Merge(email, name, password)
	if err != nil {
		return "", err
	}

	_, err = store.NamedExec("update-user", user)
	if err != nil {
		return "", err
	}

	return TokenUser(user, duration)
}

func DeleteUser(id int) error {
	_, err := store.Exec("delete-user-by-id", id)
	return err
}

func TokenUser(user *model.User, duration time.Duration) (string, error) {
	token := token.New(user, user.Email, time.Now(), time.Now().Add(duration))
	return token.Marshal()
}

func EmailTokenUser(user *model.User, duration time.Duration, referer string) error {
	tokenString, err := TokenUser(user, duration)
	if err != nil {
		return err
	}

	// email that token
	go func() {
		mergeVars := map[string]string{
			"user_id":    fmt.Sprint(user.Id),
			"user_token": tokenString,
			"referer":    referer,
		}
		mailTemplatedMessage(user.Email, user.Name, "password-reset", mergeVars)
	}()

	return nil
}
