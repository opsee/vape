package store

import (
        "os"
        "testing"
        "time"
        "encoding/json"
        "github.com/dvsekhvalnov/jose2go"
        . "gopkg.in/check.v1"
)

type UserSuite struct{}

var (
        testVapeKey = []byte{194,164,235,6,138,248,171,239,24,216,11,22,137,199,215,133}
        _ = Suite(&UserSuite{})
)

func Test(t *testing.T) { TestingT(t) }

func (s *UserSuite) SetUpTest(c *C) {
        Init(os.Getenv("TEST_POSTGRES_CONN"), testVapeKey)

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

func (s *UserSuite) TestAuthenticateUser(c *C) {
        user, err := AuthenticateUser("mark@opsee.co", "eatshit")
        c.Assert(err, IsNil)
        c.Assert(user.Name, Equals, "mark")

        user, err = AuthenticateUser("mark@opsee.co", "shiteat")
        c.Assert(err, NotNil)

        user, err = AuthenticateUser("mark@opsee.com", "eatshit")
        c.Assert(err, NotNil)
}

func (s *UserSuite) TestMarshalJwe(c *C) {
        user, err := AuthenticateUser("mark@opsee.co", "eatshit")
        c.Assert(err, IsNil)

        token, err := user.MarshalJwe()
        c.Assert(err, IsNil)
        c.Assert(token, Not(DeepEquals), "")

        // ok, let's try decoding ok
        payload, headers, err := jose.Decode(token, testVapeKey)
        c.Assert(err, IsNil)

        // we should make these assertions in a decoding library (just a reminder!)
        c.Assert(headers["alg"], DeepEquals, "A128GCMKW")
        c.Assert(headers["enc"], DeepEquals, "A128GCM")

        payloadJson := make(map[string]interface{})
        err = json.Unmarshal([]byte(payload), &payloadJson)
        if err != nil {
                c.Fatal(err)
        }

        // and also importantly, we should make the exp assertion in a decoding library
        now := time.Now()
        exp := time.Unix(int64(payloadJson["exp"].(float64)), 0)
        c.Assert(now.Before(exp), DeepEquals, true)

        // just sanity that we got data in the payload
        c.Assert(payloadJson["id"].(float64), DeepEquals, float64(1))
        c.Assert(payloadJson["email"].(string), DeepEquals, "mark@opsee.co")
}
