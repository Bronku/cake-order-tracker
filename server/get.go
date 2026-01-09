package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Bronku/cake-order-tracker/models"
)

func monthInterval(y int, m time.Month) (firstDay, lastDay time.Time) {
	firstDay = time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	lastDay = time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
	return firstDay, lastDay
}

func (h *Server) orders(_ *http.Request) (any, int, error) {
	var (
		y int
		m time.Month
	)
	y, m, _ = time.Now().Date()
	first, last := monthInterval(y, m)
	fmt.Println(first, last)
	orders, err := h.s.GetOrders(first, last)
	data := struct {
		First  string
		Last   string
		Orders []models.Order
	}{first.Format("2006-01-02"), last.Format("2006-01-02"), orders}
	return data, http.StatusOK, err
}

func (h *Server) ordersSearch(r *http.Request) (any, int, error) {
	from, err := time.Parse("2006-01-02", r.URL.Query().Get("from"))
	if err != nil {
		from = time.Time{}
	}
	to, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
	if err != nil {
		to = time.Time{}
	}
	fmt.Println(from, to)
	data, err := h.s.GetOrders(from, to)
	return data, http.StatusOK, err
}

func (h *Server) cakes(_ *http.Request) (any, int, error) {
	data, err := h.s.GetCakes()
	return data, http.StatusOK, err
}

func (h *Server) cake(r *http.Request) (any, int, error) {
	url := strings.Split(r.URL.String(), "/")
	if len(url) < 3 || url[2] == "" {
		return models.Cake{}, http.StatusOK, nil
	}

	id, err := strconv.Atoi(url[2])
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	data, err := h.s.GetCake(id)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return data, http.StatusOK, nil
}

func (h *Server) order(r *http.Request) (any, int, error) {
	type formData struct {
		Order     models.Order
		Catalogue []models.Cake
	}
	var err error
	var data formData

	data.Catalogue, err = h.s.GetCakes()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	url := strings.Split(r.URL.String(), "/")
	if len(url) < 3 || url[2] == "" {
		return data, http.StatusOK, nil
	}

	id, err := strconv.Atoi(url[2])
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	data.Order, err = h.s.GetOrder(id)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return data, http.StatusOK, nil
}
