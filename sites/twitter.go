package sites

import (
	"errors" // TODO: remove this dep
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

type twitter struct {
	urlFrom func(string) url.URL
}

var twitterImpl = twitter{
	urlFrom: func(username string) url.URL {
		return url.URL{
			Scheme: "https",
			Host:   "twitter.com",
			Path:   username,
		}
	},
}

var (
	minLength        = 1
	maxLength        = 15
	expectedPattern  = regexp.MustCompile("^[_A-Za-z0-9]+$")
	forbiddenPattern = "twitter"
)

func Twitter() ValidNameChecker {
	return &twitterImpl
}

func (*twitter) Name() string {
	return "Twitter"
}

// See https://help.twitter.com/en/managing-your-account/twitter-username-rules
func (*twitter) Validate(username string) error {
	runeCount := utf8.RuneCountInString(username)
	switch {
	case runeCount < minLength:
		return errors.New("too short") // TODO: append to []Violations
	case !expectedPattern.MatchString(username):
		return errors.New("illegal characters") // TODO: append to []Violations
	case strings.Contains(strings.ToLower(username), forbiddenPattern):
		return errors.New("illegal pattern") // TODO: append to []Violations
	case maxLength < runeCount:
		return errors.New("too long") // TODO: append to []Violations
	default:
		return nil
	}
}

func (t *twitter) Check(client Client, username string) error {
	u := t.urlFrom(username)
	statusCode, err := client.HeadStatusCode(u)
	if err != nil {
		return err
	}
	switch statusCode {
	case 200:
		return &unavailableUsernameError{
			Namer:    Twitter(),
			username: username,
		}
	case 404:
		return nil
	default:
		return &unexpectedError{err}
	}
}
