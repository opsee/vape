package servicer

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	ld "github.com/launchdarkly/go-client"
	"github.com/opsee/basic/schema"
	opsee_types "github.com/opsee/protobuf/opseeproto/types"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/store"
	"github.com/opsee/vaper"
	log "github.com/sirupsen/logrus"
	"github.com/snorecone/closeio-go"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	log.SetLevel(log.DebugLevel)
}
func VerifyUser(user *schema.User, token string) (bool, error) {
	if !VerifyToken(fmt.Sprint(user.Id), token) {
		return false, nil
	}

	user.Verified = true
	_, err := store.NamedExec("update-user", user)
	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateActiveUser(name, email, referrer string) (*schema.User, error) {
	signup, err := createSignup("", email, name, referrer, true, &opsee_types.Permission{Perm: model.AllUserPerms})
	if err != nil {
		return nil, err
	}

	user, err := ClaimSignup(signup.Id, VerificationToken(fmt.Sprint(signup.Id)), name, "", false)
	if err != nil {
		return nil, err
	}

	// send an email, create a lead and notify slack here!
	go func() {
		toke, err := TokenUser(user, 24*7*time.Hour)
		if err != nil {
			return
		}

		mergeVars := map[string]interface{}{
			"user_id":                 fmt.Sprint(user.Id),
			"user_verification_token": VerificationToken(fmt.Sprint(user.Id)),
			"user_auth_token":         toke,
		}
		mailTemplatedMessage(user.Email, "", "new-user", mergeVars)

		lead := &closeio.Lead{
			Name: user.Email,
			Contacts: []*closeio.Contact{
				{
					Name: user.Email,
					Emails: []*closeio.Email{
						{
							Type:  "work",
							Email: user.Email,
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
			"user_name":  user.Name,
			"user_email": user.Email,
		}

		// work around template shortcomings
		if referrer != "" {
			slackMap["referrer"] = referrer
		}

		notifySlack("new-signup", slackMap)
	}()

	return user, nil
}

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
		if user.Email != email {
			user.Verified = false
		}
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
	user.Perms.Name = "user"

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
	user.Perms.Name = "user"

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
	user.Perms.Name = "user"

	return user, nil
}

func UpdateUserPerms(user *schema.User, perms *opsee_types.Permission, duration time.Duration) (string, error) {
	user.Perms = perms
	_, err := store.NamedExec("update-user-perms", user)
	if err != nil {
		return "", err
	}

	return TokenUser(user, duration)
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
	user.Perms.Name = "user"

	if !user.Verified {
		SendVerification(user)
	}

	return TokenUser(user, duration)
}

func DeleteUser(id int) error {
	_, err := store.Exec("delete-user-by-id", id)
	return err
}

func TokenUser(user *schema.User, duration time.Duration) (string, error) {
	UpdateLDFlags(user)
	user.Perms.Name = "user"

	token := vaper.New(user, user.Email, time.Now(), time.Now().Add(duration))
	return token.Marshal()
}

func InviteTokenUser(user *schema.User, duration time.Duration) error {
	tokenString, err := TokenUser(user, duration)
	if err != nil {
		return err
	}

	go func() {
		mergeVars := map[string]interface{}{
			"user_id":     fmt.Sprint(user.Id),
			"user_token":  tokenString,
			"permissions": user.Perms.Permissions(),
			"name":        user.Name,
		}
		mailTemplatedMessage(user.Email, user.Name, "user-invite", mergeVars)
	}()
	return nil
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

func SendVerification(user *schema.User) {
	go func() {
		toke, err := TokenUser(user, 24*7*time.Hour)
		if err != nil {
			return
		}

		mergeVars := map[string]interface{}{
			"user_id":                 fmt.Sprint(user.Id),
			"user_verification_token": VerificationToken(fmt.Sprint(user.Id)),
			"user_auth_token":         toke,
		}
		mailTemplatedMessage(user.Email, "", "resend-verification", mergeVars)
	}()
}

// returns a launch darkly user from a user
func LDUser(user *schema.User) ld.User {
	idstr := fmt.Sprintf("%d", user.Id)
	custom := map[string]interface{}{"Admin": user.Admin}
	name := user.Name
	email := user.Email
	return ld.User{
		Key:    &idstr,
		Name:   &name,
		Email:  &email,
		Custom: &custom,
	}
}

// sets launchdarkly flags for a user
// TODO(dan) diff flags prior to setting them
func SetLDFlags(user *schema.User) error {
	if ldClient == nil {
		return fmt.Errorf("launch darkly client not initialized")
	}
	lduser := LDUser(user)
	if user.TeamFlags != nil {
		for _, flag := range user.TeamFlags.Permissions() {
			log.Debug("setting flag %s for user: %d", flag, user.Id)
			show_feature, err := ldClient.Toggle(flag, lduser, true)
			if err != nil {
				return nil
			}
			log.Debug("ld flag %s: %t", flag, show_feature)
		}
	}
	return nil
}

// Update launch darkly flags
func UpdateLDFlags(user *schema.User) error {
	lduser := LDUser(user)
	if ldClient == nil {
		return fmt.Errorf("launch darkly client not initialized")
	}
	m, err := ldClient.AllFlags(lduser)
	if err != nil {
		log.WithError(err).Debug("launch darkly client return error")
		return err
	}
	if err == nil && m == nil {
		return fmt.Errorf("launch darkly client offline")
	}

	log.Debug("fetching launchdarkly flags for user.Id %d", user.Id)
	if m != nil {
		var flags []string
		for k, _ := range m {
			log.Debugf("Flag %s", k)
			flags = append(flags, k)
		}
		teamFlags, err := opsee_types.NewPermissions("team_flags", flags...)
		log.Debug("got launch darkly flags %v for user.Id %d", flags, user.Id)
		if err != nil {
			return err
		}
		user.TeamFlags = teamFlags
		return nil
	}

	return err
}
