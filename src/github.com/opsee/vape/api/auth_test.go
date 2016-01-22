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
	servicer.Init("test.opsy.co", nil, "fffff--fffffffffffffffffffffffffffffffff", "", "")
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

	// reglar ole login
	rec, err = testReq(publicRouter, "POST", "https://vape/authenticate/password", bytes.NewBuffer([]byte(`{"email": "mark@opsee.co", "password": "eatshit"}`)), nil)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Code, DeepEquals, 200)

	response := &UserTokenResponse{}
	loadResponse(response, rec.Body)
	c.Assert(response.IntercomHMAC, DeepEquals, "6b30a7417724ab8aaa918f4a66cc6be8f1ff278ab9d9ee3bcafcdd46326f1605")
	c.Assert(response.User.Email, DeepEquals, "mark@opsee.co")

	// admin login as user id 3
	rec, err = testReq(publicRouter, "POST", "https://vape/authenticate/password", bytes.NewBuffer([]byte(`{"email": "mark@opsee.co", "password": "eatshit", "as": 3}`)), nil)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Code, DeepEquals, 200)

	response = &UserTokenResponse{}
	loadResponse(response, rec.Body)
	c.Assert(response.User.Email, DeepEquals, "dan@opsee.co")
	c.Assert(response.User.AdminId, DeepEquals, 1)

	// non-admin shouldn't be able to log in as someone else
	rec, err = testReq(publicRouter, "POST", "https://vape/authenticate/password", bytes.NewBuffer([]byte(`{"email": "dan@opsee.co", "password": "eatshit", "as": 1}`)), nil)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Code, DeepEquals, 200)

	response = &UserTokenResponse{}
	loadResponse(response, rec.Body)
	c.Assert(response.User.Email, DeepEquals, "dan@opsee.co")
	c.Assert(response.User.AdminId, DeepEquals, 0)
}

func (s *ApiSuite) TestCreateAuthToken(c *C) {
	mailer := &testMailer{}
	servicer.Init("test.opsy.co", mailer, "fffff--fffffffffffffffffffffffffffffffff", "", "")

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
