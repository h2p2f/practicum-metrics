package decryptormiddleware

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecryptMiddleware(t *testing.T) {
	// Generate RSA keys for testing
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Errorf("Error generating RSA key: %v", err)
	}
	anotherPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Errorf("Error generating RSA key: %v", err)
	}

	tests := []struct {
		name     string
		key      *rsa.PrivateKey
		path     string
		body     []byte
		expected int
	}{
		{
			name:     "Valid request 1",
			key:      privateKey,
			body:     []byte("example"),
			path:     "/update/",
			expected: http.StatusOK,
		},
		{
			name:     "Valid request 2",
			key:      privateKey,
			body:     []byte("example"),
			path:     "/updates/",
			expected: http.StatusOK,
		},
		{
			name:     "Invalid request 1",
			key:      anotherPrivateKey,
			body:     []byte("example"),
			path:     "/update/",
			expected: http.StatusUnprocessableEntity,
		},
		{
			name:     "Empty Body, wrong key and request to main page",
			key:      anotherPrivateKey,
			body:     []byte("example"),
			path:     "/",
			expected: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			// Create test request
			if tt.path == "/update/" || tt.path == "/updates/" {
				body, err = rsa.EncryptPKCS1v15(rand.Reader, &privateKey.PublicKey, tt.body)
				//encryptedBody, err := rsa.EncryptPKCS1v15(rand.Reader, &privateKey.PublicKey, []byte(reqBody))
				if err != nil {
					t.Errorf("Error encrypting test request body: %v", err)
				}
			} else {
				body = tt.body
			}
			req, err := http.NewRequest("POST", tt.path, bytes.NewReader(body))
			if err != nil {
				t.Errorf("Error creating test request: %v", err)
			}

			// Create test response recorder
			rr := httptest.NewRecorder()

			// Create middleware handler
			middlewareHandler := DecryptMiddleware(tt.key)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check if request body was decrypted
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Errorf("Error reading decrypted request body: %v", err)
				}
				if string(body) != string(tt.body) {
					t.Errorf("Request body was not decrypted properly")
				}

				// Write response
				w.WriteHeader(http.StatusOK)
			}))

			// Serve test request
			middlewareHandler.ServeHTTP(rr, req)

			// Check if response status code is OK
			if rr.Code != tt.expected {
				t.Errorf("Response status code not as expected: %v", rr.Code)
			}
		})
	}
}
