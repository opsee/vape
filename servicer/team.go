package servicer

import (
	"database/sql"

	"github.com/opsee/basic/schema"
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

// Gets subset of fields of a customer accessible to team admin
func GetTeamUsers(id string) ([]*schema.User, error) {
	users := []*schema.User{}
	err := store.Select(&users, "team-users-by-id", id)
	if err != nil && err != sql.ErrNoRows {
		log.WithError(err).Error("couldnt fetch team users")
		return nil, err
	}
	for _, user := range users {
		if user.Perms != nil {
			user.Perms.Name = "user"
		}
	}
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
