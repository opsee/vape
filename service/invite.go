package service

import (
	"fmt"

	opsee "github.com/opsee/basic/service"
	"github.com/opsee/vape/servicer"
	"golang.org/x/net/context"
)

func (s *service) InviteUser(ctx context.Context, req *opsee.InviteUserRequest) (*opsee.InviteUserResponse, error) {
	var (
		err error
	)

	// TODO(dan) this could be used as a side-channel to find valid email addresses maybe
	user, err := servicer.GetUserEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, fmt.Errorf("user exists with this email address")
	}

	teamName := ""
	team, _ := servicer.GetTeam(req.Requestor.CustomerId)
	if team != nil {
		teamName = team.Name
	}

	senderEmail := req.Requestor.Email

	invite, err := servicer.CreateActiveInvite(teamName, senderEmail, req.Requestor.CustomerId, req.Email, req.Perms)
	if err != nil {
		return nil, err
	}

	return &opsee.InviteUserResponse{
		Invite: invite,
	}, nil
}
