package api

import (
	"github.com/gocraft/web"
	"net/http"
        "github.com/opsee/vape/store"
)

type AuthContext struct {
	*Context
}

var authRouter *web.Router

func init() {
	authRouter = router.Subrouter(AuthContext{}, "/")
	authRouter.Post("/login", (*AuthContext).CreateAuth)
}

func (c *AuthContext) CreateAuth(rw web.ResponseWriter, r *web.Request) {
        postJson, err := readJson(r)
        if err != nil {
                c.Job.EventErr("create-auth", err)
                rw.WriteHeader(http.StatusBadRequest)
                return
        }

        email := postJson["email"].(string)
        password := postJson["password"].(string)
        c.Job.EventKv("create-auth.enter", map[string]string{"email": email})

        user, err := store.AuthenticateUser(email, password)
        if err != nil {
                c.Job.EventErr("authenticate-user", err)
                rw.WriteHeader(http.StatusUnauthorized)
                return
        }

        token, err := user.MarshalJwe()
        if err != nil {
                c.Job.EventErr("marshal-jwe", err)
                rw.WriteHeader(http.StatusInternalServerError)
                return
        }

        writeJson(rw, map[string]interface{}{
                "user": user,
                "token": token,
        })
}
