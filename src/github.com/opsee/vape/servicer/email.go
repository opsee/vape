package servicer

import (
	"github.com/opsee/vape/model"
)

func SendTemplatedEmail(userId int, template string, vars map[string]interface{}) (*model.User, error) {
	user, err := GetUser(userId)
	if err != nil {
		return nil, err
	}

	vars["name"] = user.Name
	vars["email"] = user.Email

	go func() {
		mailTemplatedMessage(user.Email, user.Name, template, vars)
	}()

	return user, nil
}
