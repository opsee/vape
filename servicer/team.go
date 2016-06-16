package servicer

import (
	"database/sql"

	"github.com/opsee/basic/schema"
	opsee_types "github.com/opsee/protobuf/opseeproto/types"
	"github.com/opsee/vape/store"
	log "github.com/sirupsen/logrus"
)

func MergeTeam(team *schema.Team, name, subscription string) error {
	if team.Name == "" {
		team.Name = "none"
	}
	if name != "" {
		team.Name = name
	}
	if subscription != "" {
		team.Subscription = subscription
	}

	return nil
}

func getTeamInvitedUsers(id string) ([]*schema.User, error) {
	var users []*schema.User
	signups, err := GetSignupsByCustomerId(id)
	if err != nil {
		return nil, err
	}
	for _, signup := range signups {
		if signup.Claimed == false {
			u := &schema.User{
				Id:         0, // meh
				CustomerId: id,
				Email:      signup.Email,
				Perms:      signup.Perms,
				Status:     "invited",
			}

			if u.Perms != nil {
				u.Perms.Name = "user"
			} else {
				u.Perms = &opsee_types.Permission{Name: "user", Perm: 0}
			}

			users = append(users, u)
		}
	}
	return users, nil
}

// Gets subset of fields of a customer accessible to team admin
func GetTeamUsers(id string) ([]*schema.User, error) {
	users := []*schema.User{}
	err := store.Select(&users, "team-users-by-id", id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	for _, user := range users {
		if user.Perms != nil {
			user.Perms.Name = "user"
		}
	}
	iu, err := getTeamInvitedUsers(id)
	if err != nil {
		log.WithError(err).Warn("could not get invited users. continuing")
	}
	users = append(users, iu...)
	return users, nil
}

// Gets subset of fields of a customer accessible to team admin
func GetTeam(id string) (*schema.Team, error) {
	team := new(schema.Team)
	err := store.Get(team, "team-by-id", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, TeamNotFound
		}
		return nil, err
	}

	team.Users, _ = GetTeamUsers(team.Id)

	return team, nil
}

// Updates subset of fields of a customer accessible to team admin (name, subscription)
func UpdateTeam(team *schema.Team, name string, subscription string) (*schema.Team, error) {
	currTeam, err := GetTeam(team.Id)
	if err != nil {
		return nil, TeamNotFound
	}

	err = MergeTeam(currTeam, name, subscription)
	if err != nil {
		return nil, err
	}

	_, err = store.NamedExec("update-team", currTeam)
	if err != nil {
		return nil, err
	}

	return GetTeam(team.Id)
}

// Sets customer to inactive
func DeleteTeam(id string) error {
	_, err := store.Exec("delete-team-by-id", id)
	return err
}
