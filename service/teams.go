package service

import (
	"github.com/opsee/basic/schema"
	opsee "github.com/opsee/basic/service"
	"github.com/opsee/vape/servicer"
	"golang.org/x/net/context"
)

// Fetches team, including users
func (s *service) GetTeam(ctx context.Context, req *opsee.GetTeamRequest) (*opsee.GetTeamResponse, error) {
	var (
		team *schema.Team
		err  error
	)

	team, err = servicer.GetTeam(req.Team.Id)
	if err != nil {
		return nil, err
	}

	return &opsee.GetTeamResponse{
		Team: team,
	}, nil
}

// Updates team name or subscription
func (s *service) UpdateTeam(ctx context.Context, req *opsee.UpdateTeamRequest) (*opsee.UpdateTeamResponse, error) {
	var (
		team *schema.Team
		err  error
	)

	team, err = servicer.UpdateTeam(req.Team, req.Team.Name, req.Team.Subscription)
	if err != nil {
		return nil, err
	}

	return &opsee.UpdateTeamResponse{
		Team: team,
	}, nil
}

// Sets team to inactive
func (s *service) DeleteTeam(ctx context.Context, req *opsee.DeleteTeamRequest) (*opsee.DeleteTeamResponse, error) {
	var (
		err error
	)

	err = servicer.DeleteTeam(req.Team.Id)
	if err != nil {
		return nil, err
	}

	return &opsee.DeleteTeamResponse{
		Team: req.Team,
	}, nil
}
