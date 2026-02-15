package zcrypto_test

import (
	"bytes"
	"testing"

	"github.com/zarlcorp/core/pkg/zcrypto"
)

func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, zcrypto.KeySize)
	for i := range key {
		key[i] = byte(i)
	}

	tests := []struct {
		name      string
		key       []byte
		plaintext []byte
		wantErr   bool
	}{
		{
			name:      "round trip",
			key:       key,
			plaintext: []byte("hello, zcrypto"),
		},
		{
			name:      "empty plaintext",
			key:       key,
			plaintext: []byte{},
		},
		{
			name:      "binary data",
			key:       key,
			plaintext: []byte{0x00, 0xff, 0x01, 0xfe},
		},
		{
			name:    "short key",
			key:     key[:16],
			wantErr: true,
		},
		{
			name:    "empty key",
			key:     []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct, err := zcrypto.Encrypt(tt.key, tt.plaintext)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("encrypt: %v", err)
			}

			got, err := zcrypto.Decrypt(tt.key, ct)
			if err != nil {
				t.Fatalf("decrypt: %v", err)
			}

			if !bytes.Equal(got, tt.plaintext) {
				t.Fatalf("got %q, want %q", got, tt.plaintext)
			}
		})
	}
}

func TestEncryptNonceUniqueness(t *testing.T) {
	key := make([]byte, zcrypto.KeySize)
	plaintext := []byte("same input")

	ct1, err := zcrypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt 1: %v", err)
	}

	ct2, err := zcrypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt 2: %v", err)
	}

	if bytes.Equal(ct1, ct2) {
		t.Fatal("two encryptions of same plaintext produced identical ciphertext")
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key1 := make([]byte, zcrypto.KeySize)
	key2 := make([]byte, zcrypto.KeySize)
	key2[0] = 0xff

	ct, err := zcrypto.Encrypt(key1, []byte("secret"))
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	_, err = zcrypto.Decrypt(key2, ct)
	if err == nil {
		t.Fatal("expected error decrypting with wrong key")
	}
}

func TestDecryptTamperedCiphertext(t *testing.T) {
	key := make([]byte, zcrypto.KeySize)

	ct, err := zcrypto.Encrypt(key, []byte("secret"))
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	// flip a byte in the ciphertext portion (after the nonce)
	ct[len(ct)-1] ^= 0xff

	_, err = zcrypto.Decrypt(key, ct)
	if err == nil {
		t.Fatal("expected error decrypting tampered ciphertext")
	}
}

func TestDecryptCiphertextTooShort(t *testing.T) {
	key := make([]byte, zcrypto.KeySize)

	_, err := zcrypto.Decrypt(key, []byte("short"))
	if err == nil {
		t.Fatal("expected error for short ciphertext")
	}
}

func TestDeriveKey(t *testing.T) {
	tests := []struct {
		name string
		salt []byte
	}{
		{name: "with salt", salt: bytes.Repeat([]byte{0xab}, zcrypto.SaltSize)},
		{name: "nil salt generates random", salt: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, usedSalt, err := zcrypto.DeriveKey([]byte("password"), tt.salt)
			if err != nil {
				t.Fatalf("derive key: %v", err)
			}

			if len(key) != zcrypto.KeySize {
				t.Fatalf("key length = %d, want %d", len(key), zcrypto.KeySize)
			}

			if len(usedSalt) != zcrypto.SaltSize {
				t.Fatalf("salt length = %d, want %d", len(usedSalt), zcrypto.SaltSize)
			}

			if tt.salt != nil && !bytes.Equal(usedSalt, tt.salt) {
				t.Fatal("provided salt was not used as-is")
			}
		})
	}
}

func TestDeriveKeyDeterministic(t *testing.T) {
	password := []byte("password")
	salt := bytes.Repeat([]byte{0x01}, zcrypto.SaltSize)

	k1, _, err := zcrypto.DeriveKey(password, salt)
	if err != nil {
		t.Fatalf("derive 1: %v", err)
	}

	k2, _, err := zcrypto.DeriveKey(password, salt)
	if err != nil {
		t.Fatalf("derive 2: %v", err)
	}

	if !bytes.Equal(k1, k2) {
		t.Fatal("same password + salt produced different keys")
	}
}

func TestDeriveKeyDifferentPasswords(t *testing.T) {
	salt := bytes.Repeat([]byte{0x01}, zcrypto.SaltSize)

	k1, _, err := zcrypto.DeriveKey([]byte("password1"), salt)
	if err != nil {
		t.Fatalf("derive 1: %v", err)
	}

	k2, _, err := zcrypto.DeriveKey([]byte("password2"), salt)
	if err != nil {
		t.Fatalf("derive 2: %v", err)
	}

	if bytes.Equal(k1, k2) {
		t.Fatal("different passwords produced same key")
	}
}

func TestExpandKey(t *testing.T) {
	secret := bytes.Repeat([]byte{0xaa}, 32)
	salt := bytes.Repeat([]byte{0xbb}, 16)

	tests := []struct {
		name string
		info []byte
	}{
		{name: "file-encryption", info: []byte("file-encryption")},
		{name: "auth-token", info: []byte("auth-token")},
	}

	keys := make(map[string][]byte)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := zcrypto.ExpandKey(secret, salt, tt.info)
			if err != nil {
				t.Fatalf("expand key: %v", err)
			}

			if len(k) != zcrypto.KeySize {
				t.Fatalf("key length = %d, want %d", len(k), zcrypto.KeySize)
			}

			keys[tt.name] = k
		})
	}

	// domain separation: different info values must produce different keys
	if bytes.Equal(keys["file-encryption"], keys["auth-token"]) {
		t.Fatal("different info values produced same key")
	}
}

func TestExpandKeyDeterministic(t *testing.T) {
	secret := bytes.Repeat([]byte{0xaa}, 32)
	salt := bytes.Repeat([]byte{0xbb}, 16)
	info := []byte("test")

	k1, err := zcrypto.ExpandKey(secret, salt, info)
	if err != nil {
		t.Fatalf("expand 1: %v", err)
	}

	k2, err := zcrypto.ExpandKey(secret, salt, info)
	if err != nil {
		t.Fatalf("expand 2: %v", err)
	}

	if !bytes.Equal(k1, k2) {
		t.Fatal("same inputs produced different keys")
	}
}

func TestRandBytes(t *testing.T) {
	sizes := []int{0, 1, 16, 32, 64}

	for _, n := range sizes {
		b, err := zcrypto.RandBytes(n)
		if err != nil {
			t.Fatalf("RandBytes(%d): %v", n, err)
		}
		if len(b) != n {
			t.Fatalf("RandBytes(%d) returned %d bytes", n, len(b))
		}
	}
}

func TestRandBytesUnique(t *testing.T) {
	b1, err := zcrypto.RandBytes(32)
	if err != nil {
		t.Fatalf("rand 1: %v", err)
	}

	b2, err := zcrypto.RandBytes(32)
	if err != nil {
		t.Fatalf("rand 2: %v", err)
	}

	if bytes.Equal(b1, b2) {
		t.Fatal("two random 32-byte outputs are identical")
	}
}

func TestRandHex(t *testing.T) {
	sizes := []int{1, 8, 16, 32}

	for _, n := range sizes {
		s, err := zcrypto.RandHex(n)
		if err != nil {
			t.Fatalf("RandHex(%d): %v", n, err)
		}
		if len(s) != 2*n {
			t.Fatalf("RandHex(%d) returned %d chars, want %d", n, len(s), 2*n)
		}
	}
}

func TestRandHexUnique(t *testing.T) {
	s1, err := zcrypto.RandHex(16)
	if err != nil {
		t.Fatalf("hex 1: %v", err)
	}

	s2, err := zcrypto.RandHex(16)
	if err != nil {
		t.Fatalf("hex 2: %v", err)
	}

	if s1 == s2 {
		t.Fatal("two random hex outputs are identical")
	}
}

func TestErase(t *testing.T) {
	b := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	zcrypto.Erase(b)

	for i, v := range b {
		if v != 0 {
			t.Fatalf("byte %d = %d, want 0", i, v)
		}
	}
}

func TestEraseEmpty(t *testing.T) {
	// should not panic on empty slice
	zcrypto.Erase([]byte{})
	zcrypto.Erase(nil)
}

func TestEncryptFileDecryptFile(t *testing.T) {
	key := make([]byte, zcrypto.KeySize)
	for i := range key {
		key[i] = byte(i)
	}

	original := []byte("file contents to encrypt")

	var encrypted bytes.Buffer
	if err := zcrypto.EncryptFile(key, bytes.NewReader(original), &encrypted); err != nil {
		t.Fatalf("encrypt file: %v", err)
	}

	var decrypted bytes.Buffer
	if err := zcrypto.DecryptFile(key, bytes.NewReader(encrypted.Bytes()), &decrypted); err != nil {
		t.Fatalf("decrypt file: %v", err)
	}

	if !bytes.Equal(decrypted.Bytes(), original) {
		t.Fatalf("got %q, want %q", decrypted.Bytes(), original)
	}
}

func TestEncryptFileWrongKey(t *testing.T) {
	key1 := make([]byte, zcrypto.KeySize)
	key2 := make([]byte, zcrypto.KeySize)
	key2[0] = 0xff

	var encrypted bytes.Buffer
	if err := zcrypto.EncryptFile(key1, bytes.NewReader([]byte("secret")), &encrypted); err != nil {
		t.Fatalf("encrypt file: %v", err)
	}

	var decrypted bytes.Buffer
	if err := zcrypto.DecryptFile(key2, bytes.NewReader(encrypted.Bytes()), &decrypted); err == nil {
		t.Fatal("expected error decrypting with wrong key")
	}
}

func TestEncryptFileEmptyContent(t *testing.T) {
	key := make([]byte, zcrypto.KeySize)

	var encrypted bytes.Buffer
	if err := zcrypto.EncryptFile(key, bytes.NewReader([]byte{}), &encrypted); err != nil {
		t.Fatalf("encrypt file: %v", err)
	}

	var decrypted bytes.Buffer
	if err := zcrypto.DecryptFile(key, bytes.NewReader(encrypted.Bytes()), &decrypted); err != nil {
		t.Fatalf("decrypt file: %v", err)
	}

	if len(decrypted.Bytes()) != 0 {
		t.Fatalf("expected empty output, got %d bytes", len(decrypted.Bytes()))
	}
}
