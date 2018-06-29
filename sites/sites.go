package sites

import (
	"errors"
)

type Namer interface {
	Name() string
}

type Validator interface {
	Validate(username string) []string // TODO: introduce Violation type
}

type Checker interface {
	Check(client Client, username string) (bool, error)
}

type NameChecker interface {
	Namer
	Checker
}

type ValidNameChecker interface {
	Validator
	NameChecker
}

func All() []ValidNameChecker {
	return []ValidNameChecker{
		// Facebook(),
		// GitHub(),
		// Instagram(),
		Twitter(),
	}
}

type resultsByName map[string]error

type result struct {
	nc  ValidNameChecker
	err error
}

// find better name for this method
func UniversalChecker(client Client, checkers []ValidNameChecker) func(string) resultsByName {
	n := len(checkers)
	return func(username string) resultsByName {
		ch := make(chan *result, n)

		for _, checker := range checkers {
			go func(nc ValidNameChecker) {
				violations := nc.Validate(username)
				if len(violations) > 0 {
					ch <- &result{nc, errors.New("invalid username")} // TODO: tidy up
				}
				_, err := nc.Check(client, username)
				ch <- &result{nc, err}
			}(checker)
		}

		waitChecks := n
		res := make(map[string]error)
		for waitChecks > 0 {
			r := <-ch
			res[r.nc.Name()] = r.err
			waitChecks--
		}

		return res
	}
}
