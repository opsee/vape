package servicer

import (
	"github.com/snorecone/closeio-go"
	"github.com/keighl/mandrill"
	"log"
)

type MandrillMailer interface {
	MessagesSendTemplate(*mandrill.Message, string, interface{}) ([]*mandrill.Response, error)
}

var (
	opseeHost     string
	mailClient    MandrillMailer
	intercomKey   []byte
	closeioClient *closeio.Closeio
)

func Init(host string, mailer MandrillMailer, intercom, closeioKey string) {
	opseeHost = host
	mailClient = mailer
	intercomKey = []byte(intercom)

	if closeioKey != "" {
		closeioClient = closeio.New(closeioKey)
	}
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

func createLead(lead *closeio.Lead) {
	if closeioClient != nil {
		resp, err := closeioClient.CreateLead(lead)
		if err != nil {
			log.Print(err.Error())
		} else {
			log.Printf("created closeio lead: %s", resp.Url)
		}
	}
}
