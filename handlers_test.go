package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShortenHandler(t *testing.T) {
	shortener := NewURLShortener()

	tests := []struct {
		name           string
		method         string
		body           interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "succesful shortner URL",
			method:         http.MethodPost,
			body:           ShortenRequest{URL: "https://example.com"},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name:           "incorrect method",
			method:         http.MethodGet,
			body:           nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Method not allowed",
		},
		{
			name:           "empty URL",
			method:         http.MethodPost,
			body:           ShortenRequest{URL: ""},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "URL cannot be empty",
		},
		{
			name:           "invalid URL",
			method:         http.MethodPost,
			body:           ShortenRequest{URL: "not-a-url"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid URL format",
		},
		{
			name:           "invalid JSON",
			method:         http.MethodPost,
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			switch v := tt.body.(type) {
			case ShortenRequest:
				reqBody, err = json.Marshal(v)
				if err != nil {
					t.Fatalf("Failed to marshal request: %v", err)
				}
			case string:
				reqBody = []byte(v)
			}

			req := httptest.NewRequest(tt.method, "/shorten", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := ShortenHandler(shortener)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Status = %v, want %v", rr.Code, tt.expectedStatus)
			}

			if tt.expectedError != "" {
				var resp ErrorResponse
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if resp.Error != tt.expectedError {
					t.Errorf("Error = %v, want %v", resp.Error, tt.expectedError)
				}
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	shortener := NewURLShortener()

	// creating test data
	testURL := "https://example.com"
	shortID, _ := shortener.Shorten(testURL)

	tests := []struct {
		name             string
		method           string
		path             string
		expectedStatus   int
		expectedLocation string
	}{
		{
			name:             "succesful redirect",
			method:           http.MethodGet,
			path:             "/" + shortID,
			expectedStatus:   http.StatusFound,
			expectedLocation: testURL,
		},
		{
			name:             "not exist shortner URL",
			method:           http.MethodGet,
			path:             "/notexist",
			expectedStatus:   http.StatusNotFound,
			expectedLocation: "",
		},
		{
			name:             "invalid method",
			method:           http.MethodPost,
			path:             "/" + shortID,
			expectedStatus:   http.StatusMethodNotAllowed,
			expectedLocation: "",
		},
		{
			name:             "empty way",
			method:           http.MethodGet,
			path:             "/",
			expectedStatus:   http.StatusBadRequest,
			expectedLocation: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()
			handler := RedirectHandler(shortener)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Status = %v, want %v", rr.Code, tt.expectedStatus)
			}

			if tt.expectedLocation != "" {
				location := rr.Header().Get("Location")
				if location != tt.expectedLocation {
					t.Errorf("Location = %v, want %v", location, tt.expectedLocation)
				}
			}
		})
	}
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	HealthHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status = %v, want %v", rr.Code, http.StatusOK)
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Errorf("Status = %v, want 'ok'", resp["status"])
	}
}
