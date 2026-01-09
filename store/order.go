package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Bronku/cake-order-tracker/models"
)

// #todo update labels in templates

func (s *Store) parseOrderRow(row *sql.Rows) (models.Order, error) {
	var out models.Order
	var orderDate, deliveryDate string
	err := row.Scan(&out.ID, &out.Name, &out.Surname, &out.Phone, &out.Location, &orderDate, &deliveryDate, &out.Status, &out.Paid)
	if err != nil {
		return out, err
	}
	out.Accepted, _ = time.Parse("2006-01-02 15:04", orderDate)
	out.Date, _ = time.Parse("2006-01-02 15:04", deliveryDate)

	out.Cakes = make([]models.Cake, 0)
	rows, err := s.db.Query("select cake, amount from ordered_cake where customer_order = ?;", out.ID)
	if err != nil {
		return out, err
	}
	defer rows.Close()
	for rows.Next() {
		var newCake models.Cake
		err = rows.Scan(&newCake.ID, &newCake.Amount)
		if err != nil {
			return out, err
		}
		cakeData, err := s.GetCake(newCake.ID)
		if err != nil {
			continue
		}
		newCake.Name = cakeData.Name
		newCake.Price = cakeData.Price
		out.Cakes = append(out.Cakes, newCake)
	}

	return out, nil
}

func (s *Store) getOrdersFromQuery(query string, args ...any) ([]models.Order, error) {
	var out []models.Order
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return out, err
	}
	defer rows.Close()

	for rows.Next() {
		order, err := s.parseOrderRow(rows)
		if err != nil {
			fmt.Println(err)
			continue
		}
		out = append(out, order)
	}

	return out, nil
}

func (s *Store) GetOrder(id int) (models.Order, error) {
	query := "select * from customer_order where id = ?;"
	out, err := s.getOrdersFromQuery(query, id)
	if err != nil {
		return models.Order{}, err
	}
	if len(out) < 1 {
		return models.Order{}, errors.New("no output order")
	}
	return out[0], nil
}

func (s *Store) GetOrders(from, to time.Time) ([]models.Order, error) {
	start := from.Format("2006-01-02") + " 00:00"
	end := to.Format("2006-01-02") + " 99:99"
	if to.IsZero() {
		end = "9999-99-99 99:99"
	}
	query := "select * from customer_order where status != 'done' and delivery_date >= ? and delivery_date <= ? ;"
	return s.getOrdersFromQuery(query, start, end)
}

func (s *Store) UpdateOrderContents(tx *sql.Tx, newOrder models.Order) error {
	query := "delete from ordered_cake where customer_order = ?;"
	_, err := tx.Exec(query, newOrder.ID)
	if err != nil {
		return err
	}

	query = "insert into ordered_cake(customer_order, cake, amount) values (?,?,?);"
	for _, e := range newOrder.Cakes {
		_, err := tx.Exec(query, newOrder.ID, e.ID, e.Amount)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) UpdateOrder(newOrder models.Order) error {
	_, err := s.GetOrder(newOrder.ID)
	if err != nil {
		return err
	}
	query := "update customer_order set name = ?, surname = ?, phone = ?, location = ?, order_date = ?, delivery_date = ?, status = ?, paid = ? where id = ?;"

	accepted := newOrder.Accepted.Format("2006-01-02 15:04")
	date := newOrder.Date.Format("2006-01-02 15:04")

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, newOrder.Name, newOrder.Surname, newOrder.Phone, newOrder.Location, accepted, date, newOrder.Status, newOrder.Paid, newOrder.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = s.UpdateOrderContents(tx, newOrder)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	return nil
}

func (s *Store) SaveOrder(newOrder models.Order) (int, error) {
	if newOrder.ID != 0 {
		return newOrder.ID, s.UpdateOrder(newOrder)
	}
	query := "insert into customer_order(name, surname, phone, location, order_date, delivery_date, status, paid) values (?, ?, ?, ?, ?, ?, ?, ?);"

	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	accepted := newOrder.Accepted.Format("2006-01-02 15:04")
	date := newOrder.Date.Format("2006-01-02 15:04")
	result, err := tx.Exec(query, newOrder.Name, newOrder.Surname, newOrder.Phone, newOrder.Location, accepted, date, newOrder.Status, newOrder.Paid)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return newOrder.ID, err
	}
	newOrder.ID = int(id)

	err = s.UpdateOrderContents(tx, newOrder)
	if err != nil {
		_ = tx.Rollback()
		return newOrder.ID, err
	}

	err = tx.Commit()
	return newOrder.ID, err
}
