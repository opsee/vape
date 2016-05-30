package service

import (
	opsee "github.com/opsee/basic/service"
	"github.com/opsee/vape/servicer"
	"golang.org/x/net/context"
)

func (s *service) InviteUser(ctx context.Context, req *opsee.InviteUserRequest) (*opsee.InviteUserResponse, error) {
	var (
		err error
	)

	teamName := ""
	team, _ := servicer.GetTeam(req.Requestor.CustomerId)
	if team != nil {
		teamName = team.Name
	}

	senderEmail := req.Requestor.Email

	invite, err := servicer.CreateActiveInvite(teamName, senderEmail, req.Requestor.CustomerId, req.Email, req.Name, req.Perms)
	if err != nil {
		return nil, err
	}

	return &opsee.InviteUserResponse{
		Invite: invite,
	}, nil
}
