// Package hashmiddleware реализует обертку над http.Request и http.ResponseWriter, которая считает и проверяет валидность
// хеша в запросе и выдает в ответе. Хеш передается в заголовке HashSHA256
//
// Package hashmiddleware implements a wrapper around http.Request and http.ResponseWriter that calculates and checks the validity
// hash in the request and issues a response. The hash is passed in the HashSHA256 header
package hashmiddleware

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/h2p2f/practicum-metrics/internal/server/servererrors"
)

// CheckDataHash  - функция для проверки хеша данных запроса
//
// CheckDataHash - function to check hash of request data
func checkDataHash(checkSum string, key string, data []byte) (bool, error) {
	if key == "" {
		return false, servererrors.ErrEmptyKey
	}
	requestCheckSum := sha256.Sum256(data)
	controlCheckSum := fmt.Sprintf("%x", requestCheckSum)
	if checkSum != controlCheckSum {
		fmt.Println("wrong checksum")
		fmt.Println(checkSum)
		fmt.Println(controlCheckSum)
		return false, nil
	}
	return true, nil
}

// GetHash  - функция для получения хеша данных запроса
//
// GetHash - function to get hash of request data
func GetHash(key string, value []byte) ([32]byte, error) {
	if key == "" {
		return [32]byte{}, servererrors.ErrEmptyKey
	}
	checkSum := sha256.Sum256(value)
	return checkSum, nil
}

// HashMiddleware - middleware для проверки хеша данных запроса
// и добавления хеша данных ответа
// key - секретный ключ для хеша, если пустой - хеш не будет проверяться и добавляться
//
// HashMiddleware - middleware to check hash of request data
// and add hash of response data
// key - secret key for hash, if empty - hash will not be checked and added
func HashMiddleware(log *zap.Logger, key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var buf bytes.Buffer
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			checkSum := r.Header.Get("HashSHA256")
			if checkSum != "" && key != "" {
				log.Info("checkSum", zap.String("checkSum", checkSum))
				ok, err := checkDataHash(checkSum, key, buf.Bytes())
				if err != nil {
					http.Error(w, "Bad request", http.StatusBadRequest)
					return
				}
				if !ok {
					http.Error(w, "Bad request", http.StatusBadRequest)
					return
				}
			}
			r.Body = io.NopCloser(&buf)
			capture := &responseCapture{w: w}
			next.ServeHTTP(capture, r)
			hash, err := GetHash(key, capture.body)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			w.Header().Set("HashSHA256", fmt.Sprintf("%x", hash))
		}
		return http.HandlerFunc(fn)
	}
}

// responseCapture - структура для захвата данных ответа
//
// responseCapture - struct to capture response data
type responseCapture struct {
	w    http.ResponseWriter
	body []byte
}

// Header - функция для получения заголовка ответа
//
// Header - function to get response header
func (c *responseCapture) Header() http.Header {
	return c.w.Header()
}

// Write - функция для записи данных ответа
//
// Write - function to write response data
func (c *responseCapture) Write(b []byte) (int, error) {
	c.body = append(c.body, b...)
	return c.w.Write(b)
}

// WriteHeader - функция для записи заголовка ответа
//
// WriteHeader - function to write response header
func (c *responseCapture) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}
