package sites

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

type twitter struct {
	name   string
	home   string
	scheme string
	host   string
}

var twitterImpl = twitter{
	name:   "Twitter",
	home:   "https://twitter.com",
	scheme: "https",
	host:   "twitter.com",
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
func (*twitter) CheckValid(username string) []Violation {
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
