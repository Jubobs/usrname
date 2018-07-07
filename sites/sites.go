package sites

import (
	"regexp"
	"unicode"
)

type Site interface {
	Name() string
	Home() string
	CheckValid(username string) []Violation
	CheckAvailable(client Client) func(string) (bool, error)
	IllegalPattern() *regexp.Regexp
	Whitelist() *unicode.RangeTable
}
