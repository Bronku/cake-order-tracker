package auth

import (
	"errors"

	"github.com/Bronku/cake-order-tracker/crypto"
)

func (a *Authenticator) verifyCredentials(login, password string) error {
	user, ok := a.s.GetUser(login)
	if !ok {
		return errors.New("user with this login doesn't exist")
	}
	hash := crypto.PasswordHash(password, user.Salt)
	if hash == user.Password {
		return nil
	}
	return errors.New("wrong credentials")
}
