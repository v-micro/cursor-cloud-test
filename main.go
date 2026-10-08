package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, newMux()); err != nil {
		log.Fatal(err)
	}
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("GET /users", userQueryHandler)
	mux.HandleFunc("HEAD /users", userQueryHandler)
	mux.HandleFunc("GET /users/{id}", userByIDHandler)
	mux.HandleFunc("HEAD /users/{id}", userByIDHandler)
	mux.HandleFunc("GET /orders", orderQueryHandler)
	mux.HandleFunc("HEAD /orders", orderQueryHandler)
	mux.HandleFunc("GET /orders/{id}", orderByIDHandler)
	mux.HandleFunc("HEAD /orders/{id}", orderByIDHandler)
	return mux
}
