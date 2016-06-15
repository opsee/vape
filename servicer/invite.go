package servicer

import (
	"fmt"

	"github.com/opsee/basic/schema"
	opsee_types "github.com/opsee/protobuf/opseeproto/types"
	log "github.com/sirupsen/logrus"
)

func CreateActiveInvite(teamName, senderEmail, customerId, email string, perms *opsee_types.Permission) (*schema.Invite, error) {
	referrer := ""
	signup, err := createSignup(customerId, email, "", referrer, true, perms)
	if err != nil {
		return nil, err
	}

	// fuck it, do it here too
	if signup.Perms != nil {
		signup.Perms.Name = "user"
	}

	log.Debugf("invite signup %v permissions: %v", signup, signup.Perms)
	invite := &schema.Invite{
		Id:         int32(signup.Id),
		Email:      signup.Email,
		CustomerId: signup.CustomerId,
		Perms:      signup.Perms,
	}

	// send an email, create a lead and notify slack here!
	go func() {
		mergeVars := map[string]interface{}{
			"signup_id":    fmt.Sprint(signup.Id),
			"signup_token": VerificationToken(fmt.Sprintf("%d", signup.Id)),
			"team_name":    signup.Name,
			"sender_email": senderEmail,
		}
		mailTemplatedMessage(signup.Email, signup.Name, "team-invitation", mergeVars)
	}()

	return invite, nil
}
