package service

import (
	"encoding/base64"
	"encoding/json"
	opsee "github.com/opsee/basic/service"
	"github.com/opsee/vape/servicer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type service struct {
	Server *grpc.Server
}

func New() *service {
	s := &service{}

	server := grpc.NewServer()
	opsee.RegisterVapeServer(server, s)

	s.Server = server

	return s
}

func (s *service) GetBasicToken(ctx context.Context, req *opsee.TokenRequest) (*opsee.TokenResponse, error) {
	user, err := servicer.GetUserCustID(req.CustomerId)
	if err != nil {
		return nil, err
	}

	toke, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	return &opsee.TokenResponse{
		Token: base64.StdEncoding.EncodeToString(toke),
	}, nil
}

func (s *service) ListUsers(ctx context.Context, req *opsee.ListUsersRequest) (*opsee.ListUsersResponse, error) {
	users, err := servicer.ListUsers(int(req.PerPage), int(req.Page))
	if err != nil {
		return nil, err
	}

	return &opsee.ListUsersResponse{
		PerPage: int32(users.PerPage),
		Page:    int32(users.Page),
		Total:   int32(users.Total),
		Users:   users.Users,
	}, nil
}
