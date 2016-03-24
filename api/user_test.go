package api

import (
	"bytes"
	"github.com/opsee/basic/schema"
	. "gopkg.in/check.v1"
)

func (s *ApiSuite) TestUserGet(c *C) {
	// a user viewing themselves
	rec, err := testAuthedReq(&schema.User{Id: 1, Email: "cliff@leaninto.it", Admin: false}, "GET",
		"https://vape/users/1", nil, nil)
	if err != nil {
		c.Fatal(err)
	}
	user := &schema.User{}
	err = loadResponse(user, rec.Body)
	c.Assert(user.Id, DeepEquals, int32(1))
	c.Assert(user.Email, DeepEquals, "mark@opsee.co")

	// non-admin viewing another
	rec, err = testAuthedReq(&schema.User{Id: 2, Email: "cliff@leaninto.it", Admin: false}, "GET",
		"https://vape/users/1", nil, nil)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Code, DeepEquals, 401)

	// admin viewing another
	rec, err = testAuthedReq(&schema.User{Id: 2, Email: "cliff@leaninto.it", Admin: true}, "GET",
		"https://vape/users/1", nil, nil)
	if err != nil {
		c.Fatal(err)
	}
	user = &schema.User{}
	err = loadResponse(user, rec.Body)
	c.Assert(user.Id, DeepEquals, int32(1))
	c.Assert(user.Email, DeepEquals, "mark@opsee.co")

	// not found
	rec, err = testAuthedReq(&schema.User{Id: 2, Email: "cliff@leaninto.it", Admin: true}, "GET",
		"https://vape/users/99", nil, nil)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Code, DeepEquals, 404)
}

func (s *ApiSuite) TestUserUpdate(c *C) {
	rec, err := testAuthedReq(&schema.User{Id: 1, Email: "cliff@leaninto.it", Admin: false}, "PUT",
		"https://vape/users/1", bytes.NewBuffer([]byte(`{"name": "vin diesel"}`)), nil)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(rec.Code, DeepEquals, 200)

	rec, err = testAuthedReq(&schema.User{Id: 1, Email: "cliff@leaninto.it", Admin: false}, "GET",
		"https://vape/users/1", nil, nil)
	if err != nil {
		c.Fatal(err)
	}

	user := &schema.User{}
	err = loadResponse(user, rec.Body)
	c.Assert(user.Name, DeepEquals, "vin diesel")
}

func (s *ApiSuite) TestGetListUsers(c *C) {
	// get list, fail
	rec, _ := testAuthedReq(&schema.User{Id: 1, Email: "cliff@leaninto.it", Admin: false}, "GET", "https://vape/users/", nil, nil)
	messageResponse := &MessageResponse{}
	loadResponse(messageResponse, rec.Body)
	c.Assert(messageResponse.Message, DeepEquals, Messages.UserOrAdminRequired)

	// get list
	rec, _ = testAuthedReq(&schema.User{Id: 1, Email: "cliff@leaninto.it", Admin: true}, "GET", "https://vape/users/", nil, nil)
	gotUsers := make([]*schema.User, 2)
	loadResponse(&gotUsers, rec.Body)
	c.Assert(gotUsers[len(gotUsers)-1].Name, Not(DeepEquals), "")
}
