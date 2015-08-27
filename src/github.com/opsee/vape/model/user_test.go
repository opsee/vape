package model

import (
	"github.com/opsee/vape/store"
	"github.com/opsee/vape/testutil"
	. "gopkg.in/check.v1"
	"os"
	"testing"
)

type UserSuite struct{}

var (
	_ = Suite(&UserSuite{})
)

func Test(t *testing.T) { TestingT(t) }

func (s *UserSuite) SetUpTest(c *C) {
	store.Init(os.Getenv("TEST_POSTGRES_CONN"))
	testutil.SetupFixtures(store.DB, c)
}

func (s *UserSuite) TestAuthenticate(c *C) {
	user := new(User)
	err := store.Get(user, "user-by-email-and-active", "mark@opsee.co", true)
	c.Assert(err, IsNil)

	err = user.Authenticate("shiteat")
	c.Assert(err, NotNil)

	err = user.Authenticate("eatshit")
	c.Assert(err, IsNil)
}
