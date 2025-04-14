package domain

import (
	"net/mail"
	"regexp"
)

type Email string

func (e Email) IsValid() bool {
	emailRegexp := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegexp.MatchString(string(e)) {
		return false
	}

	_, err := mail.ParseAddress(string(e))

	return err == nil
}
