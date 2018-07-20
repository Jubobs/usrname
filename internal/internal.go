package internal

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/jubobs/usrname"
)

type validate func(string) usrname.Violation

func CheckLongerThan(min int) validate {
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

func CheckOnlyContains(whitelist *unicode.RangeTable) validate {
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

func CheckIllegalPrefix(prefix string) validate {
	return func(username string) (v usrname.Violation) {
		if strings.HasPrefix(username, prefix) {
			v = &usrname.IllegalPrefix{
				Pattern: prefix,
			}
		}
		return
	}
}

func CheckIllegalSubstring(sub string) validate {
	return func(username string) (v usrname.Violation) {
		if strings.Contains(username, sub) {
			v = &usrname.IllegalSubstring{
				Pattern: sub,
			}
		}
		return
	}
}

func CheckIllegalSuffix(suffix string) validate {
	return func(username string) (v usrname.Violation) {
		if strings.HasSuffix(username, suffix) {
			v = &usrname.IllegalSuffix{
				Pattern: suffix,
			}
		}
		return
	}
}

func CheckNotMatches(re *regexp.Regexp) validate {
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

func CheckShorterThan(max int) validate {
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

func CheckAll(username string, fs ...validate) []usrname.Violation {
	vv := []usrname.Violation{}
	for _, f := range fs {
		if v := f(username); v != nil {
			vv = append(vv, v)
		}
	}
	return vv
}
