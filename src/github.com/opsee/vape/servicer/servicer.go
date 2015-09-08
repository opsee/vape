package servicer

import (
	"github.com/keighl/mandrill"
)

type MandrillMailer interface {
	MessagesSendTemplate(*mandrill.Message, string, interface{}) ([]*mandrill.Response, error)
}

var (
	mailClient MandrillMailer
)

func Init(mailer MandrillMailer) {
	mailClient = mailer
}

func mailTemplatedMessage(toEmail, toName, templateName string, mergeVars interface{}) ([]*mandrill.Response, error) {
	if mailClient == nil {
		return nil, nil
	}

	message := &mandrill.Message{}
	message.AddRecipient(toEmail, toName, "to")
	message.Merge = true
	message.MergeLanguage = "handlebars"
	message.MergeVars = []*mandrill.RcptMergeVars{mandrill.MapToRecipientVars(toEmail, mergeVars)}
	return mailClient.MessagesSendTemplate(message, templateName, map[string]string{})
}
