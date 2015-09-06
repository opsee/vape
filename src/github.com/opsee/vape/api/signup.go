package api

import (
	"github.com/gocraft/web"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/servicer"
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
// @Failure 401 {object}    MessageResponse           	 ""
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
// @Failure 409 {object}     MessageResponse           "Email was already used to sign up"
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
		c.BadRequest("name is required")
		return
	}

	if request.Email == "" {
		c.BadRequest("email is required")
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

type SignupResponse struct {
	Token string `json:"token"`
}

// @Title activateSignup
// @Description Sends the activation email for a signup. Can be called multiple times to send multiple emails.
// @Accept  json
// @Param   id               path      int   true   "The signup's id"
// @Success 200 {object}     SignupResponse         "An object with the claim token used to verify the signup (sent in email)"
// @Failure 401 {object}     MessageResponse        ""
// @Router /signups/{id}/activate [put]
func (c *SignupContext) ActivateSignup(rw web.ResponseWriter, r *web.Request) {
	if c.CurrentUser == nil || c.CurrentUser.Admin != true {
		c.Unauthorized("must be an administrator to access this resource")
		return
	}

	id, err := strconv.Atoi(r.PathParams["id"])
	if err != nil {
		c.BadRequest("need a valid id in request path")
		return
	}

	// sending referer should be temporary for developing email templates
	referer := r.Header.Get("Referer")
	signup, err := servicer.ActivateSignup(id, referer)
	if err != nil {
		if err == servicer.SignupNotFound {
			c.NotFound("signup not found")
		} else {
			c.InternalServerError("internal server error", err)
		}

		return
	}

	c.ResponseJson(&SignupResponse{Token: signup.Token()})
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
		c.Unauthorized("must be an administrator to access this resource")
		return
	}

	id, err := strconv.Atoi(r.PathParams["id"])
	if err != nil {
		c.BadRequest("need a valid id in request path")
		return
	}

	signup, err := servicer.GetSignup(id)
	if err != nil {
		if err == servicer.SignupNotFound {
			c.NotFound("signup not found")
		} else {
			c.InternalServerError("internal server error", err)
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
		c.BadRequest("malformed request", err)
		return
	}

	if request.Password == "" {
		c.BadRequest("password is required")
		return
	}

	if request.Token == "" {
		c.BadRequest("token is required")
		return
	}

	id, err := strconv.Atoi(r.PathParams["id"])
	if err != nil {
		c.BadRequest("need a valid id in request path")
		return
	}

	user, err := servicer.ClaimSignup(id, request.Token, request.Password)
	if err != nil {
		switch err {
		case servicer.SignupAlreadyClaimed:
			c.Conflict("user has already been claimed")
		case servicer.SignupNotFound:
			c.NotFound("signup not found")
		case servicer.SignupInvalidToken:
			c.Unauthorized("invalid token for signup")
		default:
			c.InternalServerError("internal server error", err)
		}
		return
	}

	tokenString, err := servicer.TokenUser(user)
	if err != nil {
		c.InternalServerError("internal server error", err)
		return
	}

	c.ResponseJson(&UserTokenResponse{Token: tokenString, User: user})
}
