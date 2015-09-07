package api

import (
	"fmt"
	"bytes"
	"time"
	. "gopkg.in/check.v1"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/servicer"
)

func (s *ApiSuite) TestCreateActivateClaimSignup(c *C) {
	mailer := &testMailer{}
	servicer.Init(mailer)

	badReqs := map[string]string{
		`{"email": "sackodonuts@hotmail.com"}`: Messages.NameRequired,
		`{"name": "sack o donuts"}`: Messages.EmailRequired,
	}

	for j, m := range badReqs {
		rec, _ := testReq(publicRouter, "POST", "https://vape/signups", bytes.NewBuffer([]byte(j)), nil)
		assertMessage(c, rec, m)
	}

	rec, _ := testReq(publicRouter, "POST", "https://vape/signups", bytes.NewBuffer([]byte(`{"email": "sackodonuts@hotmail.com", "name": "sack o donuts"}`)), nil)
	signup := &model.Signup{}
	loadResponse(signup, rec.Body)
	c.Assert(signup.Id, Not(DeepEquals), 0)
	time.Sleep(5*time.Millisecond) // wait for the goroutine to finish emailing, easier than passing a channel around somehow
	c.Assert(mailer.Template, DeepEquals, "signup-confirmation")

	rec, _ = testAuthedReq(&model.User{Id: 1, Email: "cliff@leaninto.it", Admin: true}, "PUT", "https://vape/signups/" + fmt.Sprint(signup.Id) + "/activate", nil, nil)
	c.Assert(rec.Code, DeepEquals, 200)
	time.Sleep(5*time.Millisecond)
	c.Assert(mailer.Template, DeepEquals, "beta-approval")
}

