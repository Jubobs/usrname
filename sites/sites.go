package sites

type Namer interface {
	Name() string
}

type Validator interface {
	Validate(username string) error
}

type Checker interface {
	Check(client Client, username string) error
}

type NameChecker interface {
	Namer
	Checker
}

type ValidNameChecker interface {
	Validator
	NameChecker
}

func All() []NameChecker {
	return []NameChecker{
		Facebook(),
		GitHub(),
		Twitter(),
	}
}

type resultsByName map[string]error

type result struct {
	nc  NameChecker
	err error
}

// find better name for this method
func UniversalChecker(client Client, checkers []NameChecker) func(string) resultsByName {
	n := len(checkers)
	return func(username string) resultsByName {
		ch := make(chan *result, n)

		for _, checker := range checkers {
			go func(nc NameChecker) {
				err := nc.Check(client, username)
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
