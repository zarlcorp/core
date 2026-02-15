package zcrypto

import (
	"fmt"
	"io"

	"filippo.io/age"
)

// EncryptAge encrypts src to dst using age's scrypt recipient format.
// Output is compatible with `age -d -p`.
func EncryptAge(password string, src io.Reader, dst io.Writer) error {
	r, err := age.NewScryptRecipient(password)
	if err != nil {
		return fmt.Errorf("age encrypt: %w", err)
	}

	w, err := age.Encrypt(dst, r)
	if err != nil {
		return fmt.Errorf("age encrypt: %w", err)
	}

	if _, err := io.Copy(w, src); err != nil {
		return fmt.Errorf("age encrypt: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("age encrypt: %w", err)
	}

	return nil
}

// DecryptAge decrypts an age-encrypted stream using a password.
// Input from `age -e -p` is accepted.
func DecryptAge(password string, src io.Reader, dst io.Writer) error {
	id, err := age.NewScryptIdentity(password)
	if err != nil {
		return fmt.Errorf("age decrypt: %w", err)
	}

	r, err := age.Decrypt(src, id)
	if err != nil {
		return fmt.Errorf("age decrypt: %w", err)
	}

	if _, err := io.Copy(dst, r); err != nil {
		return fmt.Errorf("age decrypt: %w", err)
	}

	return nil
}

// EncryptAgeKey encrypts src to dst for one or more age X25519 public key recipients.
// Each recipient string must be a valid age public key (e.g. "age1...").
// Output is compatible with `age -r <pubkey>`.
func EncryptAgeKey(recipients []string, src io.Reader, dst io.Writer) error {
	parsed := make([]age.Recipient, 0, len(recipients))
	for _, s := range recipients {
		r, err := age.ParseX25519Recipient(s)
		if err != nil {
			return fmt.Errorf("age encrypt: parse recipient: %w", err)
		}
		parsed = append(parsed, r)
	}

	w, err := age.Encrypt(dst, parsed...)
	if err != nil {
		return fmt.Errorf("age encrypt: %w", err)
	}

	if _, err := io.Copy(w, src); err != nil {
		return fmt.Errorf("age encrypt: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("age encrypt: %w", err)
	}

	return nil
}

// DecryptAgeKey decrypts an age-encrypted stream using an X25519 identity (private key).
// The identity string must be a valid age secret key (e.g. "AGE-SECRET-KEY-1...").
// Compatible with files encrypted via `age -r <pubkey>`.
func DecryptAgeKey(identity string, src io.Reader, dst io.Writer) error {
	id, err := age.ParseX25519Identity(identity)
	if err != nil {
		return fmt.Errorf("age decrypt: parse identity: %w", err)
	}

	r, err := age.Decrypt(src, id)
	if err != nil {
		return fmt.Errorf("age decrypt: %w", err)
	}

	if _, err := io.Copy(dst, r); err != nil {
		return fmt.Errorf("age decrypt: %w", err)
	}

	return nil
}
