package zcrypto

import (
	"strings"
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	lengths := []int{4, 8, 20, 64, 128}
	for _, n := range lengths {
		pw := GeneratePassword(n)
		if len(pw) != n {
			t.Fatalf("GeneratePassword(%d) returned length %d", n, len(pw))
		}
	}
}

func TestGeneratePasswordMinLength(t *testing.T) {
	tests := []struct {
		name   string
		length int
		want   int
	}{
		{"zero", 0, 4},
		{"one", 1, 4},
		{"three", 3, 4},
		{"negative", -5, 4},
		{"four", 4, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pw := GeneratePassword(tt.length)
			if len(pw) != tt.want {
				t.Fatalf("GeneratePassword(%d) length = %d, want %d", tt.length, len(pw), tt.want)
			}
		})
	}
}

func TestGeneratePasswordAllClasses(t *testing.T) {
	// generate many passwords and check each has all classes
	for i := 0; i < 50; i++ {
		pw := GeneratePassword(20)

		if !strings.ContainsAny(pw, lowerChars) {
			t.Fatalf("password %q missing lowercase", pw)
		}
		if !strings.ContainsAny(pw, upperChars) {
			t.Fatalf("password %q missing uppercase", pw)
		}
		if !strings.ContainsAny(pw, digitChars) {
			t.Fatalf("password %q missing digit", pw)
		}
		if !strings.ContainsAny(pw, symbolChars) {
			t.Fatalf("password %q missing symbol", pw)
		}
	}
}

func TestGeneratePasswordWithoutSymbols(t *testing.T) {
	for i := 0; i < 50; i++ {
		pw := GeneratePassword(20, WithoutSymbols())

		if strings.ContainsAny(pw, symbolChars) {
			t.Fatalf("password %q contains symbol with WithoutSymbols()", pw)
		}
		if !strings.ContainsAny(pw, lowerChars) {
			t.Fatalf("password %q missing lowercase", pw)
		}
		if !strings.ContainsAny(pw, upperChars) {
			t.Fatalf("password %q missing uppercase", pw)
		}
		if !strings.ContainsAny(pw, digitChars) {
			t.Fatalf("password %q missing digit", pw)
		}
	}
}

func TestGeneratePasswordWithoutSymbolsMinLength(t *testing.T) {
	// without symbols there are 3 classes, but minimum is still 4
	pw := GeneratePassword(2, WithoutSymbols())
	if len(pw) != 4 {
		t.Fatalf("length = %d, want 4", len(pw))
	}
}

func TestGeneratePasswordWithCharset(t *testing.T) {
	charset := "abc123"
	pw := GeneratePassword(30, WithCharset(charset))

	if len(pw) != 30 {
		t.Fatalf("length = %d, want 30", len(pw))
	}

	for _, c := range pw {
		if !strings.ContainsRune(charset, c) {
			t.Fatalf("password contains %q which is not in charset %q", string(c), charset)
		}
	}
}

func TestGeneratePasswordWithCharsetMinLength(t *testing.T) {
	// custom charset bypasses the 4-char minimum, clamps to 1
	pw := GeneratePassword(0, WithCharset("x"))
	if len(pw) != 1 {
		t.Fatalf("length = %d, want 1", len(pw))
	}
}

func TestGeneratePasswordRandomness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		pw := GeneratePassword(20)
		seen[pw] = true
	}

	// 100 random 20-char passwords should produce at least 2 unique values.
	// in practice they'll all be unique, but this catches a broken RNG
	// that returns the same thing every time.
	if len(seen) < 2 {
		t.Fatal("100 generated passwords are all identical")
	}
}

func TestGeneratePasswordCharacterDistribution(t *testing.T) {
	// verify the full charset is reachable by generating many passwords
	// and collecting unique characters
	chars := make(map[byte]bool)
	for i := 0; i < 500; i++ {
		pw := GeneratePassword(128)
		for j := 0; j < len(pw); j++ {
			chars[pw[j]] = true
		}
	}

	allChars := lowerChars + upperChars + digitChars + symbolChars
	for i := 0; i < len(allChars); i++ {
		if !chars[allChars[i]] {
			t.Fatalf("character %q never appeared in 500 passwords of length 128", string(allChars[i]))
		}
	}
}
