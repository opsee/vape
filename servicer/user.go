package servicer

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/opsee/basic/schema"
	"github.com/opsee/vape/store"
	"github.com/opsee/vaper"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func NewUser(name, email, password string) (*schema.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	return &schema.User{
		Email:        email,
		Name:         name,
		PasswordHash: string(passwordHash),
	}, nil
}

func AuthenticateUser(user *schema.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}

func MergeUser(user *schema.User, email, name, password string) error {
	if email != "" {
		user.Email = email
	}

	if name != "" {
		user.Name = name
	}

	if password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			return err
		}
		user.PasswordHash = string(passwordHash)
	}

	return nil
}

type PaginatedUsers struct {
	Page    int
	PerPage int
	Total   int
	Users   []*schema.User
}

func ListUsers(perPage int, page int) (*PaginatedUsers, error) {
	if perPage < 1 {
		perPage = 20
	}

	if page < 1 {
		page = 1
	}

	limit := perPage
	offset := (perPage * page) - perPage

	users := []*schema.User{}
	err := store.Select(&users, "list-users", limit, offset)
	var total int

	tx, err := store.Beginx()
	if err != nil {
		return nil, err
	}

	defer tx.Commit()

	err = tx.Select(&users, "list-users", limit, offset)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Get(&total, "total-users")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &PaginatedUsers{
		Page:    page,
		PerPage: perPage,
		Total:   total,
		Users:   users,
	}, nil
}

func HMACIntercomUser(user *schema.User) (string, error) {
	if intercomKey == nil {
		return "", nil
	}

	hashWriter := hmac.New(sha256.New, intercomKey)
	_, err := hashWriter.Write([]byte(user.Email))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hashWriter.Sum(nil)), nil
}

func GetUser(id int) (*schema.User, error) {
	user := new(schema.User)
	err := store.Get(user, "user-by-id", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}

		return nil, err
	}

	return user, nil
}

func GetUserCustID(id string) (*schema.User, error) {
	user := new(schema.User)
	err := store.Get(user, "user-by-cust-id", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}

		return nil, err
	}

	return user, nil
}

func GetUserEmail(email string) (*schema.User, error) {
	user := new(schema.User)
	err := store.Get(user, "user-by-email", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, UserNotFound
		}

		return nil, err
	}

	return user, nil
}

func UpdateUser(user *schema.User, email, name, password string, duration time.Duration) (string, error) {
	err := MergeUser(user, email, name, password)
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

func TokenUser(user *schema.User, duration time.Duration) (string, error) {
	token := vaper.New(user, user.Email, time.Now(), time.Now().Add(duration))
	return token.Marshal()
}

func EmailTokenUser(user *schema.User, duration time.Duration, referer string) error {
	tokenString, err := TokenUser(user, duration)
	if err != nil {
		return err
	}

	// email that token
	go func() {
		mergeVars := map[string]interface{}{
			"user_id":    fmt.Sprint(user.Id),
			"user_token": tokenString,
			"referer":    referer,
			"name":       user.Name,
		}
		mailTemplatedMessage(user.Email, user.Name, "password-reset", mergeVars)
	}()

	return nil
}

func GetUserData(id int) ([]byte, error) {
	var userdata struct {
		Data []byte
	}

	err := store.Get(&userdata, "userdata-by-id", id)
	return userdata.Data, err
}

func UpdateUserData(id int, data []byte) ([]byte, error) {
	var userdata struct {
		Data []byte
	}

	err := store.Get(&userdata, "merge-userdata", id, data)
	return userdata.Data, err
}