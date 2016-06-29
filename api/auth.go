package api

import (
	"time"

	"github.com/gocraft/web"
	"github.com/opsee/basic/schema"
	"github.com/opsee/vape/servicer"
	"github.com/opsee/vape/store"
)

type AuthContext struct {
	*Context
}

var authRouter *web.Router

// @SubApi Authentication API [/authenticate]
func init() {
	authRouter = publicRouter.Subrouter(AuthContext{}, "/authenticate")
	authRouter.Post("/password", (*AuthContext).CreateAuthPassword)
	authRouter.Post("/token", (*AuthContext).CreateAuthToken)
	authRouter.Get("/echo", (*AuthContext).Echo) // for testing
	authRouter.Put("/refresh", (*AuthContext).Refresh)
}

type UserTokenResponse struct {
	Token        string       `json:"token"`
	User         *schema.User `json:"user"`
	IntercomHMAC string       `json:"intercom_hmac"`
}

// @Title authenticateFromPassword
// @Description Authenticates a user with email and password.
// @Accept  json
// @Param   email           body    string  true         "A user's email"
// @Param   password        body    string  true         "A user's password"
// @Success 200 {object}    UserTokenResponse
// @Failure 401 {object}    MessageResponse
// @Router /authenticate/password [post]
func (c *AuthContext) CreateAuthPassword(rw web.ResponseWriter, r *web.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		As       int    `json:"as"`
	}

	tokenExp := time.Hour * 12

	err := c.RequestJson(&request)
	if err != nil {
		c.BadRequest(Messages.BadRequest, err)
		return
	}

	if request.Email == "" {
		c.BadRequest(Messages.EmailRequired)
		return
	}

	if request.Password == "" {
		c.BadRequest(Messages.PasswordRequired)
		return
	}

	user := new(schema.User)
	err = store.Get(user, "user-by-email-and-active", request.Email, true)
	if err != nil {
		c.Unauthorized(Messages.CredentialsMismatch, err)
		return
	}

	err = servicer.AuthenticateUser(user, request.Password)
	if err != nil {
		c.Unauthorized(Messages.CredentialsMismatch, err)
		return
	}

	// here is an admin requesting to log in as someone else
	if user.Admin && request.As > 0 {
		adminId := user.Id
		user = nil

		user, err = servicer.GetUser(request.As)
		if err != nil {
			if err == servicer.UserNotFound {
				c.NotFound(Messages.UserNotFound)
			} else {
				c.InternalServerError(Messages.InternalServerError, err)
			}
			return
		}

		tokenExp = time.Minute * 15
		user.AdminId = adminId
	}

	token, err := servicer.TokenUser(user, tokenExp)
	if err != nil {
		c.InternalServerError(Messages.InternalServerError, err)
		return
	}

	hmac, err := servicer.HMACIntercomUser(user)
	if err != nil {
		c.InternalServerError(Messages.InternalServerError, err)
		return
	}

	c.ResponseJson(map[string]interface{}{
		"user":          user,
		"token":         token,
		"intercom_hmac": hmac,
	})
}

// @Title authenticateFromToken
// @Description Authenticates a user by emailing a Bearer token.
// @Accept  json
// @Param   email           body    string  true         "A user's email"
// @Success 200 {object}    MessageResponse
// @Failure 401 {object}    MessageResponse
// @Router /authenticate/token [post]
func (c *AuthContext) CreateAuthToken(rw web.ResponseWriter, r *web.Request) {
	var request struct {
		Email string `json:"email"`
	}

	err := c.RequestJson(&request)
	if err != nil {
		c.BadRequest(Messages.BadRequest, err)
		return
	}

	if request.Email == "" {
		c.BadRequest(Messages.EmailRequired)
		return
	}

	user := new(schema.User)
	err = store.Get(user, "user-by-email-and-active", request.Email, true)
	if err != nil {
		c.Unauthorized(Messages.UserNotFound)
		return
	}

	referer := r.Header.Get("Origin")
	err = servicer.EmailTokenUser(user, time.Hour, referer)
	if err != nil {
		c.InternalServerError(Messages.InternalServerError, err)
		return
	}

	c.ResponseJson(&MessageResponse{Message: Messages.Ok})
}

// @Title echoSession
// @Description Echos a user session given an authentication token.
// @Accept  json
// @Param   Authorization   header   string  true         "The Bearer token"
// @Success 200 {object}    schema.User
// @Router /authenticate/echo [get]
func (c *AuthContext) Echo(rw web.ResponseWriter, r *web.Request) {
	if c.CurrentUser == nil {
		c.Unauthorized(Messages.TokenRequired)
		return
	}

	c.ResponseJson(c.CurrentUser)
}

// @Title refreshSession
// @Description Refreshes a user session given an authentication token.
// @Accept  json
// @Param   Authorization   header   string  true         "The Bearer token"
// @Success 200 {object}    UserTokenResponse
// @Failure 401 {object}    MessageResponse
// @Router /authenticate/refresh [put]
func (c *AuthContext) Refresh(rw web.ResponseWriter, r *web.Request) {
	if c.CurrentUser == nil || c.CurrentUser.Status != "active" {
		c.Unauthorized(Messages.TokenRequired)
		return
	}

	err := store.Get(c.CurrentUser, "user-by-email-and-active", c.CurrentUser.Email, true)
	if err != nil {
		c.Unauthorized(Messages.TokenRequired)
		return
	}

	token, err := servicer.TokenUser(c.CurrentUser, time.Hour*12)
	if err != nil {
		c.InternalServerError(Messages.InternalServerError, err)
		return
	}

	c.ResponseJson(map[string]interface{}{
		"user":  c.CurrentUser,
		"token": token,
	})
}
