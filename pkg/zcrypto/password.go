package zcrypto

import (
	"crypto/rand"
	"math/big"
)

// password character classes
const (
	lowerChars  = "abcdefghijklmnopqrstuvwxyz"
	upperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars  = "0123456789"
	symbolChars = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

// PasswordOption configures password generation.
type PasswordOption func(*passwordConfig)

type passwordConfig struct {
	charset string
	symbols bool
}

func defaultPasswordConfig() passwordConfig {
	return passwordConfig{
		symbols: true,
	}
}

// WithoutSymbols excludes symbol characters from generated passwords.
func WithoutSymbols() PasswordOption {
	return func(c *passwordConfig) {
		c.symbols = false
	}
}

// WithCharset overrides the full character set used for password generation.
// When set, character class guarantees (lower, upper, digit, symbol) are
// disabled — the caller controls everything.
func WithCharset(chars string) PasswordOption {
	return func(c *passwordConfig) {
		c.charset = chars
	}
}

// GeneratePassword produces a cryptographically random password of the given
// length. By default it guarantees at least one character from each class
// (lower, upper, digit, symbol) and uses Fisher-Yates shuffle for uniform
// distribution. Minimum length is 4; shorter values are clamped.
//
// Panics if crypto/rand fails (unrecoverable).
func GeneratePassword(length int, opts ...PasswordOption) string {
	cfg := defaultPasswordConfig()
	for _, o := range opts {
		o(&cfg)
	}

	// custom charset: no guarantees, just fill and shuffle
	if cfg.charset != "" {
		if length < 1 {
			length = 1
		}
		return generateFromCharset(length, cfg.charset)
	}

	if length < 4 {
		length = 4
	}

	var classes []string
	classes = append(classes, lowerChars, upperChars, digitChars)
	if cfg.symbols {
		classes = append(classes, symbolChars)
	}

	// build the full charset
	var charset string
	for _, c := range classes {
		charset += c
	}

	buf := make([]byte, length)

	// guarantee one from each class
	for i, class := range classes {
		buf[i] = pickByte(class)
	}

	// fill remainder from full charset
	for i := len(classes); i < length; i++ {
		buf[i] = pickByte(charset)
	}

	// Fisher-Yates shuffle
	shuffle(buf)

	return string(buf)
}

// generateFromCharset fills a buffer from an arbitrary charset and shuffles.
func generateFromCharset(length int, charset string) string {
	buf := make([]byte, length)
	for i := range buf {
		buf[i] = pickByte(charset)
	}
	shuffle(buf)
	return string(buf)
}

// shuffle performs Fisher-Yates shuffle using crypto/rand.
func shuffle(buf []byte) {
	for i := len(buf) - 1; i > 0; i-- {
		j := randIntn(i + 1)
		buf[i], buf[j] = buf[j], buf[i]
	}
}

// pickByte returns a random byte from a string.
func pickByte(s string) byte {
	return s[randIntn(len(s))]
}

// randIntn returns a cryptographically random int in [0, n).
func randIntn(n int) int {
	v, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		panic("crypto/rand: " + err.Error())
	}
	return int(v.Int64())
}
