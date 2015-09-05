package api

import (
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

// @SubApi Signup API [/signups]
func init() {
	signupRouter = publicRouter.Subrouter(SignupContext{}, "/signups")
	signupRouter.Post("/", (*SignupContext).CreateSignup)
	signupRouter.Get("/", (*SignupContext).ListSignups)
	signupRouter.Get("/:id", (*SignupContext).GetSignup)
	signupRouter.Post("/:id/claim", (*SignupContext).ClaimSignup)
	signupRouter.Put("/:id/activate", (*SignupContext).ActivateSignup)
}

// @Title listSignups
// @Description List all signups.
// @Accept  json
// @Param   Authorization   header   string  true        "The Bearer token - an admin user token is required"
// @Param   per_page        query    int     false       "Pagination - number of records per page"
// @Param   page            query    int     false       "Pagination - which page"
// @Success 200 {array}     model.Signup                 ""
// @Failure 401 {object}    MessageResponse           	 "Response will be empty"
// @Router /signups [get]
func (c *SignupContext) ListSignups(rw web.ResponseWriter, r *web.Request) {
	if c.CurrentUser == nil || c.CurrentUser.Admin != true {
		c.Unauthorized("must be an administrator to access this resource")
		return
	}

	var request struct {
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
	}

	err := c.RequestJson(&request)
	if err != nil {
		c.BadRequest("malformed request", err)
		return
	}

	if request.PerPage <= 0 {
		request.PerPage = 20
	}

	if request.Page <= 0 {
		request.Page = 1
	}

	signups, err := servicer.ListSignups(request.PerPage, request.Page)
	if err != nil {
		c.InternalServerError("internal server error", err)
		return
	}

	c.ResponseJson(signups)
}

// @Title createSignup
// @Description Create a new signup.
// @Accept  json
// @Param   name             body   string  true       "The user's name"
// @Param   email            body   string  true       "The user's email"
// @Success 200 {object}     model.Signup              ""
// @Router /signups [post]
func (c *SignupContext) CreateSignup(rw web.ResponseWriter, r *web.Request) {
	var request struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	err := c.RequestJson(&request)
	if err != nil {
		c.BadRequest("malformed request", err)
		return
	}

	if request.Name == "" {
		c.BadRequest("missing name")
		return
	}

	if request.Email == "" {
		c.BadRequest("missing email")
		return
	}

	// anyone is authorized for this
	signup, err := servicer.CreateSignup(request.Email, request.Name)
	if err != nil {
		if err == servicer.SignupExists {
			c.Conflict("that email address has been taken")
			return
		}

		c.InternalServerError("internal server error", err)
		return
	}

	c.ResponseJson(signup)
}

// @Title activateSignup
// @Description Sends the activation email for a signup. Can be called multiple times to send multiple emails.
// @Accept  json
// @Param   id               path      int   true   "The signup's id"
// @Success 200 {object}     interface              "An object with the claim token used to verify the signup (sent in email)"
// @Router /signups/{id}/activate [put]
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

	referer := r.Header.Get("Referer")

	err = servicer.ActivateSignup(c.Signup.Id, referer)
	if err != nil {
		c.Job.EventErr("error.fetch", err)
		return
	}

	writeJson(rw, map[string]interface{}{"token": c.Signup.Token()})
}

// @Title getSignup
// @Description Get a single signup.
// @Accept  json
// @Param   Authorization   header   string  true        "The Bearer token - an admin user token is required"
// @Param   id              path     int     true       "The signup id"
// @Success 200 {object}    model.Signup                 ""
// @Failure 401 {object}    interface           	 "Response will be empty"
// @Router /signups/{id} [get]
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

// @Title claimSignup
// @Description Claim a signup and turn it into a user (usually from a url in an activation email).
// @Accept  json
// @Param   id              path   int     true       "The signup id"
// @Param   token           body   string  true       "The signup verification token"
// @Param   password        body   string  true       "The desired plaintext password for the new user"
// @Success 200 {object}    model.User                ""
// @Failure 401 {object}    interface                 "Response will be empty"
// @Failure 409 {object}    interface                 "Response will be empty"
// @Router /signups/{id}/claim [post]
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

	tokenString, err := servicer.TokenUser(user)
	if err != nil {
		c.Job.EventErr("error.token", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJson(rw, map[string]interface{}{
		"user":  user,
		"token": tokenString,
	})
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
