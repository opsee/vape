package service

import (
	"encoding/base64"
	"encoding/json"
	"github.com/opsee/basic/schema"
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

func (s *service) GetUser(ctx context.Context, req *opsee.GetUserRequest) (*opsee.GetUserResponse, error) {
	var (
		user *schema.User
		err  error
	)

	if req.CustomerId != "" {
		user, err = servicer.GetUserCustID(req.CustomerId)
		if err != nil {
			return nil, err
		}

	} else if req.Email != "" {
		user, err = servicer.GetUserEmail(req.Email)
		if err != nil {
			return nil, err
		}

	} else {
		user, err = servicer.GetUser(int(req.Id))
		if err != nil {
			return nil, err
		}
	}

	toke, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	return &opsee.GetUserResponse{
		User:       user,
		BasicToken: base64.StdEncoding.EncodeToString(toke),
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
