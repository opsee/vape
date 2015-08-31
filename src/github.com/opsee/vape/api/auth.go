package api

import (
	"github.com/gocraft/web"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/store"
	"github.com/opsee/vape/token"
	"net/http"
	"time"
)

type AuthContext struct {
	*Context
}

const tokenExpHours = 1

var authRouter *web.Router

// @SubApi Authentication API [/authenticate]

func init() {
	authRouter = publicRouter.Subrouter(AuthContext{}, "/authenticate")
	authRouter.Post("/password", (*AuthContext).CreateAuthPassword)
	authRouter.Get("/echo", (*AuthContext).Echo) // for testing
}

// @Title authenticateFromPassword
// @Description Authenticates a user with email and password.
// @Accept  json
// @Param   email           query   string  true         "A user's email"
// @Param   password        body    string  true         "A user's password"
// @Success 200 {object}    interface                    "Response will be empty"
// @Failure 401 {object}    interface           	 "Response will be empty"
// @Router /authenticate/password [post]

func (c *AuthContext) CreateAuthPassword(rw web.ResponseWriter, r *web.Request) {
	postJson, err := readJson(r)
	if err != nil || postJson["email"] == nil || postJson["password"] == nil {
		c.Job.EventErr("create-auth", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	email := postJson["email"].(string)
	password := postJson["password"].(string)
	c.Job.EventKv("create-auth.enter", map[string]string{"email": email})

	user := new(model.User)
	err = store.Get(user, "user-by-email-and-active", email, true)
	if err != nil {
		c.Job.EventErr("get-user", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = user.Authenticate(password)
	if err != nil {
		c.Job.EventErr("authenticate-user", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := token.New(user, user.Email, time.Now(), time.Now().Add(time.Hour*tokenExpHours))
	tokenString, err := token.Marshal()
	if err != nil {
		c.Job.EventErr("token-marshal", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJson(rw, map[string]interface{}{
		"user":  user,
		"token": tokenString,
	})
}

// @Title echoSession
// @Description Echos a user session given an authentication token.
// @Accept  json
// @Param   Authorization   header   string  true         "The Bearer token"
// @Success 200 {object}    model.User
// @Router /authenticate/echo [get]

func (c *AuthContext) Echo(rw web.ResponseWriter, r *web.Request) {
	writeJson(rw, c.CurrentUser)
}
