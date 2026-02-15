package zcrypto_test

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"filippo.io/age"

	"github.com/zarlcorp/core/pkg/zcrypto"
)

func TestEncryptAgeDecryptAge(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		plaintext string
	}{
		{name: "basic round trip", password: "hunter2", plaintext: "hello, age"},
		{name: "empty plaintext", password: "pass", plaintext: ""},
		{name: "long password", password: strings.Repeat("a", 1024), plaintext: "data"},
		{name: "binary-like data", password: "pw", plaintext: "\x00\xff\x01\xfe"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var encrypted bytes.Buffer
			if err := zcrypto.EncryptAge(tt.password, strings.NewReader(tt.plaintext), &encrypted); err != nil {
				t.Fatalf("encrypt: %v", err)
			}

			var decrypted bytes.Buffer
			if err := zcrypto.DecryptAge(tt.password, &encrypted, &decrypted); err != nil {
				t.Fatalf("decrypt: %v", err)
			}

			if decrypted.String() != tt.plaintext {
				t.Fatalf("got %q, want %q", decrypted.String(), tt.plaintext)
			}
		})
	}
}

func TestDecryptAgeWrongPassword(t *testing.T) {
	var encrypted bytes.Buffer
	if err := zcrypto.EncryptAge("correct", strings.NewReader("secret"), &encrypted); err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	var decrypted bytes.Buffer
	err := zcrypto.DecryptAge("wrong", bytes.NewReader(encrypted.Bytes()), &decrypted)
	if err == nil {
		t.Fatal("expected error decrypting with wrong password")
	}
}

func TestDecryptAgeInvalidInput(t *testing.T) {
	var decrypted bytes.Buffer
	err := zcrypto.DecryptAge("pw", strings.NewReader("not an age file"), &decrypted)
	if err == nil {
		t.Fatal("expected error decrypting invalid input")
	}
}

func TestEncryptAgeKeyDecryptAgeKey(t *testing.T) {
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}

	pubkey := id.Recipient().String()
	privkey := id.String()

	tests := []struct {
		name      string
		plaintext string
	}{
		{name: "basic round trip", plaintext: "hello, age keys"},
		{name: "empty plaintext", plaintext: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var encrypted bytes.Buffer
			if err := zcrypto.EncryptAgeKey([]string{pubkey}, strings.NewReader(tt.plaintext), &encrypted); err != nil {
				t.Fatalf("encrypt: %v", err)
			}

			var decrypted bytes.Buffer
			if err := zcrypto.DecryptAgeKey(privkey, &encrypted, &decrypted); err != nil {
				t.Fatalf("decrypt: %v", err)
			}

			if decrypted.String() != tt.plaintext {
				t.Fatalf("got %q, want %q", decrypted.String(), tt.plaintext)
			}
		})
	}
}

func TestEncryptAgeKeyMultipleRecipients(t *testing.T) {
	id1, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity 1: %v", err)
	}

	id2, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity 2: %v", err)
	}

	recipients := []string{
		id1.Recipient().String(),
		id2.Recipient().String(),
	}

	plaintext := "shared secret"

	var encrypted bytes.Buffer
	if err := zcrypto.EncryptAgeKey(recipients, strings.NewReader(plaintext), &encrypted); err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	// both recipients should be able to decrypt
	for i, id := range []*age.X25519Identity{id1, id2} {
		var decrypted bytes.Buffer
		if err := zcrypto.DecryptAgeKey(id.String(), bytes.NewReader(encrypted.Bytes()), &decrypted); err != nil {
			t.Fatalf("decrypt with identity %d: %v", i+1, err)
		}
		if decrypted.String() != plaintext {
			t.Fatalf("identity %d: got %q, want %q", i+1, decrypted.String(), plaintext)
		}
	}
}

func TestDecryptAgeKeyWrongIdentity(t *testing.T) {
	id1, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity 1: %v", err)
	}

	id2, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity 2: %v", err)
	}

	var encrypted bytes.Buffer
	if err := zcrypto.EncryptAgeKey([]string{id1.Recipient().String()}, strings.NewReader("secret"), &encrypted); err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	var decrypted bytes.Buffer
	err = zcrypto.DecryptAgeKey(id2.String(), bytes.NewReader(encrypted.Bytes()), &decrypted)
	if err == nil {
		t.Fatal("expected error decrypting with wrong identity")
	}
}

func TestEncryptAgeKeyInvalidRecipient(t *testing.T) {
	err := zcrypto.EncryptAgeKey([]string{"not-a-key"}, strings.NewReader("data"), &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for invalid recipient")
	}
}

func TestDecryptAgeKeyInvalidIdentity(t *testing.T) {
	err := zcrypto.DecryptAgeKey("not-a-key", strings.NewReader("data"), &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for invalid identity")
	}
}

// interop tests require the age CLI — skip when unavailable

func ageCLI(t *testing.T) string {
	t.Helper()
	path, err := exec.LookPath("age")
	if err != nil {
		t.Skip("age CLI not found, skipping interop test")
	}
	return path
}

func TestInteropEncryptZcryptoDecryptAgeCLI(t *testing.T) {
	ageBin := ageCLI(t)

	password := "interop-test-password"
	plaintext := "encrypted by zcrypto, decrypted by age CLI"

	var encrypted bytes.Buffer
	if err := zcrypto.EncryptAge(password, strings.NewReader(plaintext), &encrypted); err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	cmd := exec.Command(ageBin, "-d")
	cmd.Stdin = &encrypted
	cmd.Env = append(cmd.Environ(), "AGE_PASSPHRASE="+password)

	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			t.Fatalf("age -d: %v\nstderr: %s", err, ee.Stderr)
		}
		t.Fatalf("age -d: %v", err)
	}

	if string(out) != plaintext {
		t.Fatalf("age CLI decrypted %q, want %q", string(out), plaintext)
	}
}

func TestInteropEncryptAgeCLIDecryptZcrypto(t *testing.T) {
	ageBin := ageCLI(t)

	password := "interop-test-password"
	plaintext := "encrypted by age CLI, decrypted by zcrypto"

	cmd := exec.Command(ageBin, "-e", "-p")
	cmd.Stdin = strings.NewReader(plaintext)
	cmd.Env = append(cmd.Environ(), "AGE_PASSPHRASE="+password)

	encrypted, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			t.Fatalf("age -e -p: %v\nstderr: %s", err, ee.Stderr)
		}
		t.Fatalf("age -e -p: %v", err)
	}

	var decrypted bytes.Buffer
	if err := zcrypto.DecryptAge(password, bytes.NewReader(encrypted), &decrypted); err != nil {
		t.Fatalf("decrypt: %v", err)
	}

	if decrypted.String() != plaintext {
		t.Fatalf("got %q, want %q", decrypted.String(), plaintext)
	}
}

func TestInteropKeyEncryptZcryptoDecryptAgeCLI(t *testing.T) {
	ageBin := ageCLI(t)

	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}

	plaintext := "key-encrypted by zcrypto"

	var encrypted bytes.Buffer
	if err := zcrypto.EncryptAgeKey([]string{id.Recipient().String()}, strings.NewReader(plaintext), &encrypted); err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	// write identity to a temp file for the age CLI
	identityFile := t.TempDir() + "/identity"
	if err := writeFile(identityFile, []byte(id.String())); err != nil {
		t.Fatalf("write identity file: %v", err)
	}

	cmd := exec.Command(ageBin, "-d", "-i", identityFile)
	cmd.Stdin = &encrypted

	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			t.Fatalf("age -d -i: %v\nstderr: %s", err, ee.Stderr)
		}
		t.Fatalf("age -d -i: %v", err)
	}

	if string(out) != plaintext {
		t.Fatalf("age CLI decrypted %q, want %q", string(out), plaintext)
	}
}

func TestInteropKeyEncryptAgeCLIDecryptZcrypto(t *testing.T) {
	ageBin := ageCLI(t)

	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}

	plaintext := "key-encrypted by age CLI"

	cmd := exec.Command(ageBin, "-e", "-r", id.Recipient().String())
	cmd.Stdin = strings.NewReader(plaintext)

	encrypted, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			t.Fatalf("age -e -r: %v\nstderr: %s", err, ee.Stderr)
		}
		t.Fatalf("age -e -r: %v", err)
	}

	var decrypted bytes.Buffer
	if err := zcrypto.DecryptAgeKey(id.String(), bytes.NewReader(encrypted), &decrypted); err != nil {
		t.Fatalf("decrypt: %v", err)
	}

	if decrypted.String() != plaintext {
		t.Fatalf("got %q, want %q", decrypted.String(), plaintext)
	}
}

func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0o600)
}
