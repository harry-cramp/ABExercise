package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"abexercise/handler"
	"abexercise/middleware"
	"abexercise/store"
)

func resetState() {
	store.Quantity.Store(1)
	store.TicketNumber.Store(0)
	middleware.IdemKeyMap = make(map[string]middleware.CachedResponse)
}

func TestInventoryHandler_OK(t *testing.T) {
	resetState()

	expectedQuantity := int64(1)

	req := httptest.NewRequest(http.MethodGet, "/inventory", nil)
	rec := httptest.NewRecorder()

	handler.InventoryHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}

	var body handler.InventoryResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body.Quantity != expectedQuantity {
		t.Fatalf("expected quantity %d, got %d", expectedQuantity, body.Quantity)
	}
}

func TestBuyHandler_OK(t *testing.T) {
	resetState()

	expectedTicket := int64(1)

	req := httptest.NewRequest(http.MethodPost, "/buy", nil)
	rec := httptest.NewRecorder()

	handler.BuyHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}

	var body handler.BuyResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body.Ticket != expectedTicket {
		t.Fatalf("expected quantity %d, got %d", expectedTicket, body.Ticket)
	}
}

func TestBuyHandler_Idempotent(t *testing.T) {
	resetState()

	buyHandler := middleware.IdemMiddleware(http.HandlerFunc(handler.BuyHandler))
	key := "test-idem-key"

	// first request
	req1 := httptest.NewRequest(http.MethodPost, "/buy", nil)
	req1.Header.Set("Idempotency-Key", key)
	rec1 := httptest.NewRecorder()

	buyHandler.ServeHTTP(rec1, req1)

	res1 := rec1.Result()
	defer res1.Body.Close()

	if res1.StatusCode != http.StatusOK {
		t.Fatalf("first request: expected 200, got %d", res1.StatusCode)
	}

	// second request with same idempotency key
	req2 := httptest.NewRequest(http.MethodPost, "/buy", nil)
	req2.Header.Set("Idempotency-Key", key)
	rec2 := httptest.NewRecorder()

	buyHandler.ServeHTTP(rec2, req2)

	res2 := rec2.Result()
	defer res2.Body.Close()

	if res2.StatusCode != http.StatusOK {
		t.Fatalf("second request: expected 200, got %d", res2.StatusCode)
	}

	// internal server state is the same after second request
	if store.GetQuantity() != int64(0) {
		t.Fatalf("expected quantity to be 0, got %d", store.GetQuantity())
	}

	if store.GetTicketNumber() != int64(1) {
		t.Fatalf("expected ticketNumber 1, got %d", store.GetTicketNumber())
	}
}

func TestBuyHandler_OutOfStock(t *testing.T) {
	store.Quantity.Store(0)

	req := httptest.NewRequest(http.MethodPost, "/buy", nil)
	rec := httptest.NewRecorder()

	handler.BuyHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, res.StatusCode)
	}

	var body handler.BuyErrorResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	
	if body.Message != "Sorry! Item is out of stock" {
		t.Fatalf("unexpected error message when buying out of stock item: %s", body.Message)
	}
}
