package api

import (
	"github.com/gocraft/web"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/servicer"
	"net/http"
	"strconv"
)

type UserContext struct {
	*Context
	Id   int
	User *model.User
}

var userRouter *web.Router

func init() {
	userRouter = publicRouter.Subrouter(UserContext{}, "/users")
	userRouter.Middleware((*UserContext).Authorized)
	userRouter.Middleware((*UserContext).FetchUser)
	userRouter.Get("/:id", (*UserContext).GetUser)
	userRouter.Put("/:id", (*UserContext).UpdateUser)
	userRouter.Delete("/:id", (*UserContext).DeleteUser)
}

func (c *UserContext) Authorized(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if c.CurrentUser == nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.PathParams["id"])
	if err != nil {
		c.Job.EventErr("error.atoi", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	c.Id = id

	if (c.Id != 0 && c.CurrentUser.Id == c.Id) || c.CurrentUser.Admin {
		next(rw, r)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
	}
}

func (c *UserContext) FetchUser(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if c.Id == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := servicer.GetUser(c.Id)
	if err != nil {
		c.Job.EventErr("error.getuser", err)
		if err == servicer.RecordNotFound {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	c.User = user
	next(rw, r)
}

func (c *UserContext) GetUser(rw web.ResponseWriter, r *web.Request) {
	writeJson(rw, c.User)
}

func (c *UserContext) UpdateUser(rw web.ResponseWriter, r *web.Request) {
	userJson, err := readJson(r)
	if err != nil {
		c.Job.EventErr("error.json", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = servicer.UpdateUser(c.User, userJson)
	if err != nil {
		c.Job.Event("error.update")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJson(rw, c.User)
}

func (c *UserContext) DeleteUser(rw web.ResponseWriter, r *web.Request) {
	err := servicer.DeleteUser(c.Id)
	if err != nil {
		c.Job.Event("error.delete")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
