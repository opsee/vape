package api

import (
	"github.com/gocraft/web"
	"github.com/opsee/vape/servicer"
)

type BastionContext struct {
	*Context
}

var bastionRouter *web.Router

func init() {
	bastionRouter = privateRouter.Subrouter(BastionContext{}, "/bastions")
	bastionRouter.Post("/", (*BastionContext).Create)
	bastionRouter.Post("/authenticate", (*BastionContext).Authenticate)
}

type BastionResponse struct {
	Id         string `json:"id"`
	Password   string `json:"password"`
	CustomerId string `json:"customer_id"`
}

func (c *BastionContext) Create(rw web.ResponseWriter, r *web.Request) {
	var request struct {
		CustomerId string `json:"customer_id"`
	}

	err := c.RequestJson(&request)
	if err != nil {
		c.BadRequest("malformed request", err)
		return
	}

	if request.CustomerId == "" {
		c.BadRequest("missing customer_id")
		return
	}

	bastion, plaintext, err := servicer.CreateBastion(request.CustomerId)
	if err != nil {
		if err == servicer.CustomerNotFound {
			c.Unauthorized("no active customer with that id")
		} else {
			c.InternalServerError("internal server error", err)
		}

		return
	}

	c.ResponseJson(&BastionResponse{Id: bastion.Id, Password: plaintext, CustomerId: request.CustomerId})
}

func (c *BastionContext) Authenticate(rw web.ResponseWriter, r *web.Request) {
	var request struct {
		Id       string `json:"id"`
		Password string `json:"password"`
	}

	err := c.RequestJson(&request)
	if err != nil {
		c.BadRequest("malformed request", err)
		return
	}

	if request.Id == "" {
		c.BadRequest("missing id")
		return
	}

	if request.Password == "" {
		c.BadRequest("missing password")
		return
	}

	if err = servicer.AuthenticateBastion(request.Id, request.Password); err != nil {
		c.Unauthorized("couldn't authenticate bastion", err, map[string]string{"id": request.Id})
		return
	}
}
