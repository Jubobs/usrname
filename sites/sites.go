package sites

type Site interface {
	Name() string
	Home() string
	Validate(username string) []Violation
	IsAvailable(client Client) func(string) (bool, error)
}

func All() []Site {
	return []Site{
		// Facebook(),
		// GitHub(),
		// Instagram(),
		Twitter(),
	}
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
