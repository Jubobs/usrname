package sites_test

import (
	"errors"
	"github.com/jubobs/username-checker/sites"
	"net/http"
	"testing"
)

var checker = sites.Twitter()

func TestTwitterName(t *testing.T) {
	const expected = "Twitter"
	actual := checker.Name()
	if actual != expected {
		template := "sites.Twitter().Name() == %q, want %q"
		t.Errorf(template, actual, expected)
	}
}

func TestTwitterValidate(t *testing.T) {
	cases := []struct {
		username string
		valid    bool
	}{
		{"", false},
		{"0", true},
		{"exotic^chars", false},
		{"underscores_ok", true},
		{"twitter_no_ok", false},
		{"not_ok_TwitteR", false},
		{"admin_fine", true},
		{"longerthan15char", false},
	}
	const template = "(IsInvalidUsernameError(Twitter().Validate(%q)) == %t, want %t"
	for _, c := range cases {
		err := checker.Validate(c.username)
		if sites.IsInvalidUsernameError(err) == c.valid {
			t.Errorf(template, c.username, err == nil, c.valid)
		}
	}
}

func TestCheckNotFound(t *testing.T) {
	client := mockClientHead(http.StatusNotFound, nil)
	const dummyUsername = "dummy"

	var expected error = nil
	actual := checker.Check(client, dummyUsername)
	if actual != nil {
		t.Errorf("Twitter().Check() == %v, want %v", actual, expected)
	}
}

func TestCheckOK(t *testing.T) {
	client := mockClientHead(http.StatusOK, nil)
	const dummyUsername = "dummy"

	const expected = true
	actual := sites.IsUnavailableUsernameError(checker.Check(client, dummyUsername))
	if actual != expected {
		const template = "(IsUnavailableUsernameError(Twitter().Validate(%q)) == %t, want %t"
		t.Errorf(template, dummyUsername, actual, expected)
	}
}

func TestCheckOther(t *testing.T) {
	const statusCode = 999 // anything other than 200 and 404
	err := errors.New("Oh no!")
	client := mockClientHead(statusCode, err)
	const dummyUsername = "dummy"
	t.Log(checker.Check(client, dummyUsername))
	const expected = true
	actual := sites.IsUnexpectedError(checker.Check(client, dummyUsername))
	if actual != expected {
		const template = "(IsUnexpectedError(Twitter().Validate(%q)) == %t, want %t"
		t.Errorf(template, dummyUsername, actual, expected)
	}
}
