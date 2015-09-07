package api

import (
	"bytes"
	. "gopkg.in/check.v1"
	"github.com/opsee/vape/model"
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
	c.Assert(user.Email, DeepEquals, "cliff@leaninto.it")
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
