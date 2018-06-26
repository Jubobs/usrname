package sites

import (
	"fmt"
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

func Instagram() NameChecker {
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
	fmt.Println(u)
	statusCode, err := client.GetStatusCode(u)
	if err != nil {
		return err
	}
	switch statusCode {
	case 200:
		return &unavailableUsernameError{
			Namer:    Instagram(),
			username: username,
		}
	case 404:
		return nil
	default:
		return &unexpectedError{err}
	}
}
