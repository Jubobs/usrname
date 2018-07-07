package sites

import "unicode"

type Site interface {
	Name() string
	Home() string
	CheckValid(username string) []Violation
	CheckAvailable(client Client) func(string) (bool, error)
	Whitelist() *unicode.RangeTable
}
