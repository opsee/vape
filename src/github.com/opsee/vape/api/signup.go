package api

import (
        "github.com/gocraft/web"
        "github.com/opsee/vape/model"
        "github.com/opsee/vape/servicer"
        "net/http"
        "strconv"
        "fmt"
)

type SignupContext struct {
        *Context
        Id int
        Signup *model.Signup
}

var signupRouter *web.Router

func init() {
        signupRouter = router.Subrouter(SignupContext{}, "/signups")
        userRouter.Post("/", (*SignupContext).CreateSignup)
        userRouter.Get("/", (*SignupContext).ListSignups)
        userRouter.Get("/:id", (*SignupContext).GetSignup)
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

        signups, err := servicer.ListSignups(perPage, page)
        if err != nil {
                rw.WriteHeader(http.StatusInternalServerError)
                return
        }

        writeJson(rw, map[string]interface{}{"signups": signups})
}

func (c *SignupContext) CreateSignup(rw web.ResponseWriter, r *web.Request) {
        // anyone is authorized for this
        signup, err := servicer.CreateSignup()
        if err != nil {
                rw.WriteHeader(http.StatusInternalServerError)
                return
        }

        writeJson(rw, map[string]interface{}{"signup": signup})
}

func (c *SignupContext) GetSignup(rw web.ResponseWriter, r *web.Request) {
        if c.CurrentUser == nil || c.CurrentUser.Admin != true {
                rw.WriteHeader(http.StatusUnauthorized)
                return
        }

        err := c.FetchSignup()
        if err != nil {
                return
        }

        writeJson(rw, map[string]interface{}{"signup": c.Signup})
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
