package internal

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/jubobs/whocanibe/sites"
)

type checker func(string) sites.Violation

func CheckLongerThan(min int) checker {
	return func(username string) (v sites.Violation) {
		count := utf8.RuneCountInString(username)
		if count < min {
			v = &sites.TooShort{
				Min:    min,
				Actual: count,
			}
		}
		return
	}
}

func CheckOnlyContains(whitelist *unicode.RangeTable) checker {
	return func(username string) (v sites.Violation) {
		var ii []int
		for i, r := range username {
			if !unicode.In(r, whitelist) {
				ii = append(ii, i)
			}
		}
		if len(ii) != 0 {
			v = &sites.IllegalChars{
				At:        ii,
				Whitelist: whitelist,
			}
		}
		return
	}
}

func CheckNotContains(sub string) checker {
	return func(username string) (v sites.Violation) {
		if i := strings.Index(strings.ToLower(username), sub); i != -1 {
			v = &sites.IllegalSubstring{
				Sub: sub,
				At:  i,
			}
		}
		return
	}
}

func CheckShorterThan(max int) checker {
	return func(username string) (v sites.Violation) {
		count := utf8.RuneCountInString(username)
		if max < count {
			v = &sites.TooLong{
				Max:    max,
				Actual: count,
			}
		}
		return
	}
}

func CheckAll(username string, fs ...checker) []sites.Violation {
	vv := []sites.Violation{}
	for _, f := range fs {
		if v := f(username); v != nil {
			vv = append(vv, v)
		}
	}
	return vv
}
