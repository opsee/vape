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

	if err = mustPresent(json, "customer_id"); err != nil {
		c.Job.EventErr("error.parse", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	customerId, ok := json["customer_id"].(string)
	if !ok {
		c.Job.Event("error.cast.string")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	bastion, plaintext, err := servicer.CreateBastion(customerId)
	if err != nil {
		c.Job.EventErr("error.create", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJson(rw, map[string]interface{}{"id": bastion.Id, "password": plaintext, "customer_id": customerId})
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
