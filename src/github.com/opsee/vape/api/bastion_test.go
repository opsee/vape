package api

import (
	"github.com/opsee/vape/store"
	"bytes"
	. "gopkg.in/check.v1"
)

func (s *ApiSuite) TestBastionAuthBadMissingCustomer(c *C) {
	rec, err := testReq(privateRouter, "POST", "https://vape/bastions", bytes.NewBuffer([]byte(`{"hey": "whatup"}`)), nil)
	if err != nil {
		c.Fatal(err)
	}

	assertMessage(c, rec, Messages.CustomerIdRequired)
}

func (s *ApiSuite) TestBastionAuthBadCustomer(c *C) {
	rec, err := testReq(privateRouter, "POST", "https://vape/bastions", bytes.NewBuffer([]byte(`{"customer_id": "4857879a-5363-11e5-81bc-4390140dd1a4"}`)), nil)
	if err != nil {
		c.Fatal(err)
	}

	assertMessage(c, rec, Messages.CustomerNotAuthorized)
}

func (s *ApiSuite) TestBastionAuth(c *C) {
	// get a valid customer id to create a bastion
	var customerId string
	err := store.DB.Get(&customerId, "select id from customers limit 1")
	if err != nil {
		c.Fatal(err)
	}

	// create a bastion auth
	rec, err := testReq(privateRouter, "POST", "https://vape/bastions", bytes.NewBuffer([]byte(`{"customer_id": "` + customerId + `"}`)), nil)
	if err != nil {
		c.Fatal(err)
	}

	bastion := &BastionResponse{}
	loadResponse(bastion, rec.Body)
	c.Assert(bastion.Password, NotNil)


	// test the auth
	rec, err = testReq(privateRouter, "POST", "https://vape/bastions/authenticate", bytes.NewBuffer([]byte(`{"id": "` + bastion.Id + `", "password": "` + bastion.Password + `"}`)), nil)
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(rec.Code, DeepEquals, 200)
}
