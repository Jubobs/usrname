package usrname

import (
	"regexp"
	"unicode"
)

type Status int

const (
	UnknownStatus Status = iota
	Invalid
	Unavailable
	Available
)

type Result struct {
	Username string
	Checker  Checker
	Status   Status
	Message  string
}

type Site interface {
	Name() string
	Url(username string) string
}

type Validator interface {
	Site
	Validate(username string) []Violation
	IllegalPattern() *regexp.Regexp
	Whitelist() *unicode.RangeTable
}

type Checker interface {
	Validator
	Check(client Client) func(string) Result
}
