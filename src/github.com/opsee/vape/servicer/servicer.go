package servicer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hoisie/mustache"
	"github.com/keighl/mandrill"
	slacktmpl "github.com/opsee/notification-templates/dist/go/slack"
	log "github.com/sirupsen/logrus"
	"github.com/snorecone/closeio-go"
	"net/http"
	"net/url"
)

type MandrillMailer interface {
	MessagesSendTemplate(*mandrill.Message, string, interface{}) ([]*mandrill.Response, error)
}

var (
	opseeHost       string
	mailClient      MandrillMailer
	intercomKey     []byte
	closeioClient   *closeio.Closeio
	slackEndpoint   string
	slackTemplates  map[string]*mustache.Template
	slackDomain     string
	slackAdminToken string
)

func init() {
	slackTemplates = make(map[string]*mustache.Template)

	tmpl, err := mustache.ParseString(slacktmpl.NewSignup)
	if err != nil {
		panic(err)
	}

	slackTemplates["new-signup"] = tmpl
}

func Init(host string, mailer MandrillMailer, intercom, closeioKey, slackUrl, slackDomain, slackAdminToken string) {
	opseeHost = host
	mailClient = mailer
	intercomKey = []byte(intercom)
	slackEndpoint = slackUrl
	slackDomain = slackDomain
	slackAdminToken = slackAdminToken

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

func inviteSlack(name, email string) {
	if slackDomain == "" || slackAdminToken == "" {
		log.Warn("not inviting user to the opsee support slack")
		return
	}

	log.Info("inviting user to the opsee support slack")

	// find out if they are already in the slack
	resp, err := http.Get(fmt.Sprintf("https://%s/api/users.list&token=%s", slackDomain, slackAdminToken))
	if err != nil {
		log.WithError(err).Error("failed to get a list of users from slack")
		return
	}

	defer resp.Body.Close()

	var usersResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&usersResponse)
	if err != nil {
		log.WithError(err).Error("failed to decode list of users from slack")
		return
	}

	members, ok := usersResponse["members"].([]map[string]interface{})
	if !ok {
		log.WithError(err).Error("failed to decode list of users from slack")
		return
	}

	for _, member := range members {
		if profile, ok := member["profile"].(map[string]interface{}); ok {
			if em, ok := profile["email"].(string); ok {
				if em == email {
					log.Warn("user with email %s was already invited to slack", em)
					return
				}
			}
		}
	}

	v := url.Values{}
	v.Set("email", email)
	v.Set("token", slackAdminToken)
	v.Set("set_active", "true")
	v.Set("extra_message", fmt.Sprintf(`Thanks for signing up for Opsee, %s. If you use Slack, you can chat with our engineering and support teams any time in the opsee-support Slack team.`, name))

	resp, err = http.PostForm(fmt.Sprintf("https://%s/api/users.admin.invite", slackDomain), v)
	if err != nil {
		log.WithError(err).Errorf("failed to send slack invitation to: %s", email)
		return
	}

	defer resp.Body.Close()
	log.Infof("slack invitation sent to %s", email)
}
