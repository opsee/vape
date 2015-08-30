package servicer

import (
	"github.com/opsee/vape/model"
	"github.com/opsee/vape/store"
)

func CreateBastion(orgId int) (*model.Bastion, string, error) {
	bastion, plaintext, err := model.NewBastion(orgId)
	if err != nil {
		return nil, "", err
	}

	// we'll want to activate the bastion as well
	bastion.Active = true

	// need to pull out the generated bastion id, so use a query instead
	rows, err := store.NamedQuery("insert-bastion", bastion)
	if err != nil {
		return nil, "", err
	}
	for rows.Next() {
		if err = rows.StructScan(bastion); err != nil {
			return nil, "", err
		}
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
