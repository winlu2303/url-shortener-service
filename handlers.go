package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ShortenRequest - structer request for /shorten
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse - structure response for /shorten
type ShortenResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// ErrorResponse - structure for mistakes
type ErrorResponse struct {
	Error string `json:"error"`
}

// ShortenHandler handler POST /shorten
func ShortenHandler(shortener *URLShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only POST method
		if r.Method != http.MethodPost {
			sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// decoding JSON
		var req ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendJSONError(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		// to short URL
		shortID, err := shortener.Shorten(req.URL)
		if err != nil {
			switch err {
			case ErrEmptyURL:
				sendJSONError(w, "URL cannot be empty", http.StatusBadRequest)
			case ErrInvalidURL:
				sendJSONError(w, "Invalid URL format", http.StatusBadRequest)
			default:
				sendJSONError(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// forming full and short URL
		shortURL := buildShortURL(r, shortID)

		// sending response
		resp := ShortenResponse{
			ShortURL:    shortURL,
			OriginalURL: req.URL,
		}
		sendJSONResponse(w, resp, http.StatusCreated)
	}
}

// RedirectHandler handler GET /{short_url}
func RedirectHandler(shortener *URLShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only GET method
		if r.Method != http.MethodGet {
			sendJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// get shortID from way
		shortID := strings.TrimPrefix(r.URL.Path, "/")
		if shortID == "" {
			sendJSONError(w, "Short URL not provided", http.StatusBadRequest)
			return
		}

		// get original URL
		originalURL, err := shortener.GetOriginal(shortID)
		if err != nil {
			if err == ErrURLNotFound {
				sendJSONError(w, "Short URL not found", http.StatusNotFound)
			} else {
				sendJSONError(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// redirect
		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}

// HealthHandler check workability service
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	sendJSONResponse(w, map[string]string{"status": "ok"}, http.StatusOK)
}

// sendJSONResponse send JSON response
func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// sendJSONError send JSON mistake
func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	sendJSONResponse(w, ErrorResponse{Error: message}, statusCode)
}

// buildShortURL build full and short URL
func buildShortURL(r *http.Request, shortID string) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host + "/" + shortID
}
