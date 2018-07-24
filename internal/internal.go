package internal

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/jubobs/usrname"
)

type validate1 func(string) usrname.Violation

func CheckLongerThan(min int) validate1 {
	return func(username string) usrname.Violation {
		count := utf8.RuneCountInString(username)
		if count < min {
			return &usrname.TooShort{
				Min:    min,
				Actual: count,
			}
		}
		return nil
	}
}

func CheckOnlyContains(whitelist *unicode.RangeTable) validate1 {
	return func(username string) usrname.Violation {
		var ii []int
		for i, r := range username {
			if !unicode.In(r, whitelist) {
				ii = append(ii, i)
			}
		}
		if len(ii) != 0 {
			return &usrname.IllegalChars{
				At:        ii,
				Whitelist: whitelist,
			}
		}
		return nil
	}
}

func CheckIllegalPrefix(prefix string) validate1 {
	return func(username string) usrname.Violation {
		if strings.HasPrefix(username, prefix) {
			return &usrname.IllegalPrefix{
				Pattern: prefix,
			}
		}
		return nil
	}
}

func CheckIllegalSubstring(sub string) validate1 {
	return func(username string) (v usrname.Violation) {
		if strings.Contains(username, sub) {
			v = &usrname.IllegalSubstring{
				Pattern: sub,
			}
		}
		return
	}
}

func CheckIllegalSuffix(suffix string) validate1 {
	return func(username string) (v usrname.Violation) {
		if strings.HasSuffix(username, suffix) {
			v = &usrname.IllegalSuffix{
				Pattern: suffix,
			}
		}
		return
	}
}

func CheckNotMatches(re *regexp.Regexp) validate1 {
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

func CheckShorterThan(max int) validate1 {
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

func CheckAll(username string, fs ...validate1) []usrname.Violation {
	vv := []usrname.Violation{}
	for _, f := range fs {
		if v := f(username); v != nil {
			vv = append(vv, v)
		}
	}
	return vv
}

func IsTimeout(err error) bool {
	type timeout interface {
		Timeout() bool
		error
	}
	err1, ok := err.(timeout)
	return ok && err1.Timeout()
}
