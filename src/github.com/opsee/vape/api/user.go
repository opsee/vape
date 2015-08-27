package api

// import (
//         "github.com/gocraft/web"
//         "github.com/opsee/vape/model"
//         "github.com/opsee/vape/store"
//         "github.com/opsee/vape/token"
//         "net/http"
//         "time"
// )
//
// type UserContext struct {
//         *Context
// }
//
// var userRouter *web.Router
//
// func init() {
//         userRouter = router.Subrouter(UserContext{}, "/users")
//         userRouter.Get("/users/:id", (*UserContext).GetUser)
//         userRouter.Put("/users/:id", (*UserContext).UpdateUser)
//         userRouter.Delete("/users/:id", (*UserContext).DeleteUser)
// }
//
// func (c *UserContext) GetUser(rw web.ResponseWriter, r *web.Request) {
//
// }
//
// func (c *UserContext) UpdateUser(rw web.ResponseWriter, r *web.Request) {
//
// }
//
// func (c *UserContext) DeleteUser(rw web.ResponseWriter, r *web.Request) {
//
// }
//
