package hash

import (
	"crypto/sha256"
	"errors"
	"fmt"
)

var ErrEmptyKey = errors.New("empty key")

func GetHash(key string, value []byte) ([32]byte, error) {
	if key == "" {
		return [32]byte{}, ErrEmptyKey
	}
	checkSum := sha256.Sum256(value)
	fmt.Sprintf("checksum %s", checkSum)
	return checkSum, nil
}
