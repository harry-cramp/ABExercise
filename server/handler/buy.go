package handler

import (
	"encoding/json"
	"net/http"
	
	"abexercise/store"
)

// BuyResponse contains the user's ticket number on purchase
type BuyResponse struct {
	Ticket int64 `json:"ticket"`
}

// BuyErrorResponse contains information in case of an error
type BuyErrorResponse struct {
	Error string `json:"error"`
	Message string `json:"message"`
}

func BuyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	if !store.AttemptBuy() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(BuyErrorResponse {
			Error: "out_of_stock",
			Message: "Sorry! Item is out of stock",
		})
		
		return
	}
	
	ticket := store.TicketNumber.Add(1)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(BuyResponse {
		Ticket: ticket,
	})
}
