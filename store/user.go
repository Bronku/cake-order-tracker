package store

import (
	"errors"

	"github.com/Bronku/cake-order-tracker/crypto"
	"github.com/Bronku/cake-order-tracker/models"
)

func (s *Store) AddUser(login, password string) error {
	_, exists := s.GetUser(login)
	if exists {
		return errors.New("the user already exists")
	}
	query := "insert into user (login, password, salt) values(?, ?, ?)"
	salt := crypto.GenerateKey()
	hash := crypto.PasswordHash(password, salt)
	_, err := s.db.Exec(query, login, hash, salt)
	if err == nil {
		s.users[login] = models.User{Password: hash, Salt: salt}
	}
	return err
}

func (s *Store) loadUsers() (map[string]models.User, error) {
	out := make(map[string]models.User)
	query := "select login, password, salt from user;"
	rows, err := s.db.Query(query)
	if err != nil {
		return out, err
	}
	defer rows.Close()
	for rows.Next() {
		var current models.User
		var login string
		err = rows.Scan(&login, &current.Password, &current.Salt)
		if err != nil {
			continue
		}
		out[login] = current
	}
	return out, nil
}

func (s *Store) GetUser(login string) (models.User, bool) {
	value, ok := s.users[login]
	return value, ok
}
