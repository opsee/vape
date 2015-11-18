package servicer

import (
	"github.com/keighl/mandrill"
)

type MandrillMailer interface {
	MessagesSendTemplate(*mandrill.Message, string, interface{}) ([]*mandrill.Response, error)
}

var (
	opseeHost   string
	mailClient  MandrillMailer
	intercomKey []byte
)

func Init(host string, mailer MandrillMailer, intercom string) {
	opseeHost = host
	mailClient = mailer
	intercomKey = []byte(intercom)
}

func mailTemplatedMessage(toEmail, toName, templateName string, mergeVars map[string]interface{}) ([]*mandrill.Response, error) {
	if mailClient == nil {
		return nil, nil
	}

	mergeVars["opsee_host"] = opseeHost

	message := &mandrill.Message{}
	message.AddRecipient(toEmail, toName, "to")
	message.Merge = true
	message.MergeLanguage = "handlebars"
	message.MergeVars = []*mandrill.RcptMergeVars{mandrill.MapToRecipientVars(toEmail, mergeVars)}
	return mailClient.MessagesSendTemplate(message, templateName, map[string]string{})
}
