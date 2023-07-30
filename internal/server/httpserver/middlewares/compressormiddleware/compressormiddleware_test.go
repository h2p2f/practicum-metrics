package compressormiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestZipMiddleware(t *testing.T) {
	type want struct {
		contentType     string
		headerContent   string
		jsonString      string
		acceptEncoding  string
		contentEncoding string
		statusCode      int
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Positive test 1",
			want: want{
				statusCode:      200,
				contentType:     "application/x-gzip",
				headerContent:   "gzip",
				jsonString:      "{\"id\": \"testID\", \"type\": \"counter\",  \"delta\": 1}",
				acceptEncoding:  "gzip",
				contentEncoding: "",
			},
		},
		{
			name: "Positive test 2",
			want: want{
				statusCode:      200,
				contentType:     "application/json",
				headerContent:   "",
				jsonString:      "{\"id\": \"testID\", \"type\": \"gauge\",  \"value\": 1.0000000000001}",
				acceptEncoding:  "",
				contentEncoding: "",
			},
		},
		{
			name: "Positive test 3",
			want: want{
				statusCode:      200,
				contentType:     "application/json",
				headerContent:   "",
				jsonString:      "some text",
				acceptEncoding:  "",
				contentEncoding: "gzip",
			},
		},
		{
			name: "Negative test 1",
			want: want{
				statusCode:      200,
				contentType:     "application/json",
				headerContent:   "",
				jsonString:      "another text",
				acceptEncoding:  "deflate",
				contentEncoding: "deflate",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/update/", nil)
			req.Header.Set("Accept-Encoding", tt.want.acceptEncoding)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Content-Encoding", "testID")
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(tt.want.jsonString)) //nolint:errcheck
			})
			ZipMiddleware(handler).ServeHTTP(rr, req)

			if rr.Header().Get("Content-Encoding") != tt.want.headerContent {
				t.Errorf("Expected response %s, but got Content-Encoding header '%s'", tt.want.headerContent, rr.Header().Get("Content-Encoding"))
			}
			if rr.Header().Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected response %s, but got Content-Type header '%s'", tt.want.contentType, rr.Header().Get("Content-Type"))
			}
			if rr.Code != tt.want.statusCode {
				t.Errorf("Expected response code %d, but got %d", tt.want.statusCode, rr.Code)
			}
		})
	}
}
