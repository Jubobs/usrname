package instagram

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"unicode"

	"github.com/jubobs/usrname"
	"github.com/jubobs/usrname/internal"
)

type instagram struct {
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

var instagramImpl = instagram{
	name:             "Instagram",
	scheme:           "https",
	host:             "www.instagram.com",
	illegalPrefix:    ".",
	illegalSuffix:    ".",
	illegalSubstring: "..",
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
	maxLength: 30,
}

func New() usrname.Checker {
	return &instagramImpl
}

func (s *instagram) Name() string {
	return s.name
}

func (*instagram) Link(username string) string {
	u := url.URL{
		Scheme: instagramImpl.scheme,
		Host:   instagramImpl.host,
		Path:   username,
	}
	return u.String() + "/" // important to avoid redirects
}

func (*instagram) IllegalPattern() *regexp.Regexp {
	return nil
}

func (v *instagram) Whitelist() *unicode.RangeTable {
	return v.whitelist
}

// See https://help.instagram.com/en/managing-your-account/instagram-username-rules
func (v *instagram) Validate(username string) []usrname.Violation {
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

func (c *instagram) Check(client usrname.Client) func(string) usrname.Result {
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
	u := instagramImpl.Link(username)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		panic(err) // should never happen
	}
	return req
}
