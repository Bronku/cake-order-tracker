package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Bronku/cake-order-tracker/models"
	"gorm.io/gorm"
)

func atoui(input string) (uint, error) {
	num, err := strconv.ParseUint(input, 10, 0)
	return uint(num), err
}

func (h *Server) postCake(r *http.Request) (any, int, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, http.StatusBadRequest, ErrInvalidForm
	}

	var postedCake models.Product
	postedCake.ID, err = atoui(r.FormValue("id"))
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("%w cakeID", ErrWrongValue)
	}
	postedCake.Name = r.FormValue("name")
	postedCake.Price, err = atoui(r.FormValue("price"))
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("%w price", ErrWrongValue)
	}
	result := h.db.Save(&postedCake)
	return postedCake, http.StatusAccepted, errors.Join(ErrSavingToDatabase, result.Error)
}

func (h *Server) parseOrder(r *http.Request) (models.Order, error) {
	var cakes []models.Product
	result := h.db.Find(&cakes)
	if result.Error != nil {
		return models.Order{}, ErrCatalogueNotFound
	}

	err := r.ParseForm()
	if err != nil {
		return models.Order{}, ErrInvalidForm
	}

	var postedOrder models.Order
	postedOrder.ID, err = atoui(r.FormValue("id"))
	if err != nil {
		return models.Order{}, fmt.Errorf("%w orderID", ErrWrongValue)
	}
	postedOrder.Prepaid, err = atoui(r.FormValue("paid"))
	if err != nil {
		return models.Order{}, fmt.Errorf("%w prepaid", ErrWrongValue)
	}
	postedOrder.Date, err = time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		return models.Order{}, fmt.Errorf("%w date", ErrWrongValue)
	}

	postedOrder.OrderItems = make([]models.OrderItem, 0)
	for _, e := range cakes {
		count, err := atoui(r.FormValue(fmt.Sprintf("cake[%d]", e.ID)))
		if err != nil {
			continue
		}
		postedOrder.OrderItems = append(postedOrder.OrderItems, models.OrderItem{Amount: count, Product: e})
	}

	postedOrder.Name = strings.TrimSpace(r.FormValue("name"))
	postedOrder.Surname = strings.TrimSpace(r.FormValue("surname"))
	postedOrder.Phone = strings.TrimSpace(r.FormValue("phone"))
	postedOrder.Location = strings.TrimSpace(r.FormValue("location"))
	postedOrder.Status = strings.TrimSpace(r.FormValue("status"))
	return postedOrder, nil
}

func (h *Server) postOrder(r *http.Request) (any, int, error) {
	postedOrder, err := h.parseOrder(r)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// i have no idea why this works, and other methods don't
	err = h.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&postedOrder).Error
		if err != nil {
			return fmt.Errorf("didn't save the order, %w", err)
		}
		err = tx.Model(&postedOrder).Association("OrderItems").Replace(postedOrder.OrderItems)
		if err != nil {
			return fmt.Errorf("didn't save the order items %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Join(ErrSavingToDatabase, err)
	}

	return postedOrder, http.StatusAccepted, nil
}
