package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/url"
	"sync"
)

// errors - calling errora as variables for testing
var (
	ErrEmptyURL     = errors.New("URL cannot be empty")
	ErrInvalidURL   = errors.New("invalid URL format")
	ErrURLNotFound  = errors.New("short URL not found")
	ErrDuplicateURL = errors.New("URL already exists")
)

// URLShortener - structure for saving URL
type URLShortener struct {
	urls      map[string]string
	reverse   map[string]string
	mu        sync.RWMutex
	urlLength int
}

// NewURLShortener creating new examplar URLShortener
func NewURLShortener() *URLShortener {
	return &URLShortener{
		urls:      make(map[string]string),
		reverse:   make(map[string]string),
		urlLength: 6,
	}
}

// Shorten is creating short ID for URL
func (us *URLShortener) Shorten(originalURL string) (string, error) {
	// 1. validation URL
	if originalURL == "" {
		return "", ErrEmptyURL
	}

	if !isValidURL(originalURL) {
		return "", ErrInvalidURL
	}

	us.mu.Lock()
	defer us.mu.Unlock()

	// 2. cheking on dublicate
	if shortID, exists := us.reverse[originalURL]; exists {
		return shortID, nil
	}

	// 3. to generate uniq and short ID
	shortID := generateShortID(us.urlLength)
	for {
		if _, exists := us.urls[shortID]; !exists {
			break
		}
		shortID = generateShortID(us.urlLength)
	}

	// 4. saving into maps
	us.urls[shortID] = originalURL
	us.reverse[originalURL] = shortID

	return shortID, nil
}

// GetOriginal return original URL by short ID
func (us *URLShortener) GetOriginal(shortID string) (string, error) {
	if shortID == "" {
		return "", ErrEmptyURL
	}

	us.mu.RLock()
	defer us.mu.RUnlock()

	originalURL, exists := us.urls[shortID]
	if !exists {
		return "", ErrURLNotFound
	}

	return originalURL, nil
}

// generateShortID generate random and short ID
func generateShortID(length int) string {
	// to generate random bytes
	bytes := make([]byte, length)
	rand.Read(bytes)

	// to code in base64 and cut by needed length
	encoded := base64.URLEncoding.EncodeToString(bytes)
	return encoded[:length]
}

// isValidURL checking is correct URL
func isValidURL(rawURL string) bool {
	if rawURL == "" {
		return false
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	// checking schema
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}

	// checking that host exist
	if parsed.Host == "" {
		return false
	}

	return true
}
