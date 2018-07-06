package twitter

import (
	"net/http"
	"net/url"
	"unicode"

	"github.com/jubobs/username-checker/sites"
)

type twitter struct {
	name             string
	home             string
	scheme           string
	host             string
	illegalSubstring string
	whitelist        *unicode.RangeTable
	minLength        int
	maxLength        int
}

func New() sites.Site {
	return &twitterImpl
}

func (t *twitter) Name() string {
	return t.name
}

func (t *twitter) Home() string {
	return t.home
}

func (t *twitter) WhitelistChars() *unicode.RangeTable {
	return t.whitelist
}

// See https://help.twitter.com/en/managing-your-account/twitter-username-rules
func (t *twitter) CheckValid(username string) []sites.Violation {
	return sites.CheckAll(
		username,
		sites.CheckLongerThan(t.minLength),
		sites.CheckOnlyContains(t.whitelist),
		sites.CheckNotContains(t.illegalSubstring),
		sites.CheckShorterThan(t.maxLength),
	)
}

func (t *twitter) CheckAvailable(client sites.Client) func(string) (bool, error) {
	return func(username string) (bool, error) {
		u := url.URL{
			Scheme: twitterImpl.scheme,
			Host:   twitterImpl.host,
			Path:   username,
		}
		req, err := http.NewRequest("HEAD", u.String(), nil)
		statusCode, err := client.Send(req)
		if err != nil {
			return false, &sites.NetworkError{err}
		}
		switch statusCode {
		case http.StatusOK:
			return false, nil
		case http.StatusNotFound:
			return true, nil
		default:
			return false, &sites.UnexpectedStatusCodeError{statusCode}
		}
	}
}

var twitterImpl = twitter{
	name:             "Twitter",
	home:             "https://twitter.com",
	scheme:           "https",
	host:             "twitter.com",
	illegalSubstring: "twitter",
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
