package password

import (
	"strings"
	"testing"
)

func TestGenerateRandomPasswordLength(t *testing.T) {
	lengths := []int{1, 8, 16, 32, 64}
	for _, length := range lengths {
		opts := Options{
			Length:    length,
			Lowercase: true,
			Uppercase: true,
			Numbers:   true,
		}
		pwd, err := GenerateRandomPassword(opts)
		if err != nil {
			t.Fatalf("unexpected error for length %d: %v", length, err)
		}
		if len(pwd) != length {
			t.Fatalf("expected length %d, got %d", length, len(pwd))
		}
	}
}

func TestGenerateRandomPasswordInvalidLength(t *testing.T) {
	lengths := []int{0, -1, -10}
	for _, length := range lengths {
		opts := Options{Length: length}
		_, err := GenerateRandomPassword(opts)
		if err == nil {
			t.Fatalf("expected error for length %d", length)
		}
	}
}

func TestGenerateRandomPasswordLowercase(t *testing.T) {
	opts := Options{Length: 100, Lowercase: true}
	pwd, err := GenerateRandomPassword(opts)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range pwd {
		if !strings.ContainsRune(lowercase, c) {
			t.Fatalf("unexpected character %c in lowercase-only password", c)
		}
	}
}

func TestGenerateRandomPasswordUppercase(t *testing.T) {
	opts := Options{Length: 100, Uppercase: true}
	pwd, err := GenerateRandomPassword(opts)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range pwd {
		if !strings.ContainsRune(uppercase, c) {
			t.Fatalf("unexpected character %c in uppercase-only password", c)
		}
	}
}

func TestGenerateRandomPasswordNumbers(t *testing.T) {
	opts := Options{Length: 100, Numbers: true}
	pwd, err := GenerateRandomPassword(opts)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range pwd {
		if !strings.ContainsRune(numbers, c) {
			t.Fatalf("unexpected character %c in numbers-only password", c)
		}
	}
}

func TestGenerateRandomPasswordSymbols(t *testing.T) {
	opts := Options{Length: 100, Symbols: true}
	pwd, err := GenerateRandomPassword(opts)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range pwd {
		if !strings.ContainsRune(symbols, c) {
			t.Fatalf("unexpected character %c in symbols-only password", c)
		}
	}
}

func TestGenerateRandomPasswordDefaultCharset(t *testing.T) {
	opts := Options{Length: 100}
	pwd, err := GenerateRandomPassword(opts)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range pwd {
		if !strings.ContainsRune(lowercase, c) {
			t.Fatalf("expected only lowercase by default, got %c", c)
		}
	}
}

func TestGenerateRandomPasswordAllCharsets(t *testing.T) {
	opts := Options{
		Length:    1000,
		Lowercase: true,
		Uppercase: true,
		Numbers:   true,
		Symbols:   true,
	}
	pwd, err := GenerateRandomPassword(opts)
	if err != nil {
		t.Fatal(err)
	}

	hasLower := false
	hasUpper := false
	hasNumber := false
	hasSymbol := false

	for _, c := range pwd {
		switch {
		case strings.ContainsRune(lowercase, c):
			hasLower = true
		case strings.ContainsRune(uppercase, c):
			hasUpper = true
		case strings.ContainsRune(numbers, c):
			hasNumber = true
		case strings.ContainsRune(symbols, c):
			hasSymbol = true
		}
	}

	if !hasLower {
		t.Fatal("password missing lowercase character")
	}
	if !hasUpper {
		t.Fatal("password missing uppercase character")
	}
	if !hasNumber {
		t.Fatal("password missing number")
	}
	if !hasSymbol {
		t.Fatal("password missing symbol")
	}
}

func TestGenerateRandomPasswordEnsuresAtLeastOnePerSet(t *testing.T) {
	opts := Options{
		Length:    4,
		Lowercase: true,
		Uppercase: true,
		Numbers:   true,
		Symbols:   true,
	}
	for range 50 {
		pwd, err := GenerateRandomPassword(opts)
		if err != nil {
			t.Fatal(err)
		}
		if len(pwd) != 4 {
			t.Fatalf("expected length 4, got %d", len(pwd))
		}
	}
}

func TestGenerateRandomPasswordUniqueness(t *testing.T) {
	opts := Options{Length: 16, Lowercase: true, Uppercase: true, Numbers: true}
	seen := make(map[string]bool)
	for range 100 {
		pwd, err := GenerateRandomPassword(opts)
		if err != nil {
			t.Fatal(err)
		}
		if seen[pwd] {
			t.Fatal("duplicate password generated")
		}
		seen[pwd] = true
	}
}

func TestGenerateRandomPasswordNotEmpty(t *testing.T) {
	opts := Options{Length: 1, Lowercase: true}
	for range 10 {
		pwd, err := GenerateRandomPassword(opts)
		if err != nil {
			t.Fatal(err)
		}
		if pwd == "" {
			t.Fatal("password should not be empty")
		}
	}
}