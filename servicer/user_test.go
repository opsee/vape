package servicer

import (
	"os"
	"testing"

	"github.com/opsee/basic/schema"
	"github.com/opsee/vape/store"
	"github.com/opsee/vape/testutil"
	. "gopkg.in/check.v1"
)

type UserSuite struct{}

var (
	_ = Suite(&UserSuite{})
)

func Test(t *testing.T) { TestingT(t) }

func (s *UserSuite) SetUpTest(c *C) {
	store.Init(os.Getenv("POSTGRES_CONN"))
	testutil.SetupFixtures(store.DB, c)
}

func (s *UserSuite) TestAuthenticate(c *C) {
	user := new(schema.User)
	err := store.Get(user, "user-by-email-and-active", "mark@opsee.co", true)
	c.Assert(err, IsNil)

	err = AuthenticateUser(user, "shiteat")
	c.Assert(err, NotNil)

	err = AuthenticateUser(user, "eatshit")
	c.Assert(err, IsNil)
}
