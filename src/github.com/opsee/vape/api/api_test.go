package api

import (
	"bytes"
	"encoding/json"
	"github.com/gocraft/web"
	"github.com/keighl/mandrill"
	"github.com/opsee/basic/schema"
	"github.com/opsee/vape/store"
	"github.com/opsee/vape/testutil"
	"github.com/opsee/vaper"
	. "gopkg.in/check.v1"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type ApiSuite struct{}

var (
	_                  = Suite(&ApiSuite{})
	testMiddlewareFunc = func(web.ResponseWriter, *web.Request) {}
	testVapeKey        = []byte{194, 164, 235, 6, 138, 248, 171, 239, 24, 216, 11, 22, 137, 199, 215, 133}
)

func Test(t *testing.T) { TestingT(t) }

func (s *ApiSuite) SetUpTest(c *C) {
	vaper.Init(testVapeKey)
	store.Init(os.Getenv("POSTGRES_CONN"))
	testutil.SetupFixtures(store.DB, c)
	// InjectLogger(os.Stdout)
}

type testMailer struct {
	Message  *mandrill.Message
	Template string
	Content  interface{}
}

func (t *testMailer) MessagesSendTemplate(msg *mandrill.Message, templateName string, templateContent interface{}) ([]*mandrill.Response, error) {
	t.Message = msg
	t.Template = templateName
	t.Content = templateContent
	return nil, nil
}

func (s *ApiSuite) TestCors(c *C) {
	rec, err := testReq(publicRouter, "POST", "https://vape/", nil, nil)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Header().Get("Access-Control-Allow-Origin"), DeepEquals, "")

	for _, o := range []string{"https://staging.opsy.co", "https://app.opsee.co", "https://app.opsee.com"} {
		rec, err = testReq(publicRouter, "POST", "https://vape/", nil, map[string]string{"Origin": o})
		if err != nil {
			c.Fatal(err)
		}
		c.Assert(rec.Header().Get("Access-Control-Allow-Origin"), DeepEquals, o)
	}

	rec, err = testReq(publicRouter, "POST", "https://vape/", nil, map[string]string{"Origin": "https://zombo.com"})
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Header().Get("Access-Control-Allow-Origin"), DeepEquals, "")
}

func testReq(router *web.Router, method, url string, body io.Reader, headers map[string]string) (*httptest.ResponseRecorder, error) {
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

func testAuthedReq(u *schema.User, method, url string, body io.Reader, headers map[string]string) (*httptest.ResponseRecorder, error) {
	if headers == nil {
		headers = make(map[string]string)
	}

	now := time.Now()
	exp := now.Add(time.Hour * 1)
	tok := vaper.New(u, u.Email, now, exp)
	tokenString, err := tok.Marshal()
	if err != nil {
		return nil, err
	}

	auth := "Bearer " + tokenString
	headers["Authorization"] = auth

	return testReq(publicRouter, method, url, body, headers)
}

func loadResponse(thing interface{}, body io.Reader) error {
	dec := json.NewDecoder(body)
	err := dec.Decode(thing)
	if err != nil {
		return err
	}
	return nil
}

func assertMessage(c *C, rec *httptest.ResponseRecorder, msg string) {
	resp := &MessageResponse{}
	loadResponse(resp, rec.Body)
	c.Assert(msg, DeepEquals, resp.Message)
}
