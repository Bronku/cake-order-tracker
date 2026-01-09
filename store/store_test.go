package store

import (
	"testing"
	"time"

	"github.com/Bronku/cake-order-tracker/models"
)

func TestStore(t *testing.T) {
	s := OpenStore("file:memdb1?mode=memory&cache=shared")
	defer s.Close()

	// ensure db is empty
	cakes, err := s.GetCakes()
	if err != nil {
		t.Errorf("Error getting cakes: %v", err)
	}
	if len(cakes) != 0 {
		t.Errorf("Expected empty cakes list, got: %v", cakes)
	}
	orders, err := s.GetOrders(time.Now(), time.Now().Add(time.Hour*24))
	if err != nil {
		t.Errorf("Error getting orders: %v", err)
	}
	if len(orders) != 0 {
		t.Errorf("Expected empty orders list, got: %v", orders)
	}

	// create new cake
	newCake := models.Cake{Name: "Chocolate Cake", Price: 2500}
	newCake.ID, err = s.SaveCake(newCake)
	if err != nil {
		t.Fatalf("Failed to save new cake: %v", err)
	}
	if newCake.ID <= 0 {
		t.Errorf("Expected positive cake ID, got: %d", newCake.ID)
	}

	// create another cake
	anotherCake := models.Cake{Name: "Another Cake", Price: 3000}
	anotherCake.ID, err = s.SaveCake(anotherCake)
	if err != nil {
		t.Fatalf("Failed to save new cake: %v", err)
	}
	if anotherCake.ID <= 0 {
		t.Errorf("Expected positive cake ID, got: %d", anotherCake.ID)
	}

	// Update Existing Cake
	newCake.Name = "Updated Cake"
	newCake.Price = 100
	newID, err := s.SaveCake(newCake)
	if newID != newCake.ID {
		t.Error("wrong id")
	}
	if err != nil {
		t.Error("error updating cake", err)
	}

	// Update Non-existing Cake
	var updatedCake models.Cake
	updatedCake.ID = 10
	_, err = s.SaveCake(updatedCake)
	if err == nil {
		t.Error("did not return an error when attempted to update non existent cake")
	}

	// get cakes
	cakes, err = s.GetCakes()
	if err != nil {
		t.Fatalf("Failed to get all cakes after creation: %v", err)
	}
	if len(cakes) != 2 {
		t.Error("Expected two cakes returned, got ", len(cakes))
	}
	for i, c := range cakes {
		if c.ID == newCake.ID {
			newCake.Amount = i
		}
		if c.ID == anotherCake.ID {
			anotherCake.Amount = i
		}
	}
	if newCake.Amount == -1 {
		t.Error("not found newCake")
	}
	if cakes[newCake.Amount].ID != newCake.ID || cakes[newCake.Amount].Name != newCake.Name || cakes[newCake.Amount].Price != newCake.Price {
		t.Errorf("Want %v\nGot  %v", newCake, cakes[newCake.Amount])
	}
	if anotherCake.Amount == -1 {
		t.Error("not found anotherCake")
	}
	if cakes[anotherCake.Amount].ID != anotherCake.ID || cakes[anotherCake.Amount].Name != anotherCake.Name || cakes[anotherCake.Amount].Price != anotherCake.Price {
		t.Errorf("Want %v\nGot  %v", anotherCake, cakes)
	}

	selectedCake, err := s.GetCake(newCake.ID)
	if err != nil {
		t.Error("error getting a cake", err)
	}
	if selectedCake.ID != newCake.ID || selectedCake.Name != newCake.Name || selectedCake.Price != newCake.Price {
		t.Errorf("Want %v\nGot  %v", newCake, selectedCake)
	}

	selectedCake, err = s.GetCake(-1)
	if err == nil {
		t.Error("no error getting a non-existent cake")
	}

	// create new order
	now := time.Now()
	newOrder := models.Order{
		Name:     "John",
		Surname:  "Doe",
		Phone:    "123-456-7890",
		Location: "Some Location",
		Accepted: now,
		Date:     now.Add(time.Hour * 24),
		Status:   "Pending",
		Paid:     1000,
		Cakes:    []models.Cake{{ID: newCake.ID, Amount: 2}, {ID: anotherCake.ID, Amount: 10}},
	}
	newOrder.ID, err = s.SaveOrder(newOrder)
	if err != nil {
		t.Fatalf("Failed to save new order: %v", err)
	}
	if newOrder.ID <= 0 {
		t.Errorf("Expected positive order ID, got: %d", newOrder.ID)
	}

	// create another order
	anotherOrder := models.Order{
		Name:     "Jane",
		Surname:  "Doe",
		Phone:    "123-456-7890",
		Location: "New Location",
		Accepted: now,
		Date:     now.Add(time.Hour * 192),
		Status:   "Accepted",
		Paid:     1500,
		Cakes:    []models.Cake{{ID: newCake.ID, Amount: 100}},
	}
	anotherOrder.ID, err = s.SaveOrder(anotherOrder)
	if err != nil {
		t.Fatalf("Failed to save new order: %v", err)
	}
	if anotherOrder.ID <= 0 {
		t.Errorf("Expected positive order ID, got: %d", anotherOrder.ID)
	}

	// update existing order
	newOrder.Name = "James"
	newOrder.Status = "Done"
	newOrder.Cakes = newOrder.Cakes[1:]
	newID, err = s.SaveOrder(newOrder)
	if newID != newOrder.ID {
		t.Error("wrong id")
	}
	if err != nil {
		t.Error("error updating order", err)
	}

	// update non existing order
	var updatedOrder models.Order
	updatedOrder.ID = 10
	_, err = s.SaveOrder(updatedOrder)
	if err == nil {
		t.Error("did not return an error when attempted to update non existent cake")
	}

	// get orders
	orders, err = s.GetOrders(time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("Failed to get all orders after creation: %v", err)
	}
	if len(orders) != 2 {
		t.Error("Expected two orders returned, got ", len(orders))
	}
	newOrderPos := -1
	anotherOrderPos := -1
	for i, o := range orders {
		if o.ID == newOrder.ID {
			newOrderPos = i
		}
		if o.ID == anotherOrder.ID {
			anotherOrderPos = i
		}
	}
	if newOrderPos == -1 {
		t.Error("not found newOrder")
	}
	if newOrder.Name != orders[newOrderPos].Name ||
		newOrder.Surname != orders[newOrderPos].Surname ||
		newOrder.Phone != orders[newOrderPos].Phone ||
		newOrder.Status != orders[newOrderPos].Status ||
		newOrder.Location != orders[newOrderPos].Location ||
		newOrder.Accepted.Format("2006-01-02 15:04") != orders[newOrderPos].Accepted.Format("2006-01-02 15:04") ||
		newOrder.Date.Format("2006-01-02 15:04") != orders[newOrderPos].Date.Format("2006-01-02 15:04") {
		t.Errorf("Want %v\nGot  %v", newOrder, orders[newOrderPos])
	}
	if !areCakeSlicesEqual(newOrder.Cakes, orders[newOrderPos].Cakes) {
		t.Errorf("Want %v\nGot  %v", newOrder.Cakes, orders[newOrderPos].Cakes)
	}

	if anotherOrderPos == -1 {
		t.Error("not found anotherOrder")
	}
	if anotherOrder.Name != orders[anotherOrderPos].Name ||
		anotherOrder.Surname != orders[anotherOrderPos].Surname ||
		anotherOrder.Phone != orders[anotherOrderPos].Phone ||
		anotherOrder.Status != orders[anotherOrderPos].Status ||
		anotherOrder.Location != orders[anotherOrderPos].Location ||
		anotherOrder.Accepted.Format("2006-01-02 15:04") != orders[anotherOrderPos].Accepted.Format("2006-01-02 15:04") ||
		anotherOrder.Date.Format("2006-01-02 15:04") != orders[anotherOrderPos].Date.Format("2006-01-02 15:04") {
		t.Errorf("Want %v\nGot  %v", anotherOrder, orders[anotherOrderPos])
	}

	if !areCakeSlicesEqual(anotherOrder.Cakes, orders[anotherOrderPos].Cakes) {
		t.Errorf("Want %v\nGot  %v", anotherOrder.Cakes, orders[anotherOrderPos].Cakes)
	}

	selectedOrder, err := s.GetOrder(newOrder.ID)
	if err != nil {
		t.Error("error getting an order by id ", err)
	}
	if newOrder.Name != selectedOrder.Name ||
		newOrder.Surname != selectedOrder.Surname ||
		newOrder.Phone != selectedOrder.Phone ||
		newOrder.Status != selectedOrder.Status ||
		newOrder.Location != selectedOrder.Location ||
		newOrder.Accepted.Format("2006-01-02 15:04") != selectedOrder.Accepted.Format("2006-01-02 15:04") ||
		newOrder.Date.Format("2006-01-02 15:04") != selectedOrder.Date.Format("2006-01-02 15:04") {
		t.Errorf("Want %v\nGot  %v", newOrder, selectedOrder)
	}

	if !areCakeSlicesEqual(newOrder.Cakes, selectedOrder.Cakes) {
		t.Errorf("Want %v\nGot  %v", newOrder.Cakes, selectedOrder.Cakes)
	}

	_, err = s.GetOrder(-1)
	if err == nil {
		t.Error("no error getting invalid order")
	}
}

func TestSliceComparison(t *testing.T) {
	a := []models.Cake{
		{Name: "ok", Price: 100, ID: 12, Amount: 10},
		{Name: "ok", Price: 11, ID: 13, Amount: 1},
		{Name: "ok", Price: 100, ID: 10, Amount: 11},
	}

	b := []models.Cake{
		{Name: "okas", Price: 110, ID: 12, Amount: 10},
		{Name: "oksa", Price: 11, ID: 13, Amount: 1},
		{Name: "oasdk", Price: 101, ID: 10, Amount: 11},
	}
	if !areCakeSlicesEqual(a, b) {
		t.Error("aren't equal")
	}
}

func areCakesEqual(a, b models.Cake) bool {
	return a.ID == b.ID && a.Amount == b.Amount
}

func areCakeSlicesEqual(a, b []models.Cake) bool {
	if len(a) != len(b) {
		return false
	}

	matched := make([]bool, len(a))

mainLoop:
	for _, e1 := range a {
		for i, e2 := range b {
			if matched[i] {
				continue
			}
			if areCakesEqual(e1, e2) {
				matched[i] = true
				continue mainLoop
			}
		}
	}
	for _, e := range matched {
		if e == false {
			return false
		}
	}
	return true
}
