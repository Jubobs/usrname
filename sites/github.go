package sites

import (
	"net/url"
)

var GitHub NameChecker = &github{
	urlFrom: func(username string) url.URL {
		return url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   username,
		}
	},
}

type github struct {
	urlFrom func(string) url.URL
}

func (*github) Name() string {
	return "GitHub"
}

func (*github) Validate(username string) error {
	return nil
}

func (g *github) Check(client Client, username string) error {
	u := g.urlFrom(username)
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
