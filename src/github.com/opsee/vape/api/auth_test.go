package api

import (
	"bytes"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/servicer"
	. "gopkg.in/check.v1"
	"time"
)

func (s *ApiSuite) TestUserSessionEcho(c *C) {
	rec, err := testAuthedReq(&model.User{Id: 1, Email: "cliff@leaninto.it", Admin: true}, "GET",
		"https://vape/authenticate/echo", nil, nil)
	if err != nil {
		c.Fatal(err)
	}

	user := &model.User{}
	err = loadResponse(user, rec.Body)
	c.Assert(user.Id, DeepEquals, 1)
	c.Assert(user.Admin, DeepEquals, true)
}

func (s *ApiSuite) TestCreateAuthPassword(c *C) {
	rec, err := testReq(publicRouter, "POST", "https://vape/authenticate/password", nil, nil)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Code, DeepEquals, 400)

	rec, err = testReq(publicRouter, "POST", "https://vape/authenticate/password", bytes.NewBuffer([]byte(`{"email": "mark@opsee.co"}`)), nil)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Code, DeepEquals, 400)

	rec, err = testReq(publicRouter, "POST", "https://vape/authenticate/password", bytes.NewBuffer([]byte(`{"email": "mark@opsee.co", "password": "hi"}`)), nil)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Code, DeepEquals, 401)
}

func (s *ApiSuite) TestCreateAuthToken(c *C) {
	mailer := &testMailer{}
	servicer.Init("test.opsy.co", mailer)

	// test a non-existent email
	rec, _ := testReq(publicRouter, "POST", "https://vape/authenticate/token", bytes.NewBuffer([]byte(`{"email": "what@rudoing.com"}`)), nil)
	c.Assert(rec.Code, DeepEquals, 401)

	// ok, this is a real user
	rec, _ = testReq(publicRouter, "POST", "https://vape/authenticate/token", bytes.NewBuffer([]byte(`{"email": "mark@opsee.co"}`)), nil)
	messageResponse := &MessageResponse{}
	loadResponse(messageResponse, rec.Body)
	c.Assert(messageResponse.Message, DeepEquals, Messages.Ok)

	// look for our token email
	time.Sleep(5 * time.Millisecond) // wait for the goroutine to finish emailing, easier than passing a channel around somehow
	c.Assert(mailer.Template, DeepEquals, "password-reset")
}
