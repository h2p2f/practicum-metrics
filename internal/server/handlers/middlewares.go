package handlers

import (
	"net/http"
	"strings"
)

func GzipHanle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originWriter := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")
		if supportGzip {
			compressedWriter := NewCompressWriter(w)
			originWriter = compressedWriter
			defer compressedWriter.Close()
		}
		contentEncoding := r.Header.Get("Content-Encoding")
		sendedGzip := strings.Contains(contentEncoding, "gzip")
		if sendedGzip {
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
