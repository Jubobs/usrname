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

const (
	minLength        = 1
	maxLength        = 15
	illegalSubstring = "twitter"
	expectedPattern  = "^[A-Za-z0-9_]*$"
)

var expectedRegexp = regexp.MustCompile(expectedPattern)

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
func (*twitter) Validate(username string) []Violation {
	runeCount := utf8.RuneCountInString(username)
	violations := []Violation{}
	if runeCount < minLength {
		v := TooShort{
			Min:    minLength,
			Actual: runeCount,
		}
		violations = append(violations, &v)
	}
	if !expectedRegexp.MatchString(username) {
		v := IllegalChars{}
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

func (t *twitter) IsAvailable(client Client) func(string) (bool, error) {
	return func(username string) (bool, error) {
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
}
