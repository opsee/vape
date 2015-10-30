package api

import (
	"github.com/gocraft/web"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/servicer"
)

type NotificationContext struct {
	*Context
}

var (
	notificationRouter     *web.Router
	testNotificationRouter *web.Router
)

func init() {
	notificationRouter = privateRouter.Subrouter(NotificationContext{}, "/notifications")
	notificationRouter.Post("/send/email", (*NotificationContext).SendEmail)

	testNotificationRouter = publicRouter.Subrouter(NotificationContext{}, "/notifications")
	testNotificationRouter.Post("/test/email", (*NotificationContext).SendTestEmail)
}

type SendEmalResponse struct {
	User *model.User `json:"user"`
}

func (c *NotificationContext) SendEmail(rw web.ResponseWriter, r *web.Request) {
	var request struct {
		UserId   int                    `json:"user_id"`
		Template string                 `json:"template"`
		Vars     map[string]interface{} `json:"vars"`
	}

	err := c.RequestJson(&request)
	if err != nil {
		c.BadRequest(Messages.BadRequest, err)
		return
	}

	if request.UserId == 0 {
		c.BadRequest(Messages.UserIdRequired)
		return
	}

	if request.Template == "" {
		c.BadRequest(Messages.TemplateRequired)
		return
	}

	user, err := servicer.SendTemplatedEmail(request.UserId, request.Template, request.Vars)
	if err != nil {
		c.InternalServerError(Messages.InternalServerError, err)
		return
	}

	c.ResponseJson(&SendEmalResponse{User: user})
}

func (c *NotificationContext) SendTestEmail(rw web.ResponseWriter, r *web.Request) {
	if c.CurrentUser == nil || c.CurrentUser.Admin != true {
		c.Unauthorized(Messages.AdminRequired)
		return
	}

	c.SendEmail(rw, r)
}
