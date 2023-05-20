package httpserver

import (
	"crypto/sha256"
	"errors"
	"fmt"
)

var ErrEmptyKey = errors.New("empty key")

func checkDataHash(checkSum string, key string, data []byte) (bool, error) {
	if key == "" {
		return false, ErrEmptyKey
	}
	h := sha256.New()
	h.Write(data)
	//controlCheckSum := h.Sum(nil)
	requestCheckSum := sha256.Sum256(data)
	controlCheckSum := fmt.Sprintf("%x", requestCheckSum)
	if checkSum != controlCheckSum {
		return false, nil
	}
	return true, nil
}

func GetHash(key string, value []byte) ([32]byte, error) {
	if key == "" {
		return [32]byte{}, ErrEmptyKey
	}
	checkSum := sha256.Sum256(value)
	return checkSum, nil
}
