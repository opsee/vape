package api

import (
	"github.com/gocraft/web"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/servicer"
	"strconv"
	"time"
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
// @Failure 401 {object}    MessageResponse           	 ""
// @Router /signups [get]
func (c *SignupContext) ListSignups(rw web.ResponseWriter, r *web.Request) {
	if c.CurrentUser == nil || c.CurrentUser.Admin != true {
		c.Unauthorized(Messages.AdminRequired)
		return
	}

	var request struct {
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
	}

	// ignore errors since all params are optional
	c.RequestJson(&request)

	if request.PerPage <= 0 {
		request.PerPage = 20
	}

	if request.Page <= 0 {
		request.Page = 1
	}

	signups, err := servicer.ListSignups(request.PerPage, request.Page)
	if err != nil {
		c.InternalServerError(Messages.InternalServerError, err)
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
// @Failure 409 {object}     MessageResponse           "Email was already used to sign up"
// @Router /signups [post]
func (c *SignupContext) CreateSignup(rw web.ResponseWriter, r *web.Request) {
	var request struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	err := c.RequestJson(&request)
	if err != nil {
		c.BadRequest(Messages.BadRequest, err)
		return
	}

	if request.Name == "" {
		c.BadRequest(Messages.NameRequired)
		return
	}

	if request.Email == "" {
		c.BadRequest(Messages.EmailRequired)
		return
	}

	// anyone is authorized for this
	signup, err := servicer.CreateSignup(request.Email, request.Name)
	if err != nil {
		if err == servicer.SignupExists {
			c.Conflict(Messages.EmailConflict)
			return
		}

		c.InternalServerError(Messages.InternalServerError, err)
		return
	}

	c.ResponseJson(signup)
}

type SignupActivationResponse struct {
	Token string `json:"token"`
}

// @Title activateSignup
// @Description Sends the activation email for a signup. Can be called multiple times to send multiple emails.
// @Accept  json
// @Param   id               path      int   true       "The signup's id"
// @Success 200 {object}     SignupActivationResponse   "An object with the claim token used to verify the signup (sent in email)"
// @Failure 401 {object}     MessageResponse            ""
// @Router /signups/{id}/activate [put]
func (c *SignupContext) ActivateSignup(rw web.ResponseWriter, r *web.Request) {
	if c.CurrentUser == nil || c.CurrentUser.Admin != true {
		c.Unauthorized(Messages.AdminRequired)
		return
	}

	id, err := strconv.Atoi(r.PathParams["id"])
	if err != nil {
		c.BadRequest(Messages.IdRequired)
		return
	}

	// sending referer should be temporary for developing email templates
	referer := r.Header.Get("Origin")
	signup, err := servicer.ActivateSignup(id, referer)
	if err != nil {
		if err == servicer.SignupNotFound {
			c.NotFound(Messages.SignupNotFound)
		} else {
			c.InternalServerError(Messages.InternalServerError, err)
		}

		return
	}

	c.ResponseJson(&SignupActivationResponse{Token: signup.Token()})
}

// @Title getSignup
// @Description Get a single signup.
// @Accept  json
// @Param   Authorization   header   string  true        "The Bearer token - an admin user token is required"
// @Param   id              path     int     true       "The signup id"
// @Success 200 {object}    model.Signup                 ""
// @Failure 401 {object}    MessageResponse           	 ""
// @Router /signups/{id} [get]
func (c *SignupContext) GetSignup(rw web.ResponseWriter, r *web.Request) {
	if c.CurrentUser == nil || c.CurrentUser.Admin != true {
		c.Unauthorized(Messages.AdminRequired)
		return
	}

	id, err := strconv.Atoi(r.PathParams["id"])
	if err != nil {
		c.BadRequest(Messages.IdRequired)
		return
	}

	signup, err := servicer.GetSignup(id)
	if err != nil {
		if err == servicer.SignupNotFound {
			c.NotFound(Messages.SignupNotFound)
		} else {
			c.InternalServerError(Messages.InternalServerError, err)
		}

		return
	}

	c.ResponseJson(signup)
}

// @Title claimSignup
// @Description Claim a signup and turn it into a user (usually from a url in an activation email).
// @Accept  json
// @Param   id              path   int     true       "The signup id"
// @Param   token           body   string  true       "The signup verification token"
// @Param   password        body   string  true       "The desired plaintext password for the new user"
// @Success 200 {object}    UserTokenResponse         ""
// @Failure 401 {object}    MessageResponse           ""
// @Failure 409 {object}    MessageResponse           ""
// @Router /signups/{id}/claim [post]
func (c *SignupContext) ClaimSignup(rw web.ResponseWriter, r *web.Request) {
	var request struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	err := c.RequestJson(&request)
	if err != nil {
		c.BadRequest(Messages.BadRequest, err)
		return
	}

	if request.Password == "" {
		c.BadRequest(Messages.PasswordRequired)
		return
	}

	if request.Token == "" {
		c.BadRequest(Messages.TokenRequired)
		return
	}

	id, err := strconv.Atoi(r.PathParams["id"])
	if err != nil {
		c.BadRequest(Messages.IdRequired)
		return
	}

	user, err := servicer.ClaimSignup(id, request.Token, request.Password)
	if err != nil {
		switch err {
		case servicer.SignupAlreadyClaimed:
			c.Conflict(Messages.UserConflict)
		case servicer.SignupNotFound:
			c.NotFound(Messages.SignupNotFound)
		case servicer.SignupInvalidToken:
			c.Unauthorized(Messages.InvalidToken)
		default:
			c.InternalServerError(Messages.InternalServerError, err)
		}
		return
	}

	tokenString, err := servicer.TokenUser(user, time.Hour*12)
	if err != nil {
		c.InternalServerError(Messages.InternalServerError, err)
		return
	}

	c.ResponseJson(&UserTokenResponse{Token: tokenString, User: user})
}
