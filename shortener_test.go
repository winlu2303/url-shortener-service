package main

import (
	"testing"
)

func TestURLShortener_Shorten(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		wantErr     error
		checkLength bool
	}{
		{
			name:        "valid HTTP URL",
			url:         "http://example.com",
			wantErr:     nil,
			checkLength: true,
		},
		{
			name:        "valid HTTPS URL",
			url:         "https://google.com/search?q=test",
			wantErr:     nil,
			checkLength: true,
		},
		{
			name:        "valid URL with port",
			url:         "https://localhost:8080/api/v1",
			wantErr:     nil,
			checkLength: true,
		},
		{
			name:        "empty line",
			url:         "",
			wantErr:     ErrEmptyURL,
			checkLength: false,
		},
		{
			name:        "ivalid URL - without schema",
			url:         "example.com",
			wantErr:     ErrInvalidURL,
			checkLength: false,
		},
		{
			name:        "invalid URL - only schema",
			url:         "https://",
			wantErr:     ErrInvalidURL,
			checkLength: false,
		},
		{
			name:        "invalid URL - unsupported schema",
			url:         "ftp://example.com",
			wantErr:     ErrInvalidURL,
			checkLength: false,
		},
	}

	shortener := NewURLShortener()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortID, err := shortener.Shorten(tt.url)

			// checking error
			if err != tt.wantErr {
				t.Errorf("Shorten() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// checking length of shortner ID
			if tt.checkLength && len(shortID) != shortener.urlLength {
				t.Errorf("Shorten() shortID length = %d, want %d", len(shortID), shortener.urlLength)
			}

			// checking that ID cost only from valid symbol
			if tt.checkLength && shortID == "" {
				t.Error("Shorten() returned empty shortID")
			}
		})
	}
}

func TestURLShortener_GetOriginal(t *testing.T) {
	tests := []struct {
		name      string
		setupURLs map[string]string // shortID -> originalURL
		shortID   string
		wantURL   string
		wantErr   error
	}{
		{
			name:      "existing short ID",
			setupURLs: map[string]string{"abc123": "https://example.com"},
			shortID:   "abc123",
			wantURL:   "https://example.com",
			wantErr:   nil,
		},
		{
			name:      "not existing short ID",
			setupURLs: map[string]string{},
			shortID:   "notexist",
			wantURL:   "",
			wantErr:   ErrURLNotFound,
		},
		{
			name:      "empty short ID",
			setupURLs: map[string]string{},
			shortID:   "",
			wantURL:   "",
			wantErr:   ErrEmptyURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortener := NewURLShortener()

			// to fill by testing data
			for shortID, originalURL := range tt.setupURLs {
				shortener.urls[shortID] = originalURL
				shortener.reverse[originalURL] = shortID
			}

			gotURL, err := shortener.GetOriginal(tt.shortID)

			if err != tt.wantErr {
				t.Errorf("GetOriginal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotURL != tt.wantURL {
				t.Errorf("GetOriginal() = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}

func TestURLShortener_DuplicateURL(t *testing.T) {
	shortener := NewURLShortener()

	originalURL := "https://example.com"

	// the first shortener
	shortID1, err := shortener.Shorten(originalURL)
	if err != nil {
		t.Fatalf("First Shorten() failed: %v", err)
	}

	// the second shortenet same URL
	shortID2, err := shortener.Shorten(originalURL)
	if err != nil {
		t.Fatalf("Second Shorten() failed: %v", err)
	}

	if shortID1 != shortID2 {
		t.Errorf("Duplicate URL returned different short IDs: %s vs %s", shortID1, shortID2)
	}
}

func TestGenerateShortID(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"length 6", 6},
		{"length 7", 7},
		{"length 8", 8},
		{"length 10", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := generateShortID(tt.length)

			if len(id) != tt.length {
				t.Errorf("generateShortID(%d) length = %d, want %d", tt.length, len(id), tt.length)
			}

			if id == "" {
				t.Error("generateShortID() returned empty string")
			}
		})
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{"valid HTTP", "http://example.com", true},
		{"valid HTTPS", "https://google.com", true},
		{"valid with port", "https://localhost:8080/api", true},
		{"valid with parameters", "https://example.com?q=test&page=1", true},
		{"empty line", "", false},
		{"without schema", "example.com", false},
		{"only schema", "https://", false},
		{"incorrect schem", "ftp://example.com", false},
		{"without host", "https://", false},
		{"only domen without schema", "google.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidURL(tt.url)
			if got != tt.want {
				t.Errorf("isValidURL(%q) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}
