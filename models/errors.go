package models

import (
	"strings"
)

const (
	ErrNotFound                 = privateError("resource not found")
	ErrIDInvalid                = privateError("ID has an invalid value")
	ErrUserPasswordHashRequired = privateError("password hash is required")
	ErrUserRememberHashRequired = privateError("remember hash is required")
	ErrUserRememberTooShort     = privateError("remember token must be at least 32 bytes long")
	ErrUserPasswordRequired     = publicError("password is required")
	ErrUserEmailRequired        = publicError("email is required")
	ErrUserNameRequired         = publicError("full name is required")
	ErrUserPasswordTooShort     = publicError("password must be at least 8 characters long")
	ErrUserEmailInvalid         = publicError("email provided has an invalid format")
	ErrUserEmailTaken           = publicError("email is already taken")
	ErrUserCredentialsInvalid   = publicError("email or password provided is invalid")
)

type privateError string

func (e privateError) Error() string {
	return string(e)
}

type publicError string

func (e publicError) Error() string {
	return string(e)
}

func (e publicError) Alert() string {
	words := strings.Split(e.Error(), " ")
	words[0] = strings.Title(words[0])
	return strings.Join(words, " ")
}
