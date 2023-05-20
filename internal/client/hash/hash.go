package hash

import (
	"crypto/sha256"
	"errors"
)

var ErrEmptyKey = errors.New("empty key")

func GetHash(key string, value []byte) ([32]byte, error) {
	if key == "" {
		return [32]byte{}, ErrEmptyKey
	}
	checkSum := sha256.Sum256(value)
	return checkSum, nil
}
