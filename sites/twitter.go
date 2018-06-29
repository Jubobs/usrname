package sites

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

type twitter struct {
	name    string
	home    string
	urlFrom func(string) url.URL
}

var twitterImpl = twitter{
	name: "Twitter",
	home: "https://twitter.com",
	urlFrom: func(username string) url.URL {
		return url.URL{
			Scheme: "https",
			Host:   "twitter.com",
			Path:   username,
		}
	},
}

var (
	minLength          = 1
	maxLength          = 15
	expectedPattern    = regexp.MustCompile("^[_A-Za-z0-9]+$")
	forbiddenSubstring = "twitter"
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
func (*twitter) Validate(username string) Violations {
	runeCount := utf8.RuneCountInString(username)
	violations := []Violation{"invalid username"} // TODO: tidy this up
	switch {
	case runeCount < minLength:
		return violations
	case !expectedPattern.MatchString(username):
		return violations
	case strings.Contains(strings.ToLower(username), forbiddenSubstring):
		return violations
	case maxLength < runeCount:
		return violations
	default:
		return []Violation{}
	}
}

func (t *twitter) CheckAvailability(client Client, username string) (bool, error) {
	u := t.urlFrom(username)
	statusCode, err := client.HeadStatusCode(u)
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
