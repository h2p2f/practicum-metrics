package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
)

type CompressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

type CompressReader struct {
	r io.ReadCloser
	z *gzip.Reader
}

func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	zw := gzip.NewWriter(w)
	return &CompressWriter{w, zw}
}

func (cw *CompressWriter) Header() http.Header {
	return cw.w.Header()
}

func (cw *CompressWriter) Write(b []byte) (int, error) {
	return cw.zw.Write(b)
}

func (cw *CompressWriter) WriteHeader(statusCode int) {
	cw.w.WriteHeader(statusCode)
	if statusCode >= 200 && statusCode < 300 {
		cw.w.Header().Set("Content-Encoding", "gzip")
	}
}
func (cw *CompressWriter) Close() error {
	return cw.zw.Close()
}
func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	z, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &CompressReader{r, z}, nil
}

func (cr *CompressReader) Read(b []byte) (int, error) {
	return cr.z.Read(b)
}

func (cr *CompressReader) Close() error {
	if err := cr.z.Close(); err != nil {
		return err
	}
	return cr.r.Close()
}
