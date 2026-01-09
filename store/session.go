package store

import (
	"time"

	"github.com/Bronku/cake-order-tracker/models"
)

func (s *Store) AddSession(token, userName string, expiration time.Time) error {
	query := "insert into session (token, user, expiration) values(?, ?, ?)"
	_, err := s.db.Exec(query, token, userName, expiration.Format("2006-01-02 15:04"))
	return err
}

func (s *Store) RevokeSession(token string) error {
	query := "delete from session where token = ?;"
	_, err := s.db.Exec(query, token)
	return err
}

func (s *Store) CleanSessions() error {
	now := time.Now().Format("2006-01-02 15:04")
	query := "delete from session where expiration < ?;"
	_, err := s.db.Exec(query, now)
	return err
}

func (s *Store) GetSessions() (map[string]models.Token, error) {
	out := make(map[string]models.Token)
	err := s.CleanSessions()
	if err != nil {
		return out, err
	}
	query := "select token, user, expiration from session;"
	rows, err := s.db.Query(query)
	if err != nil {
		return out, err
	}
	defer rows.Close()
	for rows.Next() {
		var current models.Token
		var token string
		var expiration string
		err = rows.Scan(&token, &current.User, &expiration)
		if err != nil {
			continue
		}
		current.Expiration, err = time.Parse("2006-01-02 15:04", expiration)
		if err != nil {
			continue
		}
		if token == "" {
			continue
		}
		out[token] = current
	}
	return out, nil
}
