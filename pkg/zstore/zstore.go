// Package zstore provides generic encrypted key-value storage built on
// zfilesystem and zcrypto. It uses HKDF to derive per-collection sub-keys
// from a single master key, so compromising one collection's key does not
// expose others.
//
// # Usage
//
//	fs := zfilesystem.NewMemFS()
//	s, err := zstore.Open(fs, []byte("password"))
//	if err != nil {
//	    // handle error
//	}
//	defer s.Close()
//
//	col := zstore.Collection[MyType](s, "things")
//	err = col.Put("id1", MyType{Name: "hello"})
package zstore

import (
	"errors"
	"fmt"

	"github.com/zarlcorp/core/pkg/zcrypto"
	"github.com/zarlcorp/core/pkg/zfilesystem"
)

const (
	saltFile   = "salt"
	verifyFile = "verify"
	verifyText = "zstore-verify-ok"
)

// ErrNotFound is returned when a key does not exist.
var ErrNotFound = errors.New("not found")

// ErrWrongPassword is returned when the password does not match the store.
var ErrWrongPassword = errors.New("wrong password")

// Store is an encrypted key-value store. It holds a master key derived from
// the user's password and supports multiple typed collections, each with
// its own HKDF-derived sub-key.
type Store struct {
	fs        zfilesystem.ReadWriteFileFS
	masterKey []byte
	salt      []byte
	subKeys   [][]byte
}

// Open creates or opens a store. On first run it generates a salt, derives a
// master key, creates a verification token, and stores both via the filesystem.
// On subsequent runs it reads the salt, derives the key, and verifies the
// password by decrypting the verification token.
func Open(fs zfilesystem.ReadWriteFileFS, password []byte, opts ...Option) (*Store, error) {
	_ = applyOptions(opts)

	salt, err := fs.ReadFile(saltFile)
	if err != nil {
		return initStore(fs, password)
	}

	return openStore(fs, password, salt)
}

// initStore handles first-run initialization: generate salt, derive key,
// encrypt a verification token, persist both.
func initStore(fs zfilesystem.ReadWriteFileFS, password []byte) (*Store, error) {
	salt, err := zcrypto.RandBytes(zcrypto.SaltSize)
	if err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}

	if err := fs.WriteFile(saltFile, salt, 0o600); err != nil {
		return nil, fmt.Errorf("write salt: %w", err)
	}

	key, _, err := zcrypto.DeriveKey(password, salt)
	if err != nil {
		return nil, fmt.Errorf("derive key: %w", err)
	}

	token, err := zcrypto.Encrypt(key, []byte(verifyText))
	if err != nil {
		return nil, fmt.Errorf("encrypt verification token: %w", err)
	}

	if err := fs.WriteFile(verifyFile, token, 0o600); err != nil {
		return nil, fmt.Errorf("write verification token: %w", err)
	}

	return &Store{fs: fs, masterKey: key, salt: salt}, nil
}

// openStore handles subsequent opens: derive key from existing salt, verify
// password against the stored verification token.
func openStore(fs zfilesystem.ReadWriteFileFS, password, salt []byte) (*Store, error) {
	key, _, err := zcrypto.DeriveKey(password, salt)
	if err != nil {
		return nil, fmt.Errorf("derive key: %w", err)
	}

	token, err := fs.ReadFile(verifyFile)
	if err != nil {
		return nil, fmt.Errorf("read verification token: %w", err)
	}

	plain, err := zcrypto.Decrypt(key, token)
	if err != nil || string(plain) != verifyText {
		zcrypto.Erase(key)
		return nil, ErrWrongPassword
	}

	return &Store{fs: fs, masterKey: key, salt: salt}, nil
}

// Close erases the master key and all sub-keys from memory.
func (s *Store) Close() error {
	zcrypto.Erase(s.masterKey)
	for _, k := range s.subKeys {
		zcrypto.Erase(k)
	}
	return nil
}
