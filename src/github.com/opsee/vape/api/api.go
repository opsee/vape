package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocraft/health"
	"github.com/gocraft/web"
	"github.com/nu7hatch/gouuid"
	"github.com/opsee/vape/model"
	"github.com/opsee/vaper"
	"io"
	"net/http"
	"runtime"
	"strings"
)

type Context struct {
	Job                 *health.Job
	Panic               bool
	CurrentUser         *model.User
	RequestJson         func(interface{}) error
	ResponseJson        func(interface{})
	BadRequest          func(string, ...interface{})
	Unauthorized        func(string, ...interface{})
	Conflict            func(string, ...interface{})
	NotFound            func(string, ...interface{})
	InternalServerError func(string, ...interface{})
}

var (
	stream        = health.NewStream()
	publicRouter  = web.New(Context{})
	privateRouter = web.New(Context{})
	origins       = []string{
		"http://localhost:8080",
		"https://staging.opsy.co",
		"https://opsee.co",
	}
)

type MessageResponse struct {
	Message string `json:"message"`
}

// @APIVersion 0.0.1
// @APITitle Vape API
// @APIDescription API for user/customer management and authentication

func init() {
	// we're creating a separate router instances to listen on separate ports
	// as a result, we have to be repeat ourselves
	for _, router := range []*web.Router{publicRouter, privateRouter} {
		router.Middleware((*Context).HelperFuncs)
		router.Middleware((*Context).Log)
		router.Middleware((*Context).CatchPanics)
		router.Middleware((*Context).Cors)
		router.Middleware((*Context).Options)
		router.Middleware((*Context).SetContentType)
		router.Middleware((*Context).UserSession)
		router.NotFound(notFound)
		router.Get("/health", (*Context).Health)
		router.Get("/swagger.json", (*Context).Docs)
	}
}

func InjectLogger(sink io.Writer) {
	if sink != nil {
		stream.AddSink(&health.WriterSink{sink})
	}
}

func ListenAndServe(publicAddr string, privateAddr string) {
	stream.EventKv("api.listen-and-serve", map[string]string{"public_host": publicAddr, "private_host": privateAddr})
	go http.ListenAndServe(publicAddr, publicRouter)
	http.ListenAndServe(privateAddr, privateRouter)
}

//
// endpoints
//
func (c *Context) Health(rw web.ResponseWriter, r *web.Request) {}

func (c *Context) Docs(rw web.ResponseWriter, r *web.Request) {
	rw.Write([]byte(swaggerJson))
}

//
// middleware
//
func (c *Context) HelperFuncs(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.RequestJson = func(s interface{}) error {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(s)
		if err != nil {
			return err
		}
		return nil
	}

	c.ResponseJson = func(data interface{}) {
		encoder := json.NewEncoder(rw)
		if err := encoder.Encode(data); err != nil {
			panic(err)
		}
	}

	c.BadRequest = c.responseFunc(rw, http.StatusBadRequest)
	c.Unauthorized = c.responseFunc(rw, http.StatusUnauthorized)
	c.Conflict = c.responseFunc(rw, http.StatusConflict)
	c.InternalServerError = c.responseFunc(rw, http.StatusInternalServerError)
	c.NotFound = c.responseFunc(rw, http.StatusNotFound)

	next(rw, r)
}

func (c *Context) UserSession(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	auth := r.Header.Get("Authorization")
	authslice := strings.Split(auth, " ")

	if len(authslice) >= 2 {
		switch authslice[0] {
		case "Bearer":
			tokenString := authslice[1]
			decodedToken, err := vaper.Unmarshal(tokenString)
			if err != nil {
				c.Job.EventErr("user_session.token_unmarshal", err)
				break
			}

			user := &model.User{}
			err = decodedToken.Reify(user)
			if err != nil {
				c.Job.EventErr("user_session.token_reify", err)
				break
			}

			c.CurrentUser = user
		}
	}

	next(rw, r)
}

func (c *Context) Log(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	c.Job = stream.NewJob(r.RoutePath())

	id, err := uuid.NewV4()
	if err == nil {
		c.Job.KeyValue("request-id", id.String())
	}

	path := r.URL.Path
	c.Job.EventKv("api.request", health.Kvs{"path": path})

	next(rw, r)

	code := rw.StatusCode()
	kvs := health.Kvs{
		"code": fmt.Sprint(code),
		"path": path,
	}

	// Map HTTP status code to category.
	var status health.CompletionStatus
	if c.Panic {
		status = health.Panic
	} else if code < 400 {
		status = health.Success
	} else if code == 422 {
		status = health.ValidationError
	} else if code < 500 {
		status = health.Junk // 404, 401
	} else {
		status = health.Error
	}
	c.Job.CompleteKv(status, kvs)
}

func (c *Context) CatchPanics(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	defer func() {
		if err := recover(); err != nil {
			c.Panic = true

			const size = 4096
			stack := make([]byte, size)
			stack = stack[:runtime.Stack(stack, false)]

			// err turns out to be interface{}, of actual type "runtime.errorCString"
			// The health package kinda wants an error. Luckily, the err sprints nicely via fmt.
			errorishError := errors.New(fmt.Sprint(err))

			c.Job.EventErrKv("panic", errorishError, health.Kvs{"stack": string(stack)})
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}()

	next(rw, r)
}

func (c *Context) SetContentType(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	header := rw.Header()
	header.Set("Content-Type", "application/json; charset=utf-8")
	next(rw, r)
}

func (c *Context) Cors(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	origin := r.Header.Get("Origin")

	for _, o := range origins {
		if o == origin {
			header := rw.Header()
			header.Set("Access-Control-Allow-Origin", o)
			header.Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
			header.Set("Access-Control-Allow-Headers", "Accept-Encoding,Authorization,Content-Type")
			header.Set("Access-Control-Max-Age", "1728000")
		}
	}
	next(rw, r)
}

func (c *Context) Options(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if r.Method == "OPTIONS" {
		return
	}
	next(rw, r)
}

func notFound(rw web.ResponseWriter, r *web.Request) {
	rw.WriteHeader(http.StatusNotFound)
}

func (c *Context) responseFunc(rw web.ResponseWriter, status int) func(string, ...interface{}) {
	return func(msg string, args ...interface{}) {
		rw.WriteHeader(status)
		c.ResponseJson(MessageResponse{Message: msg})

		if len(args) == 1 {
			c.Job.EventErr(msg, args[0].(error))
		}
		if len(args) == 2 {
			c.Job.EventErrKv(msg, args[0].(error), args[1].(map[string]string))
		}
	}
}
