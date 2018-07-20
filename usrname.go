package usrname

import (
	"regexp"
	"unicode"
)

type Status string

const (
	UnknownStatus Status = "unknown"
	Invalid       Status = "invalid"
	Unavailable   Status = "unavailable"
	Available     Status = "available"
)

type Result struct {
	Username string
	Checker  Checker
	Status   Status
	Message  string
}

type Site interface {
	Name() string
	Link(username string) string
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
