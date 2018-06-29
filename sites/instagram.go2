package sites

import (
	"net/http"
	"net/url"
)

type instagram struct {
	urlFrom func(string) url.URL
}

// this seems like a dirty trick... improve later
var instagramImpl = instagram{
	urlFrom: func(username string) url.URL {
		u, _ := url.Parse("https://www.instagram.com")
		relative, _ := url.Parse(username + "/")
		return *(u.ResolveReference(relative))
	},
}

func Instagram() ValidNameChecker {
	return &instagramImpl
}

func (*instagram) Name() string {
	return "Instagram"
}

func (*instagram) Validate(username string) error {
	return nil
}

func (t *instagram) Check(client Client, username string) error {
	u := t.urlFrom(username)
	statusCode, err := client.GetStatusCode(u)
	if err != nil {
		return err
	}
	switch statusCode {
	case http.StatusOK:
		return &unavailableUsernameError{
			Namer:    Instagram(),
			username: username,
		}
	case http.StatusNotFound:
		return nil
	default:
		return &unexpectedError{err}
	}
}
