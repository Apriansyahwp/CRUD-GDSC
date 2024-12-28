package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type Item struct {
	ItemID   int     `json:"item_id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type Purchase struct {
	ItemID     int     `json:"item_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}

var (
	items           []Item
	purchaseHistory []Purchase
	nextID          = 1
	itemsMu         sync.Mutex
)

func main() {
	http.HandleFunc("/items", itemsHandler)
	http.HandleFunc("/items/", itemHandler)
	http.HandleFunc("/purchase", purchaseHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server is running on http://localhost:8080/static/index.html")
	http.ListenAndServe(":8080", nil)
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getItems(w)
	case http.MethodPost:
		createOrUpdateItem(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/items/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getItemByID(w, id)
	case http.MethodPut:
		updateItem(w, r, id)
	case http.MethodDelete:
		deleteItem(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func purchaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var purchase Purchase
	if err := json.NewDecoder(r.Body).Decode(&purchase); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	itemsMu.Lock()
	defer itemsMu.Unlock()

	for i, item := range items {
		if item.ItemID == purchase.ItemID {
			if item.Quantity >= purchase.Quantity {
				totalPrice := item.Price * float64(purchase.Quantity)
				items[i].Quantity -= purchase.Quantity
				purchase.TotalPrice = totalPrice
				purchaseHistory = append(purchaseHistory, purchase)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(purchase)
				return
			}
			http.Error(w, "Insufficient stock", http.StatusBadRequest)
			return
		}
	}
	http.Error(w, "Item not found", http.StatusNotFound)
}

func getItems(w http.ResponseWriter) {
	itemsMu.Lock()
	defer itemsMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Items           []Item     `json:"items"`
		PurchaseHistory []Purchase `json:"purchase_history"`
	}{
		Items:           items,
		PurchaseHistory: purchaseHistory,
	}
	json.NewEncoder(w).Encode(response)
}

func createOrUpdateItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if item.Name == "" || item.Price <= 0 || item.Quantity < 0 {
		http.Error(w, "Invalid item data", http.StatusBadRequest)
		return
	}

	itemsMu.Lock()
	defer itemsMu.Unlock()

	for i, existingItem := range items {
		if existingItem.Name == item.Name {
			items[i].Quantity += item.Quantity
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(items[i])
			return
		}
	}

	item.ItemID = nextID
	nextID++
	items = append(items, item)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func getItemByID(w http.ResponseWriter, id int) {
	itemsMu.Lock()
	defer itemsMu.Unlock()

	for _, item := range items {
		if item.ItemID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Item not found", http.StatusNotFound)
}

func updateItem(w http.ResponseWriter, r *http.Request, id int) {
	var updatedItem Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	itemsMu.Lock()
	defer itemsMu.Unlock()

	for i, item := range items {
		if item.ItemID == id {
			items[i] = updatedItem
			items[i].ItemID = id
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(items[i])
			return
		}
	}
	http.Error(w, "Item not found", http.StatusNotFound)
}

func deleteItem(w http.ResponseWriter, id int) {
	itemsMu.Lock()
	defer itemsMu.Unlock()

	for i, item := range items {
		if item.ItemID == id {
			items = append(items[:i], items[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Item not found", http.StatusNotFound)
}
