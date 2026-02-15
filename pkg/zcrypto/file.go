package zcrypto

import (
	"fmt"
	"io"
)

// EncryptFile encrypts the contents of src and writes the result to dst.
// Uses AES-256-GCM with the given 32-byte key.
func EncryptFile(key []byte, src io.Reader, dst io.Writer) error {
	plaintext, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("read source: %w", err)
	}

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}

	if _, err := dst.Write(ciphertext); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

// DecryptFile decrypts a file encrypted by EncryptFile.
func DecryptFile(key []byte, src io.Reader, dst io.Writer) error {
	ciphertext, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("read source: %w", err)
	}

	plaintext, err := Decrypt(key, ciphertext)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}

	if _, err := dst.Write(plaintext); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}
