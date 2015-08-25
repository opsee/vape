package store

import (
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
	Init(os.Getenv("TEST_POSTGRES_CONN"))

	// teardown first since it's nice to have lingering data to play with after a test
	_, err := db.Exec("delete from logins")
	if err != nil {
		c.Fatal(err)
	}
	_, err = db.Exec("delete from orgs")
	if err != nil {
		c.Fatal(err)
	}

	// fk constraint on customer_id
	_, err = db.Exec("insert into orgs (name, subdomain) values ('markorg', 'markorg')")
	if err != nil {
		c.Fatal(err)
	}
	_, err = db.Exec(
		"insert into logins (id, email, password_hash, admin, active, verified, " +
			"customer_id, name) values (1, 'mark@opsee.co', " +
			"'$2a$10$QcgjlXDKnRys50Oc30duFuNcZW6Rmqd7pcIJX9GWheIXJExUooZ7W', true, true, true, " +
			"'markorg', 'mark')")
	if err != nil {
		c.Fatal(err)
	}
}

func (s *UserSuite) TestGetUser(c *C) {
	user, err := GetUser("by-email-and-active", "mark@opsee.co", true)
	c.Assert(err, IsNil)
	c.Assert(user.Name, Equals, "mark")

	user, err = GetUser("by-email-and-active", "mark@opsee.co", false)
	c.Assert(err, NotNil)
}

func (s *UserSuite) TestAuthenticate(c *C) {
	user, err := GetUser("by-email-and-active", "mark@opsee.co", true)
	c.Assert(err, IsNil)

	err = user.Authenticate("shiteat")
	c.Assert(err, NotNil)

	err = user.Authenticate("eatshit")
	c.Assert(err, IsNil)
}
