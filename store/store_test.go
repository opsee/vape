package store

import (
	"os"
	"testing"

	"github.com/opsee/basic/schema"
	"github.com/opsee/vape/testutil"
	. "gopkg.in/check.v1"
)

type StoreSuite struct{}

var (
	_ = Suite(&StoreSuite{})
)

func Test(t *testing.T) { TestingT(t) }

func (s *StoreSuite) SetUpTest(c *C) {
	Init(os.Getenv("POSTGRES_CONN"))
	testutil.SetupFixtures(DB, c)
}

func (s *StoreSuite) TestGetUser(c *C) {
	user := new(schema.User)
	err := Get(user, "user-by-email-and-active", "mark@opsee.co", true)
	c.Assert(err, IsNil)
	c.Assert(user.Name, Equals, "mark")
}

func (s *StoreSuite) TestTeam(c *C) {
	var markTeam = &schema.Team{
		Name:         "MarkTeam",
		Subscription: "free",
	}

	user := new(schema.User)
	err := Get(user, "user-by-email-and-active", "mark@opsee.co", true)
	c.Assert(err, IsNil)
	c.Assert(user.Name, Equals, "mark")
	markTeam.Id = user.CustomerId

	// update-team
	_, err = NamedExec("update-team", markTeam)
	c.Assert(err, IsNil)

	// team-by-id
	team := &schema.Team{}
	err = Get(team, "team-by-id", markTeam.Id)
	c.Assert(err, IsNil)
	c.Assert(team.Name, Equals, markTeam.Name)
	c.Assert(team.Subscription, Equals, "free")

	// team-by-name
	teambyname := new(schema.Team)
	err = Get(teambyname, "team-by-name", markTeam.Name)
	c.Assert(err, IsNil)
	c.Assert(teambyname.Name, Equals, markTeam.Name)
	c.Assert(teambyname.Subscription, Equals, "free")

	// team-users-by-id
	users := []*schema.User{}
	err = Select(&users, "team-users-by-id", markTeam.Id)
	c.Assert(err, IsNil)
	c.Assert(len(users), Equals, 1)
}
