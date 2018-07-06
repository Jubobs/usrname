package sites

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type Site interface {
	Name() string
	Home() string
	CheckValid(username string) []Violation
	CheckAvailable(client Client) func(string) (bool, error)
	WhitelistChars() *unicode.RangeTable
}

type checker func(string) Violation

func CheckLongerThan(min int) checker {
	return func(username string) (v Violation) {
		count := utf8.RuneCountInString(username)
		if count < min {
			v = &TooShort{
				Min:    min,
				Actual: count,
			}
		}
		return
	}
}

func CheckOnlyContains(whitelist *unicode.RangeTable) checker {
	return func(username string) (v Violation) {
		var ii []int
		for i, r := range username {
			if !unicode.In(r, whitelist) {
				ii = append(ii, i)
			}
		}
		if len(ii) != 0 {
			v = &IllegalChars{
				At:        ii,
				Whitelist: whitelist,
			}
		}
		return
	}
}

func CheckNotContains(sub string) checker {
	return func(username string) (v Violation) {
		if i := strings.Index(strings.ToLower(username), sub); i != -1 {
			v = &IllegalSubstring{
				Sub: sub,
				At:  i,
			}
		}
		return
	}
}

func CheckShorterThan(max int) checker {
	return func(username string) (v Violation) {
		count := utf8.RuneCountInString(username)
		if max < count {
			v = &TooLong{
				Max:    max,
				Actual: count,
			}
		}
		return
	}
}

func CheckAll(username string, fs ...checker) []Violation {
	vv := []Violation{}
	for _, f := range fs {
		if v := f(username); v != nil {
			vv = append(vv, v)
		}
	}
	return vv
}

// type resultsByName map[string]error

// type result struct {
// 	nc  ValidNameChecker
// 	err error
// }

// // find better name for this method
// func UniversalChecker(client Client, checkers []ValidNameChecker) func(string) resultsByName {
// 	n := len(checkers)
// 	return func(username string) resultsByName {
// 		ch := make(chan *result, n)

// 		for _, checker := range checkers {
// 			go func(nc ValidNameChecker) {
// 				violations := nc.Validate(username)
// 				if len(violations) > 0 {
// 					ch <- &result{nc, errors.New("invalid username")} // TODO: tidy up
// 				}
// 				_, err := nc.Check(client, username)
// 				ch <- &result{nc, err}
// 			}(checker)
// 		}

// 		waitChecks := n
// 		res := make(map[string]error)
// 		for waitChecks > 0 {
// 			r := <-ch
// 			res[r.nc.Name()] = r.err
// 			waitChecks--
// 		}

// 		return res
// 	}
// }
