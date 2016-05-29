package service

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/opsee/basic/schema"
	opsee "github.com/opsee/basic/service"
	opsee_types "github.com/opsee/protobuf/opseeproto/types"
	"github.com/opsee/vape/servicer"
	"golang.org/x/net/context"
)

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

// Delete a user
// TODO(dan) This also needs to delete the users subscription.
func (s *service) DeleteUser(ctx context.Context, req *opsee.DeleteUserRequest) (*opsee.DeleteUserResponse, error) {
	err := servicer.DeleteUser(int(req.User.Id))
	if err != nil {
		return nil, err
	}

	return &opsee.DeleteUserResponse{
		User: req.User,
	}, nil
}

// Update user permissions
func (s *service) UpdateUserPerms(ctx context.Context, req *opsee.UpdateUserPermsRequest) (*opsee.UserTokenResponse, error) {
	token, err := servicer.UpdateUserPerms(req.User, req.Perms, time.Hour*24)
	if err != nil {
		return nil, err
	}

	return &opsee.UserTokenResponse{
		User:  req.User,
		Token: token,
	}, nil
}

// Update a user's email, name, or password
func (s *service) UpdateUser(ctx context.Context, req *opsee.UpdateUserRequest) (*opsee.UserTokenResponse, error) {
	user, err := servicer.GetUser(int(req.User.Id))
	if err != nil {
		return nil, err
	}

	if !req.Requestor.IsOpseeAdmin() && req.Requestor.CustomerId != user.CustomerId {
		return nil, opsee_types.NewPermissionsError("must be on same team")
	}

	// TODO(dan) move this into servicer.UpdateUser/MergeUser, update api, and test
	if req.Status != "" {
		user.Status = req.Status
	}

	token, err := servicer.UpdateUser(user, req.Email, req.Name, req.Password, time.Hour*24)
	if err != nil {
		return nil, err
	}

	updatedUser, err := servicer.GetUser(int(req.User.Id))
	if err != nil {
		return nil, err
	}

	return &opsee.UserTokenResponse{
		User: &schema.User{
			Id:     updatedUser.Id,
			Name:   updatedUser.Name,
			Email:  updatedUser.Email,
			Status: updatedUser.Status,
		},
		Token: token,
	}, nil
}
