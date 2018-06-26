package sites

import (
	"net/url"
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

func Twitter() NameChecker {
	return &twitterImpl
}

func (*twitter) Name() string {
	return "Twitter"
}

func (*twitter) Validate(username string) error {
	return nil
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
