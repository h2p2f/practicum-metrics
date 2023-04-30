package handlers

import (
	"net/http"
	"strings"
)

// GzipHanle is middleware for gzip
func GzipHanle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		contentEncoding := r.Header.Get("Content-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")
		sendGzip := strings.Contains(contentEncoding, "gzip")

		if !supportGzip {
			next.ServeHTTP(w, r)
			return
		}
		if supportGzip && !sendGzip {
			originWriter := w
			compressedWriter := NewCompressWriter(w)
			originWriter = compressedWriter
			defer compressedWriter.Close()
			next.ServeHTTP(originWriter, r)
		}
		if sendGzip {
			originWriter := w
			compressedWriter := NewCompressWriter(w)
			originWriter = compressedWriter
			defer compressedWriter.Close()
			compressedReader, err := NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			r.Body = compressedReader
			defer compressedReader.Close()

			next.ServeHTTP(originWriter, r)

		}
	})

}
