package zcrypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

const (
	// totpPeriod is the time step in seconds (RFC 6238 default).
	totpPeriod = 30
)

// TOTPCode generates the current 6-digit TOTP code for the given base32-encoded secret.
// The secret is case-insensitive and padding is optional.
func TOTPCode(secret string) (string, error) {
	return TOTPCodeAt(secret, time.Now())
}

// TOTPCodeAt generates a 6-digit TOTP code for the given base32-encoded secret at time t.
// Implements RFC 6238 with HMAC-SHA1, 30-second period, and 6-digit output.
func TOTPCodeAt(secret string, t time.Time) (string, error) {
	key, err := decodeSecret(secret)
	if err != nil {
		return "", fmt.Errorf("decode secret: %w", err)
	}

	counter := uint64(t.Unix()) / totpPeriod

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], counter)

	mac := hmac.New(sha1.New, key)
	mac.Write(buf[:])
	h := mac.Sum(nil)

	offset := h[len(h)-1] & 0x0f
	code := binary.BigEndian.Uint32(h[offset:offset+4]) & 0x7fffffff
	code %= 1_000_000

	return fmt.Sprintf("%06d", code), nil
}

// decodeSecret decodes a base32 secret, handling case and optional padding.
func decodeSecret(secret string) ([]byte, error) {
	s := strings.ToUpper(strings.TrimSpace(secret))
	s = strings.ReplaceAll(s, " ", "")

	// add padding if missing
	if pad := len(s) % 8; pad != 0 {
		s += strings.Repeat("=", 8-pad)
	}

	b, err := base32.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}
