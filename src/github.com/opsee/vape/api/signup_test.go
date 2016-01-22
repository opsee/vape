package api

import (
	"bytes"
	"fmt"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/servicer"
	. "gopkg.in/check.v1"
	"time"
)

func (s *ApiSuite) TestCreateActivateClaimSignup(c *C) {
	mailer := &testMailer{}
	servicer.Init("test.opsy.go", mailer, "fffff--fffffffffffffffffffffffffffffffff", "", "")

	badReqs := map[string]string{
		`{"email": "sackodonuts@hotmail.com"}`: Messages.NameRequired,
		`{"name": "sack o donuts"}`:            Messages.EmailRequired,
	}

	// testin bad requests
	for j, m := range badReqs {
		rec, _ := testReq(publicRouter, "POST", "https://vape/signups", bytes.NewBuffer([]byte(j)), nil)
		assertMessage(c, rec, m)
	}

	// create a signup
	rec, _ := testReq(publicRouter, "POST", "https://vape/signups", bytes.NewBuffer([]byte(`{"email": "sackodonuts@hotmail.com", "name": "sack o donuts"}`)), nil)
	signup := &model.Signup{}
	loadResponse(signup, rec.Body)
	c.Assert(signup.Id, Not(DeepEquals), 0)
	time.Sleep(5 * time.Millisecond) // wait for the goroutine to finish emailing, easier than passing a channel around somehow
	c.Assert(mailer.Template, DeepEquals, "signup-confirmation")

	// activate the signup by sending the user an email/token
	rec, _ = testAuthedReq(&model.User{Id: 1, Email: "cliff@leaninto.it", Admin: true}, "PUT", "https://vape/signups/"+fmt.Sprint(signup.Id)+"/activate", nil, nil)
	activationResponse := &SignupActivationResponse{}
	loadResponse(activationResponse, rec.Body)
	c.Assert(activationResponse.Token, Not(DeepEquals), "")
	time.Sleep(5 * time.Millisecond)
	c.Assert(mailer.Template, DeepEquals, "beta-approval")

	// claim the signup, turning it into a user
	rec, _ = testReq(publicRouter, "POST", "https://vape/signups/"+fmt.Sprint(signup.Id)+"/claim", bytes.NewBuffer([]byte(`{"token": "`+activationResponse.Token+`", "password": "sackodonuts"}`)), nil)
	userTokenResponse := &UserTokenResponse{}
	loadResponse(userTokenResponse, rec.Body)
	c.Assert(userTokenResponse.User.Id, Not(DeepEquals), 0)
	c.Assert(userTokenResponse.User.Name, DeepEquals, "sack o donuts")
}

func (s *ApiSuite) TestGetListSignups(c *C) {
	rec, _ := testReq(publicRouter, "POST", "https://vape/signups", bytes.NewBuffer([]byte(`{"email": "cliff@whitehouse.gov", "name": "president moon"}`)), nil)
	signup := &model.Signup{}
	loadResponse(signup, rec.Body)

	// this should fail
	rec, _ = testAuthedReq(&model.User{Id: 1, Email: "cliff@leaninto.it", Admin: false}, "GET", "https://vape/signups/"+fmt.Sprint(signup.Id), nil, nil)
	messageResponse := &MessageResponse{}
	loadResponse(messageResponse, rec.Body)
	c.Assert(messageResponse.Message, DeepEquals, Messages.AdminRequired)

	// get 1
	rec, _ = testAuthedReq(&model.User{Id: 1, Email: "cliff@leaninto.it", Admin: true}, "GET", "https://vape/signups/"+fmt.Sprint(signup.Id), nil, nil)
	gotSignup := &model.Signup{}
	loadResponse(gotSignup, rec.Body)
	c.Assert(gotSignup.Name, DeepEquals, signup.Name)

	// get list, fail
	rec, _ = testAuthedReq(&model.User{Id: 1, Email: "cliff@leaninto.it", Admin: false}, "GET", "https://vape/signups/", nil, nil)
	messageResponse = &MessageResponse{}
	loadResponse(messageResponse, rec.Body)
	c.Assert(messageResponse.Message, DeepEquals, Messages.AdminRequired)

	// get list
	rec, _ = testAuthedReq(&model.User{Id: 1, Email: "cliff@leaninto.it", Admin: true}, "GET", "https://vape/signups/", nil, nil)
	gotSignups := make([]*model.Signup, 2)
	loadResponse(&gotSignups, rec.Body)
	c.Assert(gotSignups[len(gotSignups)-1].Name, DeepEquals, signup.Name)
}
