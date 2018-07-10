package usrname

import (
	"regexp"
	"unicode"
)

type Site interface {
	Name() string
	Home() string
	ProfilePage(username string) string
}

type Validator interface {
	Site
	Validate(username string) []Violation
	IllegalPattern() *regexp.Regexp
	Whitelist() *unicode.RangeTable
}

type Checker interface {
	Validator
	Check(client Client) func(string) (bool, error)
}
