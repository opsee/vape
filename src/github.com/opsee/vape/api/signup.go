package api

import (
	"fmt"
	"github.com/gocraft/web"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/servicer"
	"net/http"
	"strconv"
)

type SignupContext struct {
	*Context
	Id     int
	Signup *model.Signup
}

var signupRouter *web.Router

func init() {
	signupRouter = publicRouter.Subrouter(SignupContext{}, "/signups")
	signupRouter.Post("/", (*SignupContext).CreateSignup)
	signupRouter.Get("/", (*SignupContext).ListSignups)
	signupRouter.Get("/:id", (*SignupContext).GetSignup)
	signupRouter.Post("/:id/claim", (*SignupContext).ClaimSignup)
	signupRouter.Put("/:id/activate", (*SignupContext).ActivateSignup)
}

func (c *SignupContext) ListSignups(rw web.ResponseWriter, r *web.Request) {
	if c.CurrentUser == nil || c.CurrentUser.Admin != true {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	perPage, err := strconv.Atoi(r.FormValue("per_page"))
	if err != nil {
		perPage = 20
	}

	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 1
	}

	c.Job.EventKv("list.params", map[string]string{"page": fmt.Sprint(page), "per_page": fmt.Sprint(perPage)})
	signups, err := servicer.ListSignups(perPage, page)
	if err != nil {
		c.Job.EventErr("error.select", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJson(rw, signups)
}

func (c *SignupContext) CreateSignup(rw web.ResponseWriter, r *web.Request) {
	json, err := readJson(r)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = mustPresent(json, "name", "email"); err != nil {
		c.Job.EventErr("error.parse", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	email, ok := json["email"].(string)
	if !ok {
		c.Job.EventErr("error.parse", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	name, ok := json["name"].(string)
	if !ok {
		c.Job.EventErr("error.parse", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// anyone is authorized for this
	signup, err := servicer.CreateSignup(email, name)
	if err != nil {
		c.Job.EventErr("error.create", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJson(rw, signup)
}

func (c *SignupContext) ActivateSignup(rw web.ResponseWriter, r *web.Request) {
	if c.CurrentUser == nil || c.CurrentUser.Admin != true {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := c.FetchSignup(rw, r)
	if err != nil {
		c.Job.EventErr("error.fetch", err)
		return
	}

	if c.Signup == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	err = servicer.ActivateSignup(c.Signup.Id)
	if err != nil {
		c.Job.EventErr("error.fetch", err)
		return
	}

	writeJson(rw, map[string]interface{}{"token": c.Signup.Token()})
}

func (c *SignupContext) GetSignup(rw web.ResponseWriter, r *web.Request) {
	if c.CurrentUser == nil || c.CurrentUser.Admin != true {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := c.FetchSignup(rw, r)
	if err != nil {
		c.Job.EventErr("error.fetch", err)
		return
	}

	if c.Signup == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	writeJson(rw, c.Signup)
}

func (c *SignupContext) ClaimSignup(rw web.ResponseWriter, r *web.Request) {
	json, err := readJson(r)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	token, ok := json["token"]
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	password, ok := json["password"]
	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.FetchSignup(rw, r)
	if err != nil {
		return
	}

	if c.Signup == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	user, err := servicer.ClaimSignup(c.Signup, token.(string), password.(string))
	if err != nil {
		c.Job.EventErr("error.claim", err)
		switch err {
		case servicer.SignupAlreadyClaimed:
			rw.WriteHeader(http.StatusConflict)
		case servicer.RecordNotFound:
			rw.WriteHeader(http.StatusNotFound)
		case servicer.SignupInvalidToken:
			rw.WriteHeader(http.StatusUnauthorized)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	writeJson(rw, user)
}

func (c *SignupContext) FetchSignup(rw web.ResponseWriter, r *web.Request) error {
	id, err := strconv.Atoi(r.PathParams["id"])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return err
	}

	c.Id = id
	signup, err := servicer.GetSignup(id)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return err
	}

	c.Signup = signup
	return nil
}
