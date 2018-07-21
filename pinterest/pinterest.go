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

func (t *pinterest) Check(client usrname.Client) func(string) usrname.Result {
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
