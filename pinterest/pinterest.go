package pinterest

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"unicode"

	"github.com/jubobs/usrname"
	"github.com/jubobs/usrname/internal"
)

type pinterest struct {
	name          string
	scheme        string
	host          string
	illegalPrefix string
	whitelist     *unicode.RangeTable
	minLength     int
	maxLength     int
}

var pinterestImpl = pinterest{
	name:          "Pinterest",
	scheme:        "https",
	host:          "www.pinterest.com",
	illegalPrefix: "_",
	whitelist: &unicode.RangeTable{
		R16: []unicode.Range16{
			{'0', '9', 1},
			{'A', 'Z', 1},
			{'_', '_', 1},
			{'a', 'z', 1},
		},
	},
	minLength: 3,
	maxLength: 30,
}

func New() usrname.Checker {
	return &pinterestImpl
}

func (t *pinterest) Name() string {
	return t.name
}

func (t *pinterest) Link(username string) string {
	u := url.URL{
		Scheme: pinterestImpl.scheme,
		Host:   pinterestImpl.host,
		Path:   username,
	}
	return u.String() + "/" // important to avoid redirects
}

func (t *pinterest) IllegalPattern() *regexp.Regexp {
	return nil
}

func (t *pinterest) Whitelist() *unicode.RangeTable {
	return t.whitelist
}

// See https://help.pinterest.com/en/managing-your-account/pinterest-username-rules
func (t *pinterest) Validate(username string) []usrname.Violation {
	return internal.CheckAll(
		username,
		internal.CheckLongerThan(t.minLength),
		internal.CheckOnlyContains(t.whitelist),
		internal.CheckIllegalPrefix(t.illegalPrefix),
		internal.CheckShorterThan(t.maxLength),
	)
}

func (c *pinterest) Check(client usrname.Client) func(string) usrname.Result {
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
	l := pinterestImpl.Link(username)
	req, err := http.NewRequest("HEAD", l, nil)
	if err != nil {
		panic(err) // should never happen
	}
	return req
}
