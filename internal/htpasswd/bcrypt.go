package htpasswd

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const DefaultCost = 10

var (
	ErrUsernameEmpty = errors.New("username must not be empty")
	ErrUsernameColon = errors.New("username must not contain colon")
	ErrPasswordEmpty = errors.New("password must not be empty")
)

func Generate(username, password string) (string, error) {
	if username == "" {
		return "", ErrUsernameEmpty
	}
	if strings.Contains(username, ":") {
		return "", ErrUsernameColon
	}
	if password == "" {
		return "", ErrPasswordEmpty
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}

	s := string(hash)
	s = "$2y" + s[3:]

	return username + ":" + s, nil
}