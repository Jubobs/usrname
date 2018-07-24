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

func (c *medium) Check(client usrname.Client) func(string) usrname.Result {
	return func(username string) (r usrname.Result) {
		r.Username = username
		r.Checker = c

		if vv := c.Validate(username); len(vv) != 0 {
			r.Status = usrname.Invalid
			const templ = "%q is invalid on %s"
			r.Message = fmt.Sprintf(templ, username, c.Name())
			return
		}

		req := request(username)
		res, err := client.Do(req)
		if err != nil {
			r.Status = usrname.UnknownStatus
			if internal.IsTimeout(err) {
				r.Message = fmt.Sprintf("%s timed out", c.Name())
			} else {
				r.Message = "Something went wrong"
			}
			return
		}

		switch res.StatusCode {
		case http.StatusOK:
			r.Status = usrname.Unavailable
		case http.StatusNotFound:
			r.Status = usrname.Available
		default:
			r.Status = usrname.UnknownStatus
			r.Message = "Something went wrong"
		}
		return
	}
}

func request(username string) *http.Request {
	l := mediumImpl.Link(username)
	req, err := http.NewRequest("HEAD", l, nil)
	if err != nil {
		panic(err)
	}
	return req
}
