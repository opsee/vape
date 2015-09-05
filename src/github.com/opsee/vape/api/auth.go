package api

import (
	"github.com/gocraft/web"
	"github.com/opsee/vape/model"
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
	authRouter.Get("/echo", (*AuthContext).Echo) // for testing
}

type UserTokenResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
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
	}

	err := c.RequestJson(&request)
	if err != nil {
		c.BadRequest("malformed request body", err)
		return
	}

	if request.Email == "" {
		c.BadRequest("missing email")
		return
	}

	if request.Password == "" {
		c.BadRequest("missing password")
		return
	}

	user := new(model.User)
	err = store.Get(user, "user-by-email-and-active", request.Email, true)
	if err != nil {
		c.Unauthorized("credentials do not match an existing user", err)
		return
	}

	err = user.Authenticate(request.Password)
	if err != nil {
		c.Unauthorized("credentials do not match an existing user", err)
		return
	}

	token, err := servicer.TokenUser(user)
	if err != nil {
		c.InternalServerError("internal server error", err)
		return
	}

	c.ResponseJson(map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

// @Title echoSession
// @Description Echos a user session given an authentication token.
// @Accept  json
// @Param   Authorization   header   string  true         "The Bearer token"
// @Success 200 {object}    model.User
// @Router /authenticate/echo [get]
func (c *AuthContext) Echo(rw web.ResponseWriter, r *web.Request) {
	if c.CurrentUser == nil {
		c.Unauthorized("a user token is required")
		return
	}

	c.ResponseJson(c.CurrentUser)
}
