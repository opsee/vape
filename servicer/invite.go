package servicer

import (
	"github.com/opsee/basic/schema"
	opsee_types "github.com/opsee/protobuf/opseeproto/types"
)

// TODO(dan)  This currently piggy-backs on signup.  It should be moved to its own table.
// TODO(dan) Email templated message for invite
func CreateActiveInvite(customerId, email, name string, perms *opsee_types.Permission) (*schema.Invite, error) {
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

	return invite, nil
}
