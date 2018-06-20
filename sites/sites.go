package sites

import (
	"errors"
)

type Namer interface {
	Name() string
}

type Validator interface {
	Validate(username string) error
}

type Checker interface {
	Check(client Client, username string) error
}

type NameChecker interface {
	Namer
	Checker
}

type ValidNameChecker interface {
	Validator
	NameChecker
}

var (
	ErrUnavailableUsername = errors.New("unavailable username")
	ErrInvalidUserName     = errors.New("invalid username")
)

func All() []NameChecker {
	return []NameChecker{
		Facebook(),
		GitHub(),
		Twitter(),
	}
}
