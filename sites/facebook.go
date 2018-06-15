package sites

import (
	"net/url"
)

var Facebook NameChecker = &facebook{
	urlFrom: func(username string) url.URL {
		return url.URL{
			Scheme: "https",
			Host:   "www.facebook.com",
			Path:   username,
		}
	},
}

type facebook struct {
	urlFrom func(string) url.URL
}

func (*facebook) Name() string {
	return "Facebook"
}

func (*facebook) Validate(username string) error {
	return nil
}

func (f *facebook) Check(client Client, username string) error {
	u := f.urlFrom(username)
	statusCode, err := client.HeadStatusCode(u)
	if err != nil {
		return err
	}
	switch statusCode {
	case 200:
		return ErrUnavailableUsername
	case 404:
		return nil
	default:
		return err
	}
}
