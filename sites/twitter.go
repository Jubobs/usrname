package sites

import (
	"net/http"
	"net/url"
	"unicode"
)

type Twitter struct {
	name             string
	home             string
	scheme           string
	host             string
	illegalSubstring string
	whitelist        *unicode.RangeTable
	minLength        int
	maxLength        int
}

var twitterImpl = Twitter{
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

func twitterRequest(username string) (*http.Request, error) {
	u := url.URL{
		Scheme: twitterImpl.scheme,
		Host:   twitterImpl.host,
		Path:   username,
	}
	return http.NewRequest("HEAD", u.String(), nil)
}

func NewTwitter() *Twitter {
	return &twitterImpl
}

func (t *Twitter) Name() string {
	return t.name
}

func (t *Twitter) Home() string {
	return t.home
}

// See https://help.twitter.com/en/managing-your-account/twitter-username-rules
func (t *Twitter) CheckValid(username string) []Violation {
	return checkAll(
		username,
		checkLongerThan(t.minLength),
		checkOnlyContains(t.whitelist),
		checkNotContains(t.illegalSubstring),
		checkShorterThan(t.maxLength),
	)
}

func (t *Twitter) CheckAvailable(client Client) func(string) (bool, error) {
	return func(username string) (bool, error) {
		req, err := twitterRequest(username)
		statusCode, err := client.Send(req)
		if err != nil {
			return false, &networkError{err}
		}
		switch statusCode {
		case http.StatusOK:
			return false, nil
		case http.StatusNotFound:
			return true, nil
		default:
			return false, &unexpectedStatusCodeError{statusCode}
		}
	}
}
