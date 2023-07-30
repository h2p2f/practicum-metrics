package hashmiddleware

import (
	"bytes"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHashMiddleware(t *testing.T) {

	logger := zaptest.NewLogger(t)

	tests := []struct {
		name     string
		key      string
		checkSum string
		body     []byte
		expected int
	}{
		{
			name:     "Valid request 1",
			key:      "secret",
			body:     []byte("example"),
			checkSum: "50d858e0985ecc7f60418aaf0cc5ab587f42c2570a884095a9e8ccacd0f6545c",
			expected: http.StatusOK,
		},
		{
			name:     "Invalid checkSum",
			key:      "secret",
			body:     []byte("example"),
			checkSum: "invalid",
			expected: http.StatusBadRequest,
		},
		{
			name:     "Valid request 2",
			key:      "1",
			body:     []byte("example"),
			checkSum: "50d858e0985ecc7f60418aaf0cc5ab587f42c2570a884095a9e8ccacd0f6545c",
			expected: http.StatusOK,
		},
		{
			name:     "Empty Body",
			key:      "secret",
			body:     nil,
			checkSum: "",
			expected: http.StatusOK,
		},
	}

	// Iterate over test cases.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			if tt.checkSum != "" {
				req.Header.Set("HashSHA256", tt.checkSum)
			}
			rr := httptest.NewRecorder()

			HashMiddleware(logger, tt.key)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})).ServeHTTP(rr, req)

			if rr.Code != tt.expected {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.expected)
			}

			if tt.expected == http.StatusOK && rr.Header().Get("HashSHA256") == "" {
				t.Error("handler did not set HashSHA256 header")
			}
		})
	}

}
