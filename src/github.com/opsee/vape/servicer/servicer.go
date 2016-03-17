package servicer

import (
	"bytes"
	"fmt"
	"github.com/hoisie/mustache"
	"github.com/keighl/mandrill"
	slacktmpl "github.com/opsee/notification-templates/dist/go/slack"
	log "github.com/sirupsen/logrus"
	"github.com/snorecone/closeio-go"
	"net/http"
)

type MandrillMailer interface {
	MessagesSendTemplate(*mandrill.Message, string, interface{}) ([]*mandrill.Response, error)
}

var (
	opseeHost      string
	mailClient     MandrillMailer
	intercomKey    []byte
	closeioClient  *closeio.Closeio
	slackEndpoint  string
	slackTemplates map[string]*mustache.Template
)

func init() {
	slackTemplates = make(map[string]*mustache.Template)

	tmpl, err := mustache.ParseString(slacktmpl.NewSignup)
	if err != nil {
		panic(err)
	}

	slackTemplates["new-signup"] = tmpl
}

func Init(host string, mailer MandrillMailer, intercom, closeioKey, slackUrl string) {
	opseeHost = host
	mailClient = mailer
	intercomKey = []byte(intercom)
	slackEndpoint = slackUrl

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

func notifySlack(name string, vars map[string]interface{}) {
	log.Info("requested slack notification")

	template, ok := slackTemplates[name]
	if !ok {
		log.Errorf("not sending slack notification since template %s was not found", name)
		return
	}

	body := template.Render(vars)

	if slackEndpoint == "" {
		log.Warn("not sending slack notification since SLACK_ENDPOINT is not set")
		fmt.Println(body)
		return
	}

	resp, err := http.Post(slackEndpoint, "application/json", bytes.NewBufferString(body))
	if err != nil {
		log.WithError(err).Errorf("failed to send slack notification: %s", name)
		return
	}

	defer resp.Body.Close()
	log.WithField("status", resp.StatusCode).Info("sent slack request")
}
