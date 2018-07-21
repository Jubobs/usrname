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

func (t *reddit) Name() string {
	return t.name
}

func (t *reddit) Link(username string) string {
	u := url.URL{
		Scheme: redditImpl.scheme,
		Host:   redditImpl.host,
		Path:   "/user/" + username,
	}
	return u.String()
}

func (t *reddit) IllegalPattern() *regexp.Regexp {
	return nil
}

func (t *reddit) Whitelist() *unicode.RangeTable {
	return t.whitelist
}

// See https://help.reddit.com/en/managing-your-account/reddit-username-rules
func (t *reddit) Validate(username string) []usrname.Violation {
	return internal.CheckAll(
		username,
		internal.CheckLongerThan(t.minLength),
		internal.CheckOnlyContains(t.whitelist),
		internal.CheckShorterThan(t.maxLength),
	)
}

func (t *reddit) Check(client usrname.Client) func(string) usrname.Result {
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
		req.Header.Add("User-Agent", "Mozilla/5.0") // to avoid rate limiting
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
