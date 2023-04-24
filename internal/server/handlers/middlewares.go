package handlers

import (
	"net/http"
	"strings"
)

func GzipHanle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originWriter := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		contentEncoding := r.Header.Get("Content-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")
		sendGzip := strings.Contains(contentEncoding, "gzip")
		if supportGzip && sendGzip {
			compressedWriter := NewCompressWriter(w)
			originWriter = compressedWriter
			defer compressedWriter.Close()
			compressedReader, err := NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			r.Body = compressedReader
			defer compressedReader.Close()
		}

		h.ServeHTTP(originWriter, r)
	})

}
