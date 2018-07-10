package twitter

import (
	"net/http"
	"net/url"
	"regexp"
	"unicode"

	"github.com/jubobs/usrname"
	"github.com/jubobs/usrname/internal"
)

type twitter struct {
	name           string
	home           string
	scheme         string
	host           string
	illegalPattern *regexp.Regexp
	whitelist      *unicode.RangeTable
	minLength      int
	maxLength      int
}

var twitterImpl = twitter{
	name:           "Twitter",
	home:           "https://twitter.com",
	scheme:         "https",
	host:           "twitter.com",
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

func (t *twitter) Home() string {
	return t.home
}

func (t *twitter) ProfilePage(username string) string {
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

func (t *twitter) Check(client usrname.Client) func(string) (bool, error) {
	return func(username string) (bool, error) {
		u := t.ProfilePage(username)
		req, err := http.NewRequest("HEAD", u, nil)
		statusCode, err := client.Send(req)
		if err != nil {
			return false, &usrname.NetworkError{Cause: err}
		}
		switch statusCode {
		case http.StatusOK:
			return false, nil
		case http.StatusNotFound:
			return true, nil
		default:
			return false, &usrname.UnexpectedStatusCodeError{StatusCode: statusCode}
		}
	}
}