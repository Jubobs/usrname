package medium

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"unicode"

	"github.com/jubobs/usrname"
	"github.com/jubobs/usrname/internal"
)

type medium struct {
	name      string
	scheme    string
	host      string
	whitelist *unicode.RangeTable
	minLength int
	maxLength int
}

var mediumImpl = medium{
	name:   "Medium",
	scheme: "https",
	host:   "medium.com",
	whitelist: &unicode.RangeTable{
		R16: []unicode.Range16{
			{'.', '.', 1},
			{'0', '9', 1},
			{'A', 'Z', 1},
			{'_', '_', 1},
			{'a', 'z', 1},
		},
	},
	minLength: 1,
	maxLength: 16,
}

func New() usrname.Checker {
	return &mediumImpl
}

func (t *medium) Name() string {
	return t.name
}

func (t *medium) Link(username string) string {
	u := url.URL{
		Scheme: mediumImpl.scheme,
		Host:   mediumImpl.host,
		Path:   "@" + username,
	}
	return u.String()
}

func (t *medium) IllegalPattern() *regexp.Regexp {
	return nil
}

func (t *medium) Whitelist() *unicode.RangeTable {
	return t.whitelist
}

// See https://help.medium.com/en/managing-your-account/medium-username-rules
func (t *medium) Validate(username string) []usrname.Violation {
	return internal.CheckAll(
		username,
		internal.CheckLongerThan(t.minLength),
		internal.CheckOnlyContains(t.whitelist),
		internal.CheckShorterThan(t.maxLength),
	)
}

func (t *medium) Check(client usrname.Client) func(string) usrname.Result {
	return func(username string) (res usrname.Result) {
		res.Username = username
		res.Checker = t

		if vv := t.Validate(username); len(vv) != 0 {
			res.Status = usrname.Invalid
			const templ = "%q is invalid on %s"
			res.Message = fmt.Sprintf(templ, username, t.Name())
			return
		}

		u := t.Link(username)
		req, err := http.NewRequest("HEAD", u, nil)
		statusCode, err := client.Send(req)
		if err != nil {
			res.Status = usrname.UnknownStatus
			type timeout interface {
				Timeout() bool
			}
			if err, ok := err.(timeout); ok && err.Timeout() {
				res.Message = fmt.Sprintf("%s timed out", t.Name())
			} else {
				res.Message = "Something went wrong"
			}
		}
		switch statusCode {
		case http.StatusOK:
			res.Status = usrname.Unavailable
		case http.StatusNotFound:
			res.Status = usrname.Available
		default:
			res.Status = usrname.UnknownStatus
			res.Message = "Something went wrong"
		}
		return
	}
}
