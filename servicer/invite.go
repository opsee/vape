package servicer

import (
	"fmt"

	"github.com/opsee/basic/schema"
	opsee_types "github.com/opsee/protobuf/opseeproto/types"
)

// TODO(dan)  This currently piggy-backs on signup.  It should be moved to its own table.
// TODO(dan) Email templated message for invite
func CreateActiveInvite(teamName, senderEmail, customerId, email, name string, perms *opsee_types.Permission) (*schema.Invite, error) {
	referrer := ""
	signup, err := createSignup(customerId, email, name, referrer, true, perms)
	if err != nil {
		return nil, err
	}
	invite := &schema.Invite{
		Id:         int32(signup.Id),
		CustomerId: signup.CustomerId,
		Name:       signup.Name,
		Perms:      signup.Perms,
	}

	// send an email, create a lead and notify slack here!
	go func() {
		mergeVars := map[string]interface{}{
			"signup_id":    fmt.Sprint(signup.Id),
			"signup_token": signup.Token(),
			"team_name":    signup.Name,
			"name":         name,
			"sender_email": senderEmail,
		}
		mailTemplatedMessage(signup.Email, signup.Name, "team-invitation", mergeVars)
	}()

	return invite, nil
}
