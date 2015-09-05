package servicer

import (
	"github.com/keighl/mandrill"
)

const mandrillAPIKey = "V2M1onmVdOXJ42Vr8Gr_ew"

var (
	mailClient *mandrill.Client
)

func init() {
	mailClient = mandrill.ClientWithKey(mandrillAPIKey)
}

func mailTemplatedMessage(toEmail, toName, templateName string, mergeVars interface{}) ([]*mandrill.Response, error) {
	message := &mandrill.Message{}
	message.AddRecipient(toEmail, toName, "to")
	message.Merge = true
	message.MergeLanguage = "handlebars"
	message.MergeVars = []*mandrill.RcptMergeVars{mandrill.MapToRecipientVars(toEmail, mergeVars)}
	return mailClient.MessagesSendTemplate(message, templateName, map[string]string{})
}
