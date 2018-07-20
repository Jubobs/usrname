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

func (t *instagram) Name() string {
	return t.name
}

func (t *instagram) Link(username string) string {
	u := url.URL{
		Scheme: instagramImpl.scheme,
		Host:   instagramImpl.host,
		Path:   username,
	}
	return u.String() + "/" // important to avoid redirects
}

func (t *instagram) IllegalPattern() *regexp.Regexp {
	return nil
}

func (t *instagram) Whitelist() *unicode.RangeTable {
	return t.whitelist
}

// See https://help.instagram.com/en/managing-your-account/instagram-username-rules
func (t *instagram) Validate(username string) []usrname.Violation {
	return internal.CheckAll(
		username,
		internal.CheckLongerThan(t.minLength),
		internal.CheckOnlyContains(t.whitelist),
		internal.CheckIllegalPrefix(t.illegalPrefix),
		internal.CheckIllegalSubstring(t.illegalSubstring),
		internal.CheckIllegalSuffix(t.illegalSuffix),
		internal.CheckShorterThan(t.maxLength),
	)
}

func (t *instagram) Check(client usrname.Client) func(string) usrname.Result {
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
		req, err := http.NewRequest("GET", u, nil)
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
