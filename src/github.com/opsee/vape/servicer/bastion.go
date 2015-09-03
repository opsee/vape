package servicer

import (
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/store"
)

func CreateBastion(customerId string) (*model.Bastion, string, error) {
	bastion, plaintext, err := model.NewBastion(customerId)
	if err != nil {
		return nil, "", err
	}

	// we'll want to activate the bastion as well
	bastion.Active = true

	err = store.NamedInsert("insert-bastion", bastion)
	if err != nil {
		return nil, "", err
	}

	return bastion, plaintext, nil
}

func AuthenticateBastion(id, password string) error {
	bastion := new(model.Bastion)
	if err := store.Get(bastion, "bastion-by-id-and-active", id, true); err != nil {
		return err
	}
	return bastion.Authenticate(password)
}
