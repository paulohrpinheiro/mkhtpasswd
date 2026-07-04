package password

import (
	"crypto/rand"
	"errors"
	"math/big"
)

type Options struct {
	Length int

	Lowercase bool
	Uppercase bool
	Numbers   bool
	Symbols   bool
}

const (
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers   = "0123456789"
	symbols   = "!@#$%^&*()-_=+[]{}|;:,.<>?/~"
)

func GenerateRandomPassword(opts Options) (string, error) {
	if opts.Length <= 0 {
		return "", errors.New("password length must be positive")
	}

	charsets := make([]string, 0, 4)
	if opts.Lowercase {
		charsets = append(charsets, lowercase)
	}
	if opts.Uppercase {
		charsets = append(charsets, uppercase)
	}
	if opts.Numbers {
		charsets = append(charsets, numbers)
	}
	if opts.Symbols {
		charsets = append(charsets, symbols)
	}
	if len(charsets) == 0 {
		charsets = append(charsets, lowercase)
	}

	allChars := ""
	for _, cs := range charsets {
		allChars += cs
	}

	need := opts.Length

	password := make([]byte, need)

	for i := range password {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "", err
		}
		password[i] = allChars[n.Int64()]
	}

	for _, cs := range charsets {
		hasOne := false
		for _, b := range password {
			for _, c := range []byte(cs) {
				if b == c {
					hasOne = true
					break
				}
			}
			if hasOne {
				break
			}
		}
		if !hasOne {
			n, err := rand.Int(rand.Reader, big.NewInt(int64(opts.Length)))
			if err != nil {
				return "", err
			}
			cn, err := rand.Int(rand.Reader, big.NewInt(int64(len(cs))))
			if err != nil {
				return "", err
			}
			password[n.Int64()] = cs[cn.Int64()]
		}
	}

	return string(password), nil
}