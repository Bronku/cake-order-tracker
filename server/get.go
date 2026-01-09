package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Bronku/cake-order-tracker/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func monthInterval(y int, m time.Month) (time.Time, time.Time) {
	firstDay := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	lastDay := time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
	return firstDay, lastDay
}

func (h *Server) orders(_ *http.Request) (any, int, error) {
	var (
		y int
		m time.Month
	)
	y, m, _ = time.Now().Date()
	first, last := monthInterval(y, m)
	var orders []models.Order
	result := h.db.
		Preload("OrderItems.Product").
		Preload(clause.Associations).
		Where("date between ? and ?", first, last).
		Find(&orders)
	if result.Error != nil {
		return nil, http.StatusInternalServerError, result.Error
	}
	data := struct {
		First  string
		Last   string
		Orders []models.Order
	}{first.Format("2006-01-02"), last.Format("2006-01-02"), orders}
	return data, http.StatusOK, nil
}

func (h *Server) ordersSearch(r *http.Request) (any, int, error) {
	startTime, err := time.Parse("2006-01-02", r.URL.Query().Get("from"))
	if err != nil {
		startTime = time.Time{}
	}
	endTime, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
	if err != nil {
		endTime = time.Time{}
	}
	var orders []models.Order
	result := h.db.
		Preload("OrderItems.Product").
		Preload(clause.Associations).
		Where("date between ? and ?", startTime, endTime).
		Find(&orders)
	return orders, http.StatusOK, result.Error
}

func (h *Server) cakes(_ *http.Request) (any, int, error) {
	var cakes []models.Product
	result := h.db.Find(&cakes)
	return cakes, http.StatusOK, result.Error
}

func (h *Server) cake(r *http.Request) (any, int, error) {
	url := strings.Split(r.URL.String(), "/")
	if len(url) < 3 || url[2] == "" {
		return models.Product{}, http.StatusOK, nil
	}

	cakeID, err := strconv.Atoi(url[2])
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("%w cakeID", ErrWrongValue)
	}

	var cake models.Product
	result := h.db.First(&cake, cakeID)
	if result.Error != nil {
		return models.Product{}, http.StatusNotFound, gorm.ErrRecordNotFound
	}
	return cake, http.StatusOK, nil
}

func (h *Server) order(r *http.Request) (any, int, error) {
	type formData struct {
		Order     models.Order
		Catalogue []models.Product
	}
	var data formData

	result := h.db.Find(&data.Catalogue)
	if result.Error != nil {
		return nil, http.StatusInternalServerError, ErrCatalogueNotFound
	}

	url := strings.Split(r.URL.String(), "/")
	if len(url) < 3 || url[2] == "" {
		return data, http.StatusOK, nil
	}

	orderID, err := strconv.Atoi(url[2])
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("%w orderID", ErrWrongValue)
	}

	result = h.db.Preload("OrderItems.Product").Preload(clause.Associations).Find(&data.Order, orderID)
	if result.Error != nil {
		return nil, http.StatusNotFound, gorm.ErrRecordNotFound
	}

	return data, http.StatusOK, nil
}
