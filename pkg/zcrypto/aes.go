// Package zcrypto provides encryption primitives for zarlcorp privacy tools.
//
// It composes proven Go stdlib and x/crypto primitives — no custom cryptography.
// All error paths return errors; the package never panics.
//
// # Features
//
//   - AES-256-GCM symmetric encryption
//   - Argon2id password-based key derivation
//   - HKDF-SHA256 key expansion
//   - Cryptographic random generation
//   - Secure memory erasure
//   - File encryption/decryption helpers
//
// # Usage
//
//	key, salt, err := zcrypto.DeriveKey([]byte("passphrase"), nil)
//	if err != nil {
//	    // handle error
//	}
//
//	ciphertext, err := zcrypto.Encrypt(key, []byte("secret"))
//	if err != nil {
//	    // handle error
//	}
//
//	plaintext, err := zcrypto.Decrypt(key, ciphertext)
//	if err != nil {
//	    // handle error
//	}
package zcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
)

const (
	// KeySize is the required key length for AES-256.
	KeySize = 32

	// NonceSize is the standard nonce length for AES-GCM.
	NonceSize = 12
)

// Encrypt encrypts plaintext using AES-256-GCM with the given key.
// Key must be exactly 32 bytes. Returns ciphertext with nonce prepended.
func Encrypt(key, plaintext []byte) ([]byte, error) {
	if len(key) != KeySize {
		return nil, fmt.Errorf("key must be %d bytes, got %d", KeySize, len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext produced by Encrypt.
// Key must be exactly 32 bytes. Expects nonce prepended to ciphertext.
func Decrypt(key, ciphertext []byte) ([]byte, error) {
	if len(key) != KeySize {
		return nil, fmt.Errorf("key must be %d bytes, got %d", KeySize, len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ct := ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	return plaintext, nil
}
