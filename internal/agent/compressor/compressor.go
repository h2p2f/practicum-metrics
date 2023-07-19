// Package compressor реализует логику сжатия данных в gzip.
//
// Package compressor implements gzip data compression logic.
package compressor

import (
	"bytes"
	"compress/gzip"
)

// Compress - функция для сжатия данных в gzip.
//
// Compress - function for compressing data into gzip.
func Compress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(buf)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
