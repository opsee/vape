package api

import (
	"bytes"
	"github.com/gocraft/web"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/servicer"
	"strconv"
	"time"
)

type UserContext struct {
	*Context
	Id   int
	User *model.User
}

var userRouter *web.Router

// @SubApi User API [/users]
func init() {
	userRouter = publicRouter.Subrouter(UserContext{}, "/users")
	userRouter.Middleware((*UserContext).Authorized)
	userRouter.Middleware((*UserContext).FetchUser)
	userRouter.Get("/", (*UserContext).ListUsers)
	userRouter.Get("/:id", (*UserContext).GetUser)
	userRouter.Put("/:id", (*UserContext).UpdateUser)
	userRouter.Delete("/:id", (*UserContext).DeleteUser)
	userRouter.Get("/:id/data", (*UserContext).GetUserData)
	userRouter.Put("/:id/data", (*UserContext).UpdateUserData)
}

func (c *UserContext) Authorized(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if c.CurrentUser == nil {
		c.Unauthorized(Messages.UserOrAdminRequired)
		return
	}

	if r.Path != "/" {
		id, err := strconv.Atoi(r.PathParams["id"])
		if err != nil {
			c.BadRequest(Messages.IdRequired)
			return
		}
		c.Id = id
	}

	if (c.Id != 0 && c.CurrentUser.Id == c.Id) || c.CurrentUser.Admin {
		next(rw, r)
	} else {
		c.Unauthorized(Messages.UserOrAdminRequired)
	}
}

func (c *UserContext) FetchUser(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if r.Path != "/" && c.Id == 0 {
		c.BadRequest(Messages.IdRequired)
		return
	}

	user, err := servicer.GetUser(c.Id)
	if err != nil {
		if err == servicer.UserNotFound {
			c.NotFound(Messages.UserNotFound)
		} else {
			c.InternalServerError(Messages.InternalServerError, err)
		}
		return
	}

	c.User = user
	next(rw, r)
}

func (c *UserContext) ListUsers(rw web.ResponseWriter, r *web.Request) {
	perPage, err := strconv.Atoi(r.FormValue("per_page"))
	if err != nil || perPage <= 0 {
		perPage = 20
	}

	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	users, err := servicer.ListUsers(perPage, page)
	if err != nil {
		c.InternalServerError(Messages.InternalServerError, err)
		return
	}

	c.ResponseJson(users)
}

// @Title getUser
// @Description Get a single user.
// @Accept  json
// @Param   Authorization    header string  true        "The Bearer token - an admin user token or a token with matching id is required"
// @Param   id               path   int     true        "The user id"
// @Success 200 {object}     model.User                 ""
// @Failure 401 {object}     MessageResponse           	""
// @Router /users/{id} [get]
func (c *UserContext) GetUser(rw web.ResponseWriter, r *web.Request) {
	c.ResponseJson(c.User)
}

// @Title updateUser
// @Description Update a single user.
// @Accept  json
// @Param   Authorization    header string  true        "The Bearer token - an admin user token or a token with matching id is required"
// @Param   id               path   int     true        "The user id"
// @Param   email            body   string  false       "A new email address"
// @Param   name             body   string  false       "A new name"
// @Param   password         body   string  false       "A new password"
// @Success 200 {object}     UserTokenResponse          ""
// @Failure 401 {object}     MessageResponse     	 ""
// @Router /users/{id} [put]
func (c *UserContext) UpdateUser(rw web.ResponseWriter, r *web.Request) {
	var request struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	err := c.RequestJson(&request)
	if err != nil {
		c.BadRequest(Messages.BadRequest)
		return
	}

	tokenString, err := servicer.UpdateUser(c.User, request.Email, request.Name, request.Password, time.Hour*12)
	if err != nil {
		c.InternalServerError(Messages.InternalServerError, err)
		return
	}

	c.ResponseJson(&UserTokenResponse{Token: tokenString, User: c.User})
}

// @Title deleteUser
// @Description Update a single user.
// @Accept  json
// @Param   Authorization    header string  true        "The Bearer token - an admin user token or a token with matching id is required"
// @Param   id               path   int     true        "The user id"
// @Success 200 {object}     MessageResponse            ""
// @Failure 401 {object}     MessageResponse           	""
// @Router /users/{id} [delete]
func (c *UserContext) DeleteUser(rw web.ResponseWriter, r *web.Request) {
	err := servicer.DeleteUser(c.Id)
	if err != nil {
		c.InternalServerError(Messages.InternalServerError, err)
		return
	}

	c.ResponseJson(&MessageResponse{Message: Messages.UserDeleted})
}

type UserDataResponse map[string]interface{}

// @Title getUserData
// @Description Update a single user.
// @Accept  json
// @Param   Authorization    header string  true        "The Bearer token - an admin user token or a token with matching id is required"
// @Param   id               path   int     true        "The user id"
// @Success 200 {object}     UserDataResponse           ""
// @Failure 401 {object}     MessageResponse    	 ""
// @Router /users/{id}/data [get]
func (c *UserContext) GetUserData(rw web.ResponseWriter, r *web.Request) {
	data, err := servicer.GetUserData(c.Id)
	if err != nil {
		c.InternalServerError(Messages.InternalServerError, err)
	}

	rw.Write(data)
}

// @Title updateUserData
// @Description Update a single user.
// @Accept  json
// @Param   Authorization    header string  true        "The Bearer token - an admin user token or a token with matching id is required"
// @Param   id               path   int     true        "The user id"
// @Success 200 {object}     UserDataResponse           ""
// @Failure 401 {object}     MessageResponse    	 ""
// @Router /users/{id}/data [put]
func (c *UserContext) UpdateUserData(rw web.ResponseWriter, r *web.Request) {
	buf := bytes.NewBuffer([]byte{})
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		c.BadRequest(Messages.BadRequest)
		return
	}

	data, err := servicer.UpdateUserData(c.Id, buf.Bytes())
	if err != nil {
		c.InternalServerError(Messages.InternalServerError, err)
	}

	rw.Write(data)
}
