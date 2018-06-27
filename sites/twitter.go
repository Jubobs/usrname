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
	minLength        = 1
	maxLength        = 15
	expectedPattern  = regexp.MustCompile("^[_A-Za-z0-9]+$")
	forbiddenPattern = "twitter"
)

func Twitter() ValidNameChecker {
	return &twitterImpl
}

func (t *twitter) Name() string {
	return t.name
}

func (t *twitter) Home() string {
	return t.home
}

// See https://help.twitter.com/en/managing-your-account/twitter-username-rules
func (*twitter) Validate(username string) error {
	runeCount := utf8.RuneCountInString(username)
	err := &invalidUsernameError{
		Namer:    Twitter(),
		username: username,
	}
	switch {
	case runeCount < minLength:
		return err
	case !expectedPattern.MatchString(username):
		return err
	case strings.Contains(strings.ToLower(username), forbiddenPattern):
		return err
	case maxLength < runeCount:
		return err
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
	case http.StatusOK:
		return &unavailableUsernameError{
			Namer:    Twitter(),
			username: username,
		}
	case http.StatusNotFound:
		return nil
	default:
		return &unexpectedError{err}
	}
}
