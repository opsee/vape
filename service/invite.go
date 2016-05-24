package service

import (
	"github.com/opsee/basic/schema"
	opsee "github.com/opsee/basic/service"
	"github.com/opsee/vape/servicer"
	"golang.org/x/net/context"
)

func (s *service) InviteUser(ctx context.Context, req *opsee.InviteUserRequest) (*opsee.InviteUserResponse, error) {
	// Only team admins or an OpseeAdmin can perform this action
	if err := schema.CheckPermission(req.Requestor, "admin"); err != nil {
		return nil, err
	}

	var (
		err error
	)

	name := ""
	invite, err := servicer.CreateActiveInvite(req.Requestor.CustomerId, req.Email, name, req.Perms)
	if err != nil {
		return nil, err
	}

	return &opsee.InviteUserResponse{
		Invite: invite,
	}, nil
}
