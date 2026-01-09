package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Bronku/cake-order-tracker/models"
)

func (h *Server) postCake(r *http.Request) (any, int, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	var n models.Cake
	n.ID, err = strconv.Atoi(r.FormValue("id"))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	n.Name = r.FormValue("name")
	n.Price, err = strconv.Atoi(r.FormValue("price"))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	n.Category = r.FormValue("category")
	n.Availability = r.FormValue("availability")
	n.ID, err = h.s.SaveCake(n)
	fmt.Println(err)
	return n, http.StatusAccepted, err
}

func (h *Server) postOrder(r *http.Request) (any, int, error) {
	cakes, err := h.s.GetCakes()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	err = r.ParseForm()
	fmt.Println(r.PostForm)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var n models.Order
	n.ID, err = strconv.Atoi(r.FormValue("id"))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	n.Paid, err = strconv.Atoi(r.FormValue("paid"))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	n.Date, err = time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	n.Cakes = make([]models.Cake, 0)
	for _, e := range cakes {
		e.Amount, err = strconv.Atoi(r.FormValue(fmt.Sprintf("cake[%d]", e.ID)))
		if err != nil {
			continue
		}
		n.Cakes = append(n.Cakes, e)
	}

	n.Accepted = time.Now()
	n.Name = strings.TrimSpace(r.FormValue("name"))
	n.Surname = strings.TrimSpace(r.FormValue("surname"))
	n.Phone = strings.TrimSpace(r.FormValue("phone"))
	n.Location = strings.TrimSpace(r.FormValue("location"))
	n.Status = strings.TrimSpace(r.FormValue("status"))

	n.ID, err = h.s.SaveOrder(n)
	fmt.Println(n)
	return n, http.StatusAccepted, err
}
