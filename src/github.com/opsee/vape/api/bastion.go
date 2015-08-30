package api

import (
	"github.com/gocraft/web"
	"github.com/opsee/vape/servicer"
	"net/http"
)

type BastionContext struct {
	*Context
}

var bastionRouter *web.Router

func init() {
	bastionRouter = privateRouter.Subrouter(BastionContext{}, "/bastions")
	bastionRouter.Post("/", (*BastionContext).Create)
	bastionRouter.Post("/authenticate", (*BastionContext).Authenticate)
}

func (c *BastionContext) Create(rw web.ResponseWriter, r *web.Request) {
	json, err := readJson(r)
	if err != nil {
		c.Job.EventErr("error.parse", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = mustPresent(json, "org_id"); err != nil {
		c.Job.EventErr("error.parse", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	orgId, ok := json["org_id"].(float64)
	if !ok {
		c.Job.Event("error.cast.float64")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	bastion, plaintext, err := servicer.CreateBastion(int(orgId))
	if err != nil {
		c.Job.EventErr("error.create", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJson(rw, map[string]interface{}{"id": bastion.Id, "password": plaintext, "org_id": orgId})
}

func (c *BastionContext) Authenticate(rw web.ResponseWriter, r *web.Request) {
	json, err := readJson(r)
	if err != nil {
		c.Job.EventErr("error.parse", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = mustPresent(json, "id", "password"); err != nil {
		c.Job.EventErr("error.parse", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = servicer.AuthenticateBastion(json["id"].(string), json["password"].(string)); err != nil {
		c.Job.EventErr("error.auth", err)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
}
