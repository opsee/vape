package store

import (
	"os"
	"testing"

	"github.com/opsee/vape/model"
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
	user := new(model.User)
	err := Get(user, "user-by-email-and-active", "mark@opsee.co", true)
	c.Assert(err, IsNil)
	c.Assert(user.Name, Equals, "mark")

	err = Get(user, "user-by-email-and-active", "mark@opsee.co", false)
	c.Assert(err, NotNil)
}
