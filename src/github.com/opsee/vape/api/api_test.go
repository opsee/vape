package api

import (
	"bytes"
        "github.com/opsee/vape/token"
        "github.com/opsee/vape/model"
	"github.com/opsee/vape/store"
	"github.com/opsee/vape/testutil"
        "github.com/gocraft/web"
	. "gopkg.in/check.v1"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
        "encoding/json"
        "time"
)

type ApiSuite struct{}

var (
	_ = Suite(&ApiSuite{})
        testMiddlewareFunc = func(web.ResponseWriter, *web.Request) {}
        testVapeKey = []byte{194, 164, 235, 6, 138, 248, 171, 239, 24, 216, 11, 22, 137, 199, 215, 133}
)

func Test(t *testing.T) { TestingT(t) }

func (s *ApiSuite) SetUpTest(c *C) {
        token.Init(testVapeKey)
	store.Init(os.Getenv("TEST_POSTGRES_CONN"))
	testutil.SetupFixtures(store.DB, c)
}

func (s *ApiSuite) TestCors(c *C) {
	rec, err := testReq("POST", "https://vape/", nil, nil)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Header().Get("Access-Control-Allow-Origin"), DeepEquals, "")

	for _, o := range []string{"https://staging.opsy.co", "https://opsee.co"} {
		rec, err = testReq("POST", "https://vape/", nil, map[string]string{"Origin": o})
		if err != nil {
			c.Fatal(err)
		}
		c.Assert(rec.Header().Get("Access-Control-Allow-Origin"), DeepEquals, o)
	}

	rec, err = testReq("POST", "https://vape/", nil, map[string]string{"Origin": "https://zombo.com"})
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Header().Get("Access-Control-Allow-Origin"), DeepEquals, "")
}

func (s *ApiSuite) TestUserSessionEcho(c *C) {
	now := time.Now()
	exp := now.Add(time.Hour * 1)
	u := &model.User{
		Id:                 1,
		Email:              "cliff@leaninto.it",
		Admin:              true,
	}
        tok := token.New(u, "cliff@leaninto.it", now, exp)
	tokenString, err := tok.Marshal()
	if err != nil {
		c.Fatal(err)
	}

        auth := "Bearer " + tokenString
	rec, err := testReq("GET", "https://vape/authenticate/echo", nil, map[string]string{"Authorization": auth})
	if err != nil {
		c.Fatal(err)
	}

        user := &model.User{}
        dec := json.NewDecoder(rec.Body)
        err = dec.Decode(user)

        c.Assert(user.Id, DeepEquals, 1)
        c.Assert(user.Email, DeepEquals, "cliff@leaninto.it")
        c.Assert(user.Admin, DeepEquals, true)
}

func (s *ApiSuite) TestCreateAuthPassword(c *C) {
	rec, err := testReq("POST", "https://vape/authenticate/password", nil, nil)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Code, DeepEquals, 400)

	rec, err = testReq("POST", "https://vape/authenticate/password", bytes.NewBuffer([]byte(`{"email": "mark@opsee.co"}`)), nil)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Code, DeepEquals, 400)

	rec, err = testReq("POST", "https://vape/authenticate/password", bytes.NewBuffer([]byte(`{"email": "mark@opsee.co", "password": "hi"}`)), nil)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Code, DeepEquals, 401)
}

func testReq(method, url string, body io.Reader, headers map[string]string) (*httptest.ResponseRecorder, error) {
	if body == nil {
		body = bytes.NewBuffer([]byte{})
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec, nil
}
