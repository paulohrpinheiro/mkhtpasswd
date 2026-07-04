package cli

import (
	"testing"
)

func TestParseEmptyArgs(t *testing.T) {
	cfg, ok := Parse([]string{})
	if !ok {
		t.Fatal("expected ok for empty args")
	}
	if cfg.ForceGUI {
		t.Fatal("expected ForceGUI to be false")
	}
	if cfg.User != "" {
		t.Fatal("expected empty user")
	}
}

func TestParseWithUserAndPassword(t *testing.T) {
	cfg, ok := Parse([]string{"--user", "paulo", "--password", "secret"})
	if !ok {
		t.Fatal("expected ok")
	}
	if cfg.User != "paulo" {
		t.Fatalf("expected user 'paulo', got %q", cfg.User)
	}
	if cfg.Password != "secret" {
		t.Fatalf("expected password 'secret', got %q", cfg.Password)
	}
}

func TestParseRandomDefaults(t *testing.T) {
	cfg, ok := Parse([]string{"--random"})
	if !ok {
		t.Fatal("expected ok")
	}
	if !cfg.Random {
		t.Fatal("expected random")
	}
	if !cfg.Lower {
		t.Fatal("expected lower default true")
	}
	if !cfg.Upper {
		t.Fatal("expected upper default true")
	}
	if !cfg.Numbers {
		t.Fatal("expected numbers default true")
	}
	if cfg.Symbols {
		t.Fatal("expected symbols default false")
	}
}

func TestParseRandomWithSymbols(t *testing.T) {
	cfg, ok := Parse([]string{"--random", "--symbols", "--length", "30"})
	if !ok {
		t.Fatal("expected ok")
	}
	if !cfg.Symbols {
		t.Fatal("expected symbols true")
	}
	if cfg.Length != 30 {
		t.Fatalf("expected length 30, got %d", cfg.Length)
	}
}

func TestParseShowVersion(t *testing.T) {
	cfg, ok := Parse([]string{"--version"})
	if !ok {
		t.Fatal("expected ok")
	}
	if !cfg.ShowVersion {
		t.Fatal("expected ShowVersion true")
	}
}

func TestParseForceGUI(t *testing.T) {
	cfg, ok := Parse([]string{"--gui"})
	if !ok {
		t.Fatal("expected ok")
	}
	if !cfg.ForceGUI {
		t.Fatal("expected ForceGUI true")
	}
}

func TestParseShowPassword(t *testing.T) {
	cfg, ok := Parse([]string{"--random", "--show-password"})
	if !ok {
		t.Fatal("expected ok")
	}
	if !cfg.ShowPass {
		t.Fatal("expected ShowPass true")
	}
}

func TestParseCost(t *testing.T) {
	cfg, ok := Parse([]string{"--cost", "12"})
	if !ok {
		t.Fatal("expected ok")
	}
	if cfg.Cost != 12 {
		t.Fatalf("expected cost 12, got %d", cfg.Cost)
	}
}

func TestRunWithUserAndPassword(t *testing.T) {
	cfg := Config{
		User:     "testuser",
		Password: "testpass",
		Cost:     10,
	}
	res, err := Run(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Line == "" {
		t.Fatal("expected non-empty line")
	}
	if res.Password != "testpass" {
		t.Fatalf("expected password 'testpass', got %q", res.Password)
	}
}

func TestRunRandom(t *testing.T) {
	cfg := Config{
		User:    "testuser",
		Random:  true,
		Length:  16,
		Lower:   true,
		Upper:   true,
		Numbers: true,
		Cost:    10,
	}
	res, err := Run(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Line == "" {
		t.Fatal("expected non-empty line")
	}
	if res.Password == "" {
		t.Fatal("expected non-empty password")
	}
}

func TestRunEmptyUser(t *testing.T) {
	cfg := Config{
		Password: "secret",
	}
	_, err := Run(cfg)
	if err == nil {
		t.Fatal("expected error for empty user")
	}
}

func TestRunEmptyPassword(t *testing.T) {
	cfg := Config{
		User: "user",
	}
	_, err := Run(cfg)
	if err == nil {
		t.Fatal("expected error for empty password")
	}
}

func TestRunInvalidUserColon(t *testing.T) {
	cfg := Config{
		User:     "user:name",
		Password: "secret",
		Cost:     10,
	}
	_, err := Run(cfg)
	if err == nil {
		t.Fatal("expected error for username with colon")
	}
}

func TestRunRandomExceedsMaxLength(t *testing.T) {
	cfg := Config{
		User:   "testuser",
		Random: true,
		Length: 1000,
		Lower:  true,
	}
	_, err := Run(cfg)
	if err == nil {
		t.Fatal("expected error for length exceeding max")
	}
}

func TestParsePasswordStdin(t *testing.T) {
	cfg, ok := Parse([]string{"--password-stdin"})
	if !ok {
		t.Fatal("expected ok")
	}
	if !cfg.PasswordStdin {
		t.Fatal("expected PasswordStdin true")
	}
}

func TestRunPasswordAndStdinMutuallyExclusive(t *testing.T) {
	cfg := Config{
		User:          "testuser",
		Password:      "secret",
		PasswordStdin: true,
		Cost:          10,
	}
	_, err := Run(cfg)
	if err == nil {
		t.Fatal("expected error for --password and --password-stdin together")
	}
}

func TestRunRandomAndStdinMutuallyExclusive(t *testing.T) {
	cfg := Config{
		User:          "testuser",
		PasswordStdin: true,
		Random:        true,
		Cost:          10,
	}
	_, err := Run(cfg)
	if err == nil {
		t.Fatal("expected error for --random and --password-stdin together")
	}
}

func TestParsePasswordStdinFlag(t *testing.T) {
	cfg, ok := Parse([]string{"--user", "u", "--password-stdin"})
	if !ok {
		t.Fatal("expected ok")
	}
	if !cfg.PasswordStdin {
		t.Fatal("expected PasswordStdin true")
	}
}