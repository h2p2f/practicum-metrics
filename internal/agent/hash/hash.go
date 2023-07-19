// Package hash реализует логику получения хеша данных запроса.
//
// Package hash implements the logic of getting the hash of the request data.
package hash

import (
	"crypto/sha256"
	"errors"
)

// ErrEmptyKey - ошибка пустого ключа
//
// ErrEmptyKey - empty key error
var ErrEmptyKey = errors.New("empty key")

// GetHash - функция для получения хеша данных запроса
//
// GetHash - function to get hash of request data
func GetHash(key string, value []byte) [32]byte {
	if key == "" {
		return [32]byte{}
	}
	checkSum := sha256.Sum256(value)

	return checkSum
}
