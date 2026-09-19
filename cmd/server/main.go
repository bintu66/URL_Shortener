package main

import (
	"fmt"
	"net/http"
	"url_shortener/internal/handler"
	"url_shortener/internal/service"
	"url_shortener/internal/storage"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Server is running OK!")
}

func main() {

	store := storage.NewMemoryStore()

	urlService := service.NewURLService(store)

	urlHandler := handler.NewURLHandler(urlService)

	// using http.NewServeMux() instead of the default one because it supports
	// method-based routing (Go 1.22+). This means we can say "only accept GET"
	// or "only accept POST" for each route, and avoid conflicts.
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", urlHandler.HomePage)       
	mux.HandleFunc("GET /health", healthHandler)           
	mux.HandleFunc("POST /shorten", urlHandler.ShortenURL) 
	mux.HandleFunc("GET /{code}", urlHandler.RedirectURL)

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}