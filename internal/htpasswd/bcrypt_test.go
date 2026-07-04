package htpasswd

import (
	"strings"
	"testing"
)

func TestGenerateSuccess(t *testing.T) {
	line, err := Generate("paulo", "mysecret", DefaultCost)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		t.Fatalf("expected user:hash format, got %q", line)
	}
	if parts[0] != "paulo" {
		t.Fatalf("expected username 'paulo', got %q", parts[0])
	}
	if !strings.HasPrefix(parts[1], "$2y$10$") {
		t.Fatalf("expected bcrypt hash with cost 10, got %q", parts[1])
	}
}

func TestGenerateSuccessCustomCost(t *testing.T) {
	line, err := Generate("paulo", "mysecret", 12)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parts := strings.SplitN(line, ":", 2)
	if !strings.HasPrefix(parts[1], "$2y$12$") {
		t.Fatalf("expected bcrypt hash with cost 12, got %q", parts[1])
	}
}

func TestGenerateInvalidCostBelow(t *testing.T) {
	_, err := Generate("user", "password", MinCost-1)
	if err != ErrCostInvalid {
		t.Fatalf("expected ErrCostInvalid, got %v", err)
	}
}

func TestGenerateInvalidCostAbove(t *testing.T) {
	_, err := Generate("user", "password", MaxCost+1)
	if err != ErrCostInvalid {
		t.Fatalf("expected ErrCostInvalid, got %v", err)
	}
}

func TestGenerateEmptyUsername(t *testing.T) {
	_, err := Generate("", "password", DefaultCost)
	if err != ErrUsernameEmpty {
		t.Fatalf("expected ErrUsernameEmpty, got %v", err)
	}
}

func TestGenerateEmptyPassword(t *testing.T) {
	_, err := Generate("user", "", DefaultCost)
	if err != ErrPasswordEmpty {
		t.Fatalf("expected ErrPasswordEmpty, got %v", err)
	}
}

func TestGenerateUsernameWithColon(t *testing.T) {
	_, err := Generate("user:name", "password", DefaultCost)
	if err != ErrUsernameColon {
		t.Fatalf("expected ErrUsernameColon, got %v", err)
	}
}

func TestGenerateUsernameWithColonAtEnd(t *testing.T) {
	_, err := Generate("username:", "password", DefaultCost)
	if err != ErrUsernameColon {
		t.Fatalf("expected ErrUsernameColon, got %v", err)
	}
}

func TestGenerateDifferentPasswords(t *testing.T) {
	line1, err := Generate("user", "password1", DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	line2, err := Generate("user", "password2", DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	hash1 := strings.SplitN(line1, ":", 2)[1]
	hash2 := strings.SplitN(line2, ":", 2)[1]

	if hash1 == hash2 {
		t.Fatal("different passwords should produce different hashes")
	}
}

func TestGenerateDeterministicHash(t *testing.T) {
	line1, err := Generate("user", "samepassword", DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	line2, err := Generate("user", "samepassword", DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	hash1 := strings.SplitN(line1, ":", 2)[1]
	hash2 := strings.SplitN(line2, ":", 2)[1]

	if hash1 == hash2 {
		t.Fatal("bcrypt should produce different hashes even for same input")
	}
}

func TestGenerateSpecialCharsInUsername(t *testing.T) {
	usernames := []string{
		"usér",
		"user name",
		"user.name",
		"user_name",
		"user-name",
	}
	for _, u := range usernames {
		line, err := Generate(u, "password", DefaultCost)
		if err != nil {
			t.Fatalf("unexpected error for username %q: %v", u, err)
		}
		if !strings.HasPrefix(line, u+":") {
			t.Fatalf("expected prefix %q, got %q", u+":", line)
		}
	}
}

func TestGenerateOutputFormat(t *testing.T) {
	line, err := Generate("testuser", "testpass", DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(line, ":") {
		t.Fatal("output must contain a colon separator")
	}

	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		t.Fatal("output must have exactly two parts separated by colon")
	}
	if parts[0] != "testuser" {
		t.Fatalf("expected username 'testuser', got %q", parts[0])
	}
	if len(parts[1]) == 0 {
		t.Fatal("hash must not be empty")
	}
}