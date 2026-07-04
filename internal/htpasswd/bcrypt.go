package htpasswd

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/unicode/norm"
)

const (
	DefaultCost = 10
	MinCost     = 4
	MaxCost     = 31
)

var (
	ErrUsernameEmpty = errors.New("username must not be empty")
	ErrUsernameColon = errors.New("username must not contain colon")
	ErrPasswordEmpty = errors.New("password must not be empty")
	ErrCostInvalid   = fmt.Errorf("cost must be between %d and %d", MinCost, MaxCost)
)

func Generate(username, password string, cost int) (string, error) {
	username = norm.NFC.String(username)
	password = norm.NFC.String(password)

	if username == "" {
		return "", ErrUsernameEmpty
	}
	if strings.Contains(username, ":") {
		return "", ErrUsernameColon
	}
	if password == "" {
		return "", ErrPasswordEmpty
	}
	if cost < MinCost || cost > MaxCost {
		return "", ErrCostInvalid
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	s := string(hash)
	s = "$2y" + s[3:]

	return username + ":" + s, nil
}