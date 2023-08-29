// Package hash реализует логику получения хеша данных запроса.
//
// Package hash implements the logic of getting the hash of the request data.
package hash

import (
	"crypto/sha256"
)

// GetHash - function to get hash of request data
func GetHash(value []byte) [32]byte {
	checkSum := sha256.Sum256(value)
	return checkSum
}
