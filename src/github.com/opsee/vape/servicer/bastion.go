package servicer

import (
        "github.com/opsee/vape/model"
        "github.com/opsee/vape/store"
)

func CreateBastion() (*model.Bastion, string, error) {
        bastion, plaintext, err := model.NewBastion()
        if err != nil {
                return nil, "", err
        }

        var bastionId string
        if err = store.Get(&bastionId, "insert-bastion", bastion.PasswordHash); err != nil {
                return nil, "", err
        }

        bastion.Id = bastionId
        return bastion, plaintext, nil
}

func AuthenticateBastion(id, password string) error {
        bastion := new(model.Bastion)
        if err := store.Get(bastion, "bastion-by-id", id); err != nil {
                return err
        }
        return bastion.Authenticate(password)
}
