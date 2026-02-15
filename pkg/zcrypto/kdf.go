package zcrypto

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/hkdf"
)

const (
	// Argon2id parameters.
	argon2Time    = 1
	argon2Memory  = 64 * 1024 // 64 MB
	argon2Threads = 4

	// SaltSize is the default salt length for Argon2id.
	SaltSize = 16
)

// DeriveKey derives a 32-byte key from a password and salt using Argon2id.
// If salt is nil, a random 16-byte salt is generated.
// Returns the derived key and the salt used.
func DeriveKey(password, salt []byte) (key, usedSalt []byte, err error) {
	if salt == nil {
		salt = make([]byte, SaltSize)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, fmt.Errorf("generate salt: %w", err)
		}
	}

	k := argon2.IDKey(password, salt, argon2Time, argon2Memory, argon2Threads, KeySize)
	return k, salt, nil
}

// ExpandKey derives a new key from existing key material using HKDF-SHA256.
// info provides context for domain separation (e.g. "file-encryption", "auth-token").
func ExpandKey(secret, salt, info []byte) ([]byte, error) {
	r := hkdf.New(sha256.New, secret, salt, info)

	key := make([]byte, KeySize)
	if _, err := io.ReadFull(r, key); err != nil {
		return nil, fmt.Errorf("expand key: %w", err)
	}

	return key, nil
}
