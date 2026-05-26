package main

import (
	"log"
	"net/http"
)

func main() {
	// creating examplar of service
	shortener := NewURLShortener()

	// setting roadmap
	http.HandleFunc("/shorten", ShortenHandler(shortener))
	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/", RedirectHandler(shortener))

	// starting server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
