package domain

import (
	"errors"
	"net/mail"
	"regexp"
)

var (
	validChars    = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+=\-]+$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9._]{1,18})[a-zA-Z0-9]$`)
)

func (user *User) Validate() error {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return errors.New("Invalid request")
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return errors.New("The mail is incorrect")
	}

	if len([]rune(user.Password)) < 6 {
		return errors.New("The password must be at least 6 characters long.")
	}

	if !validChars.MatchString(user.Password) {
		return errors.New("Invalid password")
	}

	if !usernameRegex.MatchString(user.Username) {
		return errors.New("Invalid username")
	}

	return nil
}
