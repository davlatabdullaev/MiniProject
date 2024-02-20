package check

import (
	"errors"
	"unicode"
)

func PhoneNumber(phone string) bool {
	for _, r := range phone {
		if r == '+' {
			continue
		} else if !unicode.IsNumber(r) {
			return false
		}
	}

	return true
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password length should be more than 8")
	}

	return nil
}
