package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"mkhtpasswd/internal/htpasswd"
	"mkhtpasswd/internal/password"
)

const version = "0.1.0"

type Result struct {
	Line     string
	Password string
}

type Config struct {
	ShowVersion bool
	ForceGUI    bool
	User        string
	Password    string
	Random      bool
	Length      int
	Lower       bool
	Upper       bool
	Numbers     bool
	Symbols     bool
	ShowPass    bool
	Cost        int
}

func Parse(args []string) (Config, bool) {
	if len(args) == 0 {
		return Config{}, true
	}

	fs := flag.NewFlagSet("mkhtpasswd", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: mkhtpasswd [flags]\n\n")
		fs.PrintDefaults()
	}

	gui := fs.Bool("gui", false, "Force opening the GUI")
	user := fs.String("user", "", "Username")
	password := fs.String("password", "", "Password to encrypt")
	random := fs.Bool("random", false, "Generate a random password")
	length := fs.Int("length", 20, "Password length")
	lower := fs.Bool("lower", true, "Use lowercase letters")
	upper := fs.Bool("upper", true, "Use uppercase letters")
	numbers := fs.Bool("numbers", true, "Use numbers")
	symbols := fs.Bool("symbols", false, "Use symbols")
	showPass := fs.Bool("show-password", false, "Display the generated password")
	cost := fs.Int("cost", 10, "bcrypt cost")
	showVersion := fs.Bool("version", false, "Show version")

	if err := fs.Parse(args); err != nil {
		return Config{}, false
	}

	return Config{
		ShowVersion: *showVersion,
		ForceGUI:    *gui,
		User:        *user,
		Password:    *password,
		Random:      *random,
		Length:      *length,
		Lower:       *lower,
		Upper:       *upper,
		Numbers:     *numbers,
		Symbols:     *symbols,
		ShowPass:    *showPass,
		Cost:        *cost,
	}, true
}

func Run(cfg Config) (Result, error) {
	if cfg.ShowVersion {
		fmt.Println("mkhtpasswd version", version)
		os.Exit(0)
	}

	var pwd string

	if cfg.Random {
		opts := password.Options{
			Length:    cfg.Length,
			Lowercase: cfg.Lower,
			Uppercase: cfg.Upper,
			Numbers:   cfg.Numbers,
			Symbols:   cfg.Symbols,
		}
		var err error
		pwd, err = password.GenerateRandomPassword(opts)
		if err != nil {
			return Result{}, fmt.Errorf("generating random password: %w", err)
		}
	} else {
		pwd = cfg.Password
	}

	if cfg.User == "" {
		return Result{}, errors.New("username is required")
	}
	if pwd == "" {
		return Result{}, errors.New("password is required")
	}

	line, err := htpasswd.Generate(cfg.User, pwd)
	if err != nil {
		return Result{}, err
	}

	return Result{Line: line, Password: pwd}, nil
}

func ShowResult(cfg Config, res Result) {
	fmt.Println(res.Line)
	if cfg.Random && cfg.ShowPass {
		fmt.Fprintln(os.Stderr, "Generated password:", res.Password)
	}
}