package twitter

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"unicode"

	"github.com/jubobs/usrname"
	"github.com/jubobs/usrname/internal"
)

type twitter struct {
	name           string
	scheme         string
	host           string
	suspended      string
	illegalPattern *regexp.Regexp
	whitelist      *unicode.RangeTable
	minLength      int
	maxLength      int
}

var twitterImpl = twitter{
	name:           "Twitter",
	scheme:         "https",
	host:           "twitter.com",
	suspended:      "https://twitter.com/account/suspended",
	illegalPattern: regexp.MustCompile("(?i)twitter"),
	whitelist: &unicode.RangeTable{
		R16: []unicode.Range16{
			{'0', '9', 1},
			{'A', 'Z', 1},
			{'_', '_', 1},
			{'a', 'z', 1},
		},
	},
	minLength: 1,
	maxLength: 15,
}

func New() usrname.Checker {
	return &twitterImpl
}

func (t *twitter) Name() string {
	return t.name
}

func (t *twitter) Link(username string) string {
	u := url.URL{
		Scheme: twitterImpl.scheme,
		Host:   twitterImpl.host,
		Path:   username,
	}
	return u.String()
}

func (t *twitter) IllegalPattern() *regexp.Regexp {
	return t.illegalPattern
}

func (t *twitter) Whitelist() *unicode.RangeTable {
	return t.whitelist
}

// See https://help.twitter.com/en/managing-your-account/twitter-username-rules
func (t *twitter) Validate(username string) []usrname.Violation {
	return internal.CheckAll(
		username,
		internal.CheckLongerThan(t.minLength),
		internal.CheckOnlyContains(t.whitelist),
		internal.CheckNotMatches(t.illegalPattern),
		internal.CheckShorterThan(t.maxLength),
	)
}

func (t *twitter) Check(client usrname.Client) func(string) usrname.Result {
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
			r.Message = "account unavailable"
		case http.StatusFound:
			if loc := res.Header["location"]; len(loc) == 1 && loc[0] == t.suspended {
				r.Status = usrname.Unavailable
				r.Message = "account suspended"
			} else {
				r.Status = usrname.UnknownStatus
				r.Message = "303 Found, but unexpected 'location'"
			}
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
	u := twitterImpl.Link(username)
	req, err := http.NewRequest("HEAD", u, nil)
	if err != nil {
		panic(err) // should never happen
	}
	return req
}
