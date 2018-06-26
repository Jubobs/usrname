package sites

import (
	"fmt"
)

type unavailableUsernameError struct {
	Namer
	username string
}

func (err *unavailableUsernameError) Error() string {
	return fmt.Sprintf("sites: %s is unavailable on %s", err.username, err.Name())
}

func IsUnavailableUsernameError(err error) bool {
	_, ok := err.(*unavailableUsernameError)
	return ok
}

type invalidUsernameError struct {
	Namer
	username string
}

func (err *invalidUsernameError) Error() string {
	return fmt.Sprintf("sites: %s is invalid on %s", err.username, err.Name())
}

func IsInvalidUsernameError(err error) bool {
	_, ok := err.(*invalidUsernameError)
	return ok
}

type unexpectedError struct {
	cause error
}

func (err *unexpectedError) Error() string {
	return fmt.Sprintf("sites: unexpected error: %v", err.cause)
}

func IsUnexpectedError(err error) bool {
	_, ok := err.(*invalidUsernameError)
	return ok
}
