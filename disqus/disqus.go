package disqus

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"unicode"

	"github.com/jubobs/usrname"
	"github.com/jubobs/usrname/internal"
)

type disqus struct {
	name          string
	scheme        string
	host          string
	illegalPrefix string
	illegalSuffix string
	whitelist     *unicode.RangeTable
	minLength     int
	maxLength     int
}

var disqusImpl = disqus{
	name:          "Disqus",
	scheme:        "https",
	host:          "disqus.com",
	illegalPrefix: "_",
	illegalSuffix: "_",
	whitelist: &unicode.RangeTable{
		R16: []unicode.Range16{
			{'0', '9', 1},
			{'A', 'Z', 1},
			{'_', '_', 1},
			{'a', 'z', 1},
		},
	},
	minLength: 2,
	maxLength: 30,
}

func New() usrname.Checker {
	return &disqusImpl
}

func (t *disqus) Name() string {
	return t.name
}

func (t *disqus) Link(username string) string {
	u := url.URL{
		Scheme: disqusImpl.scheme,
		Host:   disqusImpl.host,
		Path:   "/by/" + username,
	}
	return u.String() + "/" // important to avoid redirects
}

func (t *disqus) IllegalPattern() *regexp.Regexp {
	return nil
}

func (t *disqus) Whitelist() *unicode.RangeTable {
	return t.whitelist
}

// See https://help.disqus.com/en/managing-your-account/disqus-username-rules
func (t *disqus) Validate(username string) []usrname.Violation {
	return internal.CheckAll(
		username,
		internal.CheckLongerThan(t.minLength),
		internal.CheckOnlyContains(t.whitelist),
		internal.CheckIllegalPrefix(t.illegalPrefix),
		internal.CheckIllegalSuffix(t.illegalSuffix),
		internal.CheckShorterThan(t.maxLength),
	)
}

func (t *disqus) Check(client usrname.Client) func(string) usrname.Result {
	return func(username string) (r usrname.Result) {
		r.Username = username
		r.Checker = t

		if vv := t.Validate(username); len(vv) != 0 {
			r.Status = usrname.Invalid
			const templ = "%q is invalid on %s"
			r.Message = fmt.Sprintf(templ, username, t.Name())
			return
		}

		req := request(username)
		res, err := client.Do(req)
		if err != nil {
			r.Status = usrname.UnknownStatus
			if internal.IsTimeout(err) {
				r.Message = fmt.Sprintf("%s timed out", t.Name())
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
			r.Message = fmt.Sprintf("unexpected status code %d", res.StatusCode)
		}
		return
	}
}

func request(username string) *http.Request {
	l := disqusImpl.Link(username)
	req, err := http.NewRequest("HEAD", l, nil)
	if err != nil {
		panic(err) // should never happen
	}
	return req
}
