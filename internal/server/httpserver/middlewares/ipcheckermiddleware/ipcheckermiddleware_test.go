package ipcheckermiddleware

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap/zaptest"
)

func TestIpCheckMiddleware(t *testing.T) {
	logger := zaptest.NewLogger(t)
	tests := []struct {
		name     string
		ip       string
		subnet   string
		expected int
	}{
		{
			name:     "Valid IP",
			ip:       "10.1.23.2",
			subnet:   "10.1.23.0/16",
			expected: http.StatusOK,
		},
		{
			name:     "Invalid IP",
			ip:       "10.2.23.2",
			subnet:   "10.1.23.0/16",
			expected: http.StatusForbidden,
		},
		{
			name:     "Empty IP",
			ip:       "",
			subnet:   "10.1.23.0/16",
			expected: http.StatusForbidden,
		},
		{
			name:     "Empty Subnet",
			ip:       "10.1.23.2",
			subnet:   "",
			expected: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subnet := &net.IPNet{}
			err := error(nil)
			if tt.subnet != "" {
				_, subnet, err = net.ParseCIDR(tt.subnet)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				subnet = nil
			}
			req, err := http.NewRequest(http.MethodPost, "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.ip != "" {
				req.Header.Set("X-Real-IP", tt.ip)
			}
			rr := httptest.NewRecorder()
			handler := IPCheckMiddleware(logger, subnet)
			handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})).ServeHTTP(rr, req)
			if status := rr.Code; status != tt.expected {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expected)
			}
		})
	}
}
