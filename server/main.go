package main

import (
	"net/http"
	"log"

	"abexercise/handler"
	"abexercise/middleware"
	"abexercise/store"
)


func main() {
	store.Quantity.Store(10)
	store.TicketNumber.Store(0)
	middleware.IdemKeyMap = make(map[string]middleware.CachedResponse)

	mux := http.NewServeMux()

	mux.HandleFunc("/inventory", handler.InventoryHandler)
	mux.Handle("/buy", middleware.IdemMiddleware(http.HandlerFunc(handler.BuyHandler)))

	server := &http.Server{
		Addr:    ":8080",
		Handler: middleware.CorsMiddleware(mux),
	}

	log.Println("New Phones server running on localhost:8080")
	log.Fatal(server.ListenAndServe())
}
