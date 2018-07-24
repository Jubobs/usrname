package github

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"unicode"

	"github.com/jubobs/usrname"
	"github.com/jubobs/usrname/internal"
)

type github struct {
	name             string
	scheme           string
	host             string
	illegalPrefix    string
	illegalSuffix    string
	illegalSubstring string
	whitelist        *unicode.RangeTable
	minLength        int
	maxLength        int
}

var githubImpl = github{
	name:             "GitHub",
	scheme:           "https",
	host:             "github.com",
	illegalPrefix:    "-",
	illegalSuffix:    "-",
	illegalSubstring: "--",
	whitelist: &unicode.RangeTable{
		R16: []unicode.Range16{
			{'-', '-', 1},
			{'0', '9', 1},
			{'A', 'Z', 1},
			{'a', 'z', 1},
		},
	},
	minLength: 1,
	maxLength: 39,
}

func New() usrname.Checker {
	return &githubImpl
}

func (s *github) Name() string {
	return s.name
}

func (*github) Link(username string) string {
	u := url.URL{
		Scheme: githubImpl.scheme,
		Host:   githubImpl.host,
		Path:   username,
	}
	return u.String()
}

func (*github) IllegalPattern() *regexp.Regexp {
	return nil
}

func (v *github) Whitelist() *unicode.RangeTable {
	return v.whitelist
}

// See https://help.github.com/en/managing-your-account/github-username-rules
func (v *github) Validate(username string) []usrname.Violation {
	return internal.CheckAll(
		username,
		internal.CheckLongerThan(v.minLength),
		internal.CheckOnlyContains(v.whitelist),
		internal.CheckIllegalPrefix(v.illegalPrefix),
		internal.CheckIllegalSubstring(v.illegalSubstring),
		internal.CheckIllegalSuffix(v.illegalSuffix),
		internal.CheckShorterThan(v.maxLength),
	)
}

func (c *github) Check(client usrname.Client) func(string) usrname.Result {
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
	l := githubImpl.Link(username)
	req, err := http.NewRequest("HEAD", l, nil)
	if err != nil {
		panic(err) // should never happen
	}
	return req
}
