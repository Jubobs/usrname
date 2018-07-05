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
}

func All() []Site {
	return []Site{
		// Facebook(),
		// GitHub(),
		// Instagram(),
		NewTwitter(),
	}
}

func checkTooShort(s string, min int) (v Violation) {
	count := utf8.RuneCountInString(s)
	if count < min {
		v = &TooShort{
			Min:    min,
			Actual: count,
		}
	}
	return
}

func checkIllegalChars(s string, whitelist *unicode.RangeTable) (v Violation) {
	var inds []int
	for i, r := range s {
		if !unicode.In(r, whitelist) {
			inds = append(inds, i)
		}
	}
	if len(inds) != 0 {
		v = &IllegalChars{
			At:        inds,
			Whitelist: whitelist,
		}
	}
	return
}

func checkIllegalSubstring(s string, sub string) (v Violation) {
	if i := strings.Index(strings.ToLower(s), sub); i != -1 {
		v = &IllegalSubstring{
			Sub: sub,
			At:  i,
		}
	}
	return
}

func checkTooLong(s string, max int) (v Violation) {
	count := utf8.RuneCountInString(s)
	if max < count {
		v = &TooLong{
			Max:    max,
			Actual: count,
		}
	}
	return
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
