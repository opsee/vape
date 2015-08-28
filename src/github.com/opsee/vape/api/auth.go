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

func init() {
	authRouter = router.Subrouter(AuthContext{}, "/authenticate")
	authRouter.Post("/password", (*AuthContext).CreateAuthPassword)
	authRouter.Get("/echo", (*AuthContext).Echo) // for testing
}

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

func (c *AuthContext) Echo(rw web.ResponseWriter, r *web.Request) {
	writeJson(rw, c.CurrentUser)
}
