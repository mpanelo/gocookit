package models

const (
	ErrNotFound                 = privateError("resource not found")
	ErrIDInvalid                = privateError("ID has an invalid value")
	ErrUserPasswordHashRequired = privateError("password hash is required")
	ErrUserPasswordRequired     = privateError("password is required")                        // TODO convert to a public error
	ErrUserEmailRequired        = privateError("email is required")                           // TODO convert to a public error
	ErrUserNameRequired         = privateError("full name is required")                       // TODO convert to a public error
	ErrUserPasswordTooShort     = privateError("password must be at least 8 characters long") // TODO convert to a public error
	ErrUserEmailInvalid         = privateError("email provided has an invalid format")        // TODO convert to a public error
	ErrUserEmailTaken           = privateError("email is already taken")                      // TODO convert to a public error
)

type privateError string

func (me privateError) Error() string {
	return string(me)
}
