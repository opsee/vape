package api

import (
	"bytes"
	"fmt"
	"time"

	"github.com/opsee/basic/schema"
	opsee_types "github.com/opsee/protobuf/opseeproto/types"
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/servicer"
	"github.com/opsee/vape/testutil"
	. "gopkg.in/check.v1"
)

func (s *ApiSuite) TestCreateActivateClaimSignup(c *C) {
	adminPerms, _ := opsee_types.NewPermissions("user", "admin", "edit", "billing")

	mailer := &testMailer{}
	servicer.Init("test.opsy.co", mailer, "fffff--fffffffffffffffffffffffffffffffff", "", "", "", "", "", testutil.LDTestToken)

	badReqs := map[string]string{
		`{"name": "sack o donuts"}`: Messages.EmailRequired,
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
	time.Sleep(10 * time.Millisecond) // wait for the goroutine to finish emailing, easier than passing a channel around somehow
	c.Assert(mailer.Template, DeepEquals, "instant-approval")

	// try to create a signup by altering the email case errors
	rec, _ = testReq(publicRouter, "POST", "https://vape/signups", bytes.NewBuffer([]byte(`{"email": "sackodonuTS@hotmail.com", "name": "sack o donuts"}`)), nil)
	c.Assert(rec.Code, DeepEquals, 409)

	// activate the signup by sending the user an email/token
	rec, _ = testAuthedReq(&schema.User{Id: 1, Email: "cliff@leaninto.it", Admin: true, Perms: adminPerms}, "PUT", "https://vape/signups/"+fmt.Sprint(signup.Id)+"/activate", nil, nil)
	activationResponse := &SignupActivationResponse{}
	loadResponse(activationResponse, rec.Body)
	c.Assert(activationResponse.Token, Not(DeepEquals), "")
	time.Sleep(5 * time.Millisecond)
	c.Assert(mailer.Template, DeepEquals, "beta-approval")

	// claim the signup, turning it into a user
	rec, _ = testReq(publicRouter, "POST", "https://vape/signups/"+fmt.Sprint(signup.Id)+"/claim", bytes.NewBuffer([]byte(`{"token": "`+activationResponse.Token+`", "password": "sackodonuts", "name": "sack o donuts", "invite": true}`)), nil)
	userTokenResponse := &UserTokenResponse{}
	loadResponse(userTokenResponse, rec.Body)
	c.Assert(userTokenResponse.User.Id, Not(DeepEquals), 0)
	c.Assert(userTokenResponse.User.Name, DeepEquals, "sack o donuts")

	// test creating already activated signup - works for producthunt
	rec, _ = testReq(publicRouter, "POST", "https://vape/signups/new", bytes.NewBuffer([]byte(`{"email": "sackobanane@hotmail.com", "name": "sack o banane", "referrer": "producthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthuntproducthunt"}`)), nil)
	signup = &model.Signup{}
	loadResponse(signup, rec.Body)
	c.Assert(signup.Id, Not(DeepEquals), 0)
	c.Assert(signup.Activated, DeepEquals, true)
}

func (s *ApiSuite) TestGetListSignups(c *C) {
	nonAdminPerms, _ := opsee_types.NewPermissions("user", "edit", "billing")
	adminPerms, _ := opsee_types.NewPermissions("user", "admin", "edit", "billing")

	rec, _ := testReq(publicRouter, "POST", "https://vape/signups", bytes.NewBuffer([]byte(`{"email": "cliff@whitehouse.gov", "name": "president moon"}`)), nil)
	signup := &model.Signup{}
	loadResponse(signup, rec.Body)

	// this should fail
	rec, _ = testAuthedReq(&schema.User{Id: 1, Email: "cliff@leaninto.it", Admin: false, Perms: nonAdminPerms}, "GET", "https://vape/signups/"+fmt.Sprint(signup.Id), nil, nil)
	messageResponse := &MessageResponse{}
	loadResponse(messageResponse, rec.Body)
	c.Assert(messageResponse.Message, DeepEquals, Messages.AdminRequired)

	// get 1
	rec, _ = testAuthedReq(&schema.User{Id: 1, Email: "cliff@leaninto.it", Admin: true, Perms: adminPerms}, "GET", "https://vape/signups/"+fmt.Sprint(signup.Id), nil, nil)
	gotSignup := &model.Signup{}
	loadResponse(gotSignup, rec.Body)
	c.Assert(gotSignup.Name, DeepEquals, signup.Name)

	// get list, fail
	rec, _ = testAuthedReq(&schema.User{Id: 1, Email: "cliff@leaninto.it", Admin: false, Perms: nonAdminPerms}, "GET", "https://vape/signups/", nil, nil)
	messageResponse = &MessageResponse{}
	loadResponse(messageResponse, rec.Body)
	c.Assert(messageResponse.Message, DeepEquals, Messages.AdminRequired)

	// get list
	rec, _ = testAuthedReq(&schema.User{Id: 1, Email: "cliff@leaninto.it", Admin: true, Perms: adminPerms}, "GET", "https://vape/signups/", nil, nil)
	gotSignups := make([]*model.Signup, 2)
	loadResponse(&gotSignups, rec.Body)
	c.Assert(gotSignups[len(gotSignups)-1].Name, DeepEquals, signup.Name)
}
