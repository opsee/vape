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
		c.BadRequest(Messages.BadRequest, err)
		return
	}

	if request.CustomerId == "" {
		c.BadRequest(Messages.CustomerIdRequired)
		return
	}

	bastion, plaintext, err := servicer.CreateBastion(request.CustomerId)
	if err != nil {
		if err == servicer.CustomerNotFound {
			c.Unauthorized(Messages.CustomerNotAuthorized)
		} else {
			c.InternalServerError(Messages.InternalServerError, err)
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
		c.BadRequest(Messages.BadRequest, err)
		return
	}

	if request.Id == "" {
		c.BadRequest(Messages.IdRequired)
		return
	}

	if request.Password == "" {
		c.BadRequest(Messages.PasswordRequired)
		return
	}

	if err = servicer.AuthenticateBastion(request.Id, request.Password); err != nil {
		c.Unauthorized(Messages.BastionCredentialsMismatch, err, map[string]string{"id": request.Id})
		return
	}
}
