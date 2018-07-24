package reddit

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"unicode"

	"github.com/jubobs/usrname"
	"github.com/jubobs/usrname/internal"
)

type reddit struct {
	name      string
	scheme    string
	host      string
	whitelist *unicode.RangeTable
	minLength int
	maxLength int
}

var redditImpl = reddit{
	name:   "reddit",
	scheme: "https",
	host:   "www.reddit.com",
	whitelist: &unicode.RangeTable{
		R16: []unicode.Range16{
			{'0', '9', 1},
			{'A', 'Z', 1},
			{'_', '_', 1},
			{'a', 'z', 1},
		},
	},
	minLength: 3,
	maxLength: 20,
}

func New() usrname.Checker {
	return &redditImpl
}

func (s *reddit) Name() string {
	return s.name
}

func (*reddit) Link(username string) string {
	u := url.URL{
		Scheme: redditImpl.scheme,
		Host:   redditImpl.host,
		Path:   "/user/" + username,
	}
	return u.String()
}

func (*reddit) IllegalPattern() *regexp.Regexp {
	return nil
}

func (v *reddit) Whitelist() *unicode.RangeTable {
	return v.whitelist
}

// See https://help.reddit.com/en/managing-your-account/reddit-username-rules
func (v *reddit) Validate(username string) []usrname.Violation {
	return internal.CheckAll(
		username,
		internal.CheckLongerThan(v.minLength),
		internal.CheckOnlyContains(v.whitelist),
		internal.CheckShorterThan(v.maxLength),
	)
}

func (c *reddit) Check(client usrname.Client) func(string) usrname.Result {
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
			r.Message = fmt.Sprintf("unsupported status code %d", res.StatusCode)
		}
		return
	}
}

func request(username string) *http.Request {
	l := redditImpl.Link(username)
	req, err := http.NewRequest("HEAD", l, nil)
	if err != nil {
		panic(err) // should never happen
	}
	req.Header.Add("User-Agent", "Mozilla/5.0") // to avoid rate limiting
	return req
}
