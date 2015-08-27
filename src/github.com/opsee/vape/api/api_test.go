package api

import (
	"bytes"
	"github.com/opsee/vape/store"
	"github.com/opsee/vape/testutil"
	. "gopkg.in/check.v1"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type ApiSuite struct{}

var (
	_ = Suite(&ApiSuite{})
)

func Test(t *testing.T) { TestingT(t) }

func (s *ApiSuite) SetUpTest(c *C) {
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
