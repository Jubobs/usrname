package sites

import (
	"net/http"
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"
)

type twitter struct {
	name      string
	home      string
	scheme    string
	host      string
	whitelist *unicode.RangeTable
}

var twitterImpl = twitter{
	name:   "Twitter",
	home:   "https://twitter.com",
	scheme: "https",
	host:   "twitter.com",
	whitelist: &unicode.RangeTable{
		R16: []unicode.Range16{
			{'0', '9', 1},
			{'A', 'Z', 1},
			{'_', '_', 1},
			{'a', 'z', 1},
		},
	},
}

func twitterRequest(username string) (*http.Request, error) {
	u := url.URL{
		Scheme: twitterImpl.scheme,
		Host:   twitterImpl.host,
		Path:   username,
	}
	return http.NewRequest("HEAD", u.String(), nil)
}

const (
	minLength        = 1
	maxLength        = 15
	illegalSubstring = "twitter"
)

func Twitter() Site {
	return &twitterImpl
}

func (t *twitter) Name() string {
	return t.name
}

func (t *twitter) Home() string {
	return t.home
}

// See https://help.twitter.com/en/managing-your-account/twitter-username-rules
func (t *twitter) CheckValid(username string) []Violation {
	runeCount := utf8.RuneCountInString(username)
	violations := []Violation{}
	if runeCount < minLength {
		v := TooShort{
			Min:    minLength,
			Actual: runeCount,
		}
		violations = append(violations, &v)
	}

	var inds []int
	for i, r := range username {
		if !unicode.In(r, t.whitelist) {
			inds = append(inds, i)
		}
	}
	if len(inds) != 0 {
		v := IllegalChars{
			At:        inds,
			Whitelist: t.whitelist,
		}
		violations = append(violations, &v)
	}

	if i := strings.Index(strings.ToLower(username), illegalSubstring); i != -1 {
		v := IllegalSubstring{
			Sub: illegalSubstring,
			At:  i,
		}
		violations = append(violations, &v)
	}
	if maxLength < runeCount {
		v := TooLong{
			Max:    maxLength,
			Actual: runeCount,
		}
		violations = append(violations, &v)
	}
	return violations
}

func (t *twitter) CheckAvailable(client Client) func(string) (bool, error) {
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
