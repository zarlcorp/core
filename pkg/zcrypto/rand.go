package zcrypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// RandBytes returns n cryptographically random bytes.
func RandBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("read random bytes: %w", err)
	}
	return b, nil
}

// RandHex returns a hex-encoded string of n random bytes (2n chars).
func RandHex(n int) (string, error) {
	b, err := RandBytes(n)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
