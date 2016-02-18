package api

import (
	"github.com/gocraft/web"
	"github.com/opsee/vape/servicer"
)

type AdminContext struct {
	*Context
}

var adminRouter *web.Router

// @SubApi Admin API [/admin]
func init() {
	adminRouter = privateRouter.Subrouter(AdminContext{}, "/admin")
	adminRouter.Get("/users/:cust_id", (*AdminContext).GetUserCustID)
}

func (c *AdminContext) GetUserCustID(rw web.ResponseWriter, r *web.Request) {
	// TODO probably should sanitize custID before passing to db
	custID, exists := r.PathParams["cust_id"]
	if !exists {
		c.BadRequest(Messages.CustomerIdRequired)
		return
	}

	user, err := servicer.GetUserCustID(custID)
	if err != nil {
		c.InternalServerError(Messages.InternalServerError, err)
		return
	}

	c.ResponseJson(user)
}
