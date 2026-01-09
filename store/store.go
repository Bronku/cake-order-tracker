package store

import (
	"database/sql"
	"log"

	"github.com/Bronku/cake-order-tracker/models"
	_ "github.com/knaka/go-sqlite3-fts5"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db    *sql.DB
	cakes []models.Cake
	users map[string]models.User
}

func OpenStore(filename string) *Store {
	var out Store
	var err error

	out.db, err = sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal("Can't open the database", filename, err)
	}

	out.loadMigrations()
	out.cakes, err = out.loadCakes()
	log.Println(out.cakes)
	if err != nil {
		log.Fatal(err)
	}
	out.users, err = out.loadUsers()
	if err != nil {
		log.Fatal(err)
	}
	return &out
}

func (s *Store) Close() {
	if s.db != nil {
		s.db.Close()
	}
}
