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
//   - age-compatible password and key-based encryption
//
// # AES-256-GCM
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
//
// # Age Password Encryption
//
// Password-based age encryption produces output compatible with the age CLI.
//
//	err := zcrypto.EncryptAge("passphrase", src, dst)
//	if err != nil {
//	    // handle error
//	}
//
//	err = zcrypto.DecryptAge("passphrase", src, dst)
//	if err != nil {
//	    // handle error
//	}
//
// # Age Key Encryption
//
// Key-based age encryption uses X25519 public keys and is compatible with
// age -r and age -i.
//
//	err := zcrypto.EncryptAgeKey([]string{"age1..."}, src, dst)
//	if err != nil {
//	    // handle error
//	}
//
//	err = zcrypto.DecryptAgeKey("AGE-SECRET-KEY-1...", src, dst)
//	if err != nil {
//	    // handle error
//	}
package zcrypto
