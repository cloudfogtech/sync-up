package check

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func UsernamePassword(username, password string) error {
	return validation.Errors{
		"username": validation.Validate(username, validation.Required),
		"password": validation.Validate(password, validation.Required),
	}.Filter()
}
