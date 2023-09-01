// Package decryptormiddleware implements http.Handler wrapper, which decrypts request body if
// RSA key is present in config
// Hard limitations:
// decrypts only request body, if RSA key is present, doesn't decrypt headers
// works only with /update/ and /updates/ endpoints

package decryptormiddleware

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"
)

// DecryptMiddleware - http.Handler wrapper, which decrypts request body if
// RSA key is present in config
func DecryptMiddleware(RSAKey *rsa.PrivateKey) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if RSAKey != nil && (r.URL.Path == "/update/" || r.URL.Path == "/updates/") {
				var buf bytes.Buffer
				_, err := buf.ReadFrom(r.Body)
				if err != nil {
					http.Error(w, "Bad request", http.StatusBadRequest)
					return
				}
				data, err := rsa.DecryptPKCS1v15(rand.Reader, RSAKey, buf.Bytes())
				if err != nil {
					http.Error(w, "Unprocessable entity", http.StatusUnprocessableEntity)
					return
				}
				r.Body = io.NopCloser(bytes.NewReader(data))
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
