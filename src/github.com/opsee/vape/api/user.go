package api

import (
        "github.com/gocraft/web"
        "github.com/opsee/vape/model"
        "github.com/opsee/vape/servicer"
        "net/http"
        "strconv"
        "fmt"
)

type UserContext struct {
        *Context
        Id int
        User *model.User
}

var userRouter *web.Router

func init() {
        userRouter = router.Subrouter(UserContext{}, "/users")
        userRouter.Middleware((*UserContext).Authorized)
        userRouter.Middleware((*UserContext).SetUserContext)
        userRouter.Get("/:id", (*UserContext).GetUser)
        userRouter.Put("/:id", (*UserContext).UpdateUser)
        userRouter.Delete("/:id", (*UserContext).DeleteUser)
}

func (c *UserContext) Authorized(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
        id, ok := r.PathParams["id"]
        if ok {
                userId, err := strconv.Atoi(id)
                if err != nil {
                        c.Job.EventErr("user-authorized", err)
                        rw.WriteHeader(http.StatusBadRequest)
                        return
                }

                c.Id = userId
        }

        if (c.Id != 0 && c.CurrentUser.Id == c.Id) || c.CurrentUser.Admin {
                next(rw, r)
        } else {
                c.Job.EventKv("user-authorized", map[string]string{"user_id": fmt.Sprintf("%s", c.Id)})
                rw.WriteHeader(http.StatusUnauthorized)
        }
}

func (c *UserContext) SetUserContext(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
        if c.Id == 0 {
                c.Job.EventKv("user-get", map[string]string{"user_id": fmt.Sprintf("%s", c.Id)})
                rw.WriteHeader(http.StatusBadRequest)
                return
        }

        user, err := servicer.GetUser(c.Id)
        if err != nil {
                c.Job.EventErr("user-get", err)
                rw.WriteHeader(http.StatusInternalServerError)
                return
        }

        if user == nil {
                rw.WriteHeader(http.StatusNotFound)
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

