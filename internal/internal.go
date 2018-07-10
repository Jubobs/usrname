package internal

import (
	"regexp"
	"unicode"
	"unicode/utf8"

	"github.com/jubobs/usrname"
)

type checker func(string) usrname.Violation

func CheckLongerThan(min int) checker {
	return func(username string) (v usrname.Violation) {
		count := utf8.RuneCountInString(username)
		if count < min {
			v = &usrname.TooShort{
				Min:    min,
				Actual: count,
			}
		}
		return
	}
}

func CheckOnlyContains(whitelist *unicode.RangeTable) checker {
	return func(username string) (v usrname.Violation) {
		var ii []int
		for i, r := range username {
			if !unicode.In(r, whitelist) {
				ii = append(ii, i)
			}
		}
		if len(ii) != 0 {
			v = &usrname.IllegalChars{
				At:        ii,
				Whitelist: whitelist,
			}
		}
		return
	}
}

func CheckNotMatches(re *regexp.Regexp) checker {
	return func(username string) (v usrname.Violation) {
		if ii := re.FindStringIndex(username); ii != nil {
			v = &usrname.IllegalSubstring{
				Pattern: re.String(),
				At:      ii,
			}
		}
		return
	}
}

func CheckShorterThan(max int) checker {
	return func(username string) (v usrname.Violation) {
		count := utf8.RuneCountInString(username)
		if max < count {
			v = &usrname.TooLong{
				Max:    max,
				Actual: count,
			}
		}
		return
	}
}

func CheckAll(username string, fs ...checker) []usrname.Violation {
	vv := []usrname.Violation{}
	for _, f := range fs {
		if v := f(username); v != nil {
			vv = append(vv, v)
		}
	}
	return vv
}
