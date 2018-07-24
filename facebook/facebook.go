package facebook

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/jubobs/usrname"
	"github.com/jubobs/usrname/internal"
)

type facebook struct {
	name      string
	scheme    string
	host      string
	whitelist *unicode.RangeTable
	minLength int
	maxLength int
}

var facebookImpl = facebook{
	name:   "facebook",
	scheme: "https",
	host:   "www.facebook.com",
	whitelist: &unicode.RangeTable{
		R16: []unicode.Range16{
			{'.', '.', 1},
			{'0', '9', 1},
			{'A', 'Z', 1},
			{'a', 'z', 1},
		},
	},
	minLength: 5,
	maxLength: 50,
}

func New() usrname.Checker {
	return &facebookImpl
}

func (s *facebook) Name() string {
	return s.name
}

func (s *facebook) Link(username string) string {
	u := url.URL{
		Scheme: s.scheme,
		Host:   s.host,
		Path:   username,
	}
	return u.String()
}

func (*facebook) IllegalPattern() *regexp.Regexp {
	return nil
}

func (v *facebook) Whitelist() *unicode.RangeTable {
	return v.whitelist
}

// See https://help.facebook.com/en/managing-your-account/facebook-username-rules
func (v *facebook) Validate(username string) []usrname.Violation {
	return internal.CheckAll(
		username,
		internal.CheckLongerThan(v.minLength),
		internal.CheckOnlyContains(v.whitelist),
		internal.CheckShorterThan(v.maxLength),
	)
}

func (c *facebook) Check(client usrname.Client) func(string) usrname.Result {
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
			r.Message = "account unavailable"
		case http.StatusNotFound:
			r.Status = usrname.Available
		case http.StatusFound:
			if loc := res.Header["location"]; len(loc) == 1 && checkRedirect(username, loc[0]) {
				r.Status = usrname.Unavailable
				r.Message = "account suspended"
			} else {
				r.Status = usrname.UnknownStatus
				r.Message = "302 Found, but unexpected 'location'"
			}
		default:
			r.Status = usrname.UnknownStatus
			r.Message = fmt.Sprintf("unsupported status code %d", res.StatusCode)
		}
		return
	}
}

func request(username string) *http.Request {
	l := facebookImpl.Link(username)
	req, err := http.NewRequest("HEAD", l, nil)
	if err != nil {
		panic(err) // should never happen
	}
	return req
}

func checkRedirect(username string, location string) bool {
	root := facebookImpl.Link("") + "/"
	ss := strings.SplitAfterN(location, root, 2)
	return len(ss) == 2 && strings.EqualFold(normalize(ss[1]), normalize(username))
}

func normalize(s string) string {
	return strings.Replace(s, ".", "", -1)
}
