package handler

import (
	"encoding/json"
	"net/http"
	
	"abexercise/store"
)

// InventoryResponse holds the remaining quantity of the item
type InventoryResponse struct {
	Quantity int64 `json:"quantity"`
}

func InventoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(InventoryResponse {
		Quantity: store.Quantity.Load(),
	})
}
