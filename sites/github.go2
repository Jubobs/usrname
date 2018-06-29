package sites

import (
	"net/http"
	"net/url"
)

type github struct {
	urlFrom func(string) url.URL
}

var githubImpl = github{
	func(username string) url.URL {
		return url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   username,
		}
	},
}

func GitHub() ValidNameChecker {
	return &githubImpl
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
	case http.StatusOK:
		return &unavailableUsernameError{
			Namer:    GitHub(),
			username: username,
		}
	case http.StatusNotFound:
		return nil
	default:
		return &unexpectedError{err}
	}
}
