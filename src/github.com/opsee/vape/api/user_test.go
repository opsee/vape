package api

import (
	. "gopkg.in/check.v1"
	"github.com/opsee/vape/model"
)

func (s *ApiSuite) TestUserGet(c *C) {
	// a user viewing themselves
	rec, err := testAuthedReq(&model.User{Id: 1, Email: "cliff@leaninto.it", Admin: false}, "GET",
		"https://vape/users/1", nil, nil)
	if err != nil {
		c.Fatal(err)
	}
	user := &model.User{}
	err = loadResponse(user, rec.Body)
	c.Assert(user.Id, DeepEquals, 1)
	c.Assert(user.Email, DeepEquals, "mark@opsee.co")

	// non-admin viewing another
	rec, err = testAuthedReq(&model.User{Id: 2, Email: "cliff@leaninto.it", Admin: false}, "GET",
		"https://vape/users/1", nil, nil)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Code, DeepEquals, 401)

	// admin viewing another
	rec, err = testAuthedReq(&model.User{Id: 2, Email: "cliff@leaninto.it", Admin: true}, "GET",
		"https://vape/users/1", nil, nil)
	if err != nil {
		c.Fatal(err)
	}
	user = &model.User{}
	err = loadResponse(user, rec.Body)
	c.Assert(user.Id, DeepEquals, 1)
	c.Assert(user.Email, DeepEquals, "mark@opsee.co")

	// not found
	rec, err = testAuthedReq(&model.User{Id: 2, Email: "cliff@leaninto.it", Admin: true}, "GET",
		"https://vape/users/99", nil, nil)
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(rec.Code, DeepEquals, 404)
}
