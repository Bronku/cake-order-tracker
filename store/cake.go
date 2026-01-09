package store

import (
	"errors"
	"fmt"

	"github.com/Bronku/cake-order-tracker/models"
)

func (s *Store) cakeCount() int {
	out := 0
	rows, err := s.db.Query("select count(*) from cake;")
	if err != nil {
		return out
	}
	defer rows.Close()
	_ = rows.Next()
	_ = rows.Scan(&out)
	return out
}

func (s *Store) loadCakes() ([]models.Cake, error) {
	out := make([]models.Cake, 0, s.cakeCount())

	rows, err := s.db.Query("select id, name, price, category, availability from cake order by id;")
	if err != nil {
		return out, err
	}
	defer rows.Close()

	for rows.Next() {
		var cake models.Cake
		err = rows.Scan(&cake.ID, &cake.Name, &cake.Price, &cake.Category, &cake.Availability)
		if err != nil {
			continue
		}
		out = append(out, cake)
	}

	return out, nil
}

func (s *Store) searchCakes(id int) (int, error) {
	for i, e := range s.cakes {
		if e.ID != id {
			continue
		}
		return i, nil
	}
	return 0, errors.New("cake not found")
}

func (s *Store) GetCake(id int) (models.Cake, error) {
	if id <= 0 {
		return models.Cake{}, errors.New("invalid cake id")
	}

	i, err := s.searchCakes(id)
	if err != nil {
		return models.Cake{}, err
	}
	return s.cakes[i], err
}

func (s *Store) GetCakes() ([]models.Cake, error) {
	result := make([]models.Cake, len(s.cakes))
	copy(result, s.cakes)
	return result, nil
}

func (s *Store) updateCake(newCake models.Cake) error {
	query := "update cake set name = ? , price = ?, category = ?, availability = ?  where id = ?"
	_, err := s.db.Exec(query, newCake.Name, newCake.Price, newCake.Category, newCake.Availability, newCake.ID)
	if err != nil {
		return err
	}

	i, err := s.searchCakes(newCake.ID)
	if err != nil {
		return err
	}
	s.cakes[i] = newCake
	return nil
}

func (s *Store) SaveCake(newCake models.Cake) (int, error) {
	if newCake.ID != 0 {
		return newCake.ID, s.updateCake(newCake)
	}
	fmt.Println("adding a new cake", newCake)
	query := "insert into cake(name, price, category, availability) values (?, ?, ?, ?);"
	result, err := s.db.Exec(query, newCake.Name, newCake.Price, newCake.Category, newCake.Availability)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	newCake.ID = int(id)

	s.cakes = append(s.cakes, newCake)
	return newCake.ID, nil
}
