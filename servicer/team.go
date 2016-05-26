package servicer

import (
	"database/sql"

	"github.com/opsee/basic/schema"
	"github.com/opsee/vape/store"
)

func MergeTeam(team *schema.Team, name, subscription string) error {
	if name != "" {
		team.Name = name
	}

	if subscription != "" {
		team.Subscription = subscription
	}

	return nil
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
	users := []*schema.User{}
	err = store.Get(team, "team-users-by-id", id)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err

		}
	}
	team.Users = users

	return team, nil
}

// Updates subset of fields of a customer accessible to team admin (name, subscription)
func UpdateTeam(team *schema.Team, name string, subscription string) (*schema.Team, error) {
	err := MergeTeam(team, name, subscription)
	if err != nil {
		return nil, err
	}

	_, err = store.NamedExec("update-team", team)
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
