package sites_test

import (
	"errors"
	"github.com/jubobs/username-checker/sites"
	"net/http"
	"testing"
)

var s = sites.Twitter()

func TestTwitterName(t *testing.T) {
	const expected = "Twitter"
	actual := s.Name()
	if actual != expected {
		template := "sites.Twitter().Name() == %q, want %q"
		t.Errorf(template, actual, expected)
	}
}

func TestTwitterHome(t *testing.T) {
	const expected = "https://twitter.com"
	actual := s.Home()
	if actual != expected {
		template := "sites.Twitter().Home() == %q, want %q"
		t.Errorf(template, actual, expected)
	}
}

func TestTwitterValidate(t *testing.T) {
	cases := []struct {
		username       string
		noOfViolations int // TODO: refine when I introduce Violation type
	}{
		{"", 1},
		{"0", 0},
		{"exotic^chars", 1},
		{"underscores_ok", 0},
		{"twitter_no_ok", 1},
		{"not_ok_TwitteR", 1},
		{"admin_fine", 0},
		{"longerthan15char", 1},
		{"exotic^chars_and_too_long", 2},
	}
	const template = "(len(Twitter().Validate(%q))) is %d, but expected %d"
	for _, c := range cases {
		if vs := s.Validate(c.username); len(vs) != c.noOfViolations {
			t.Errorf(template, c.username, len(vs), c.noOfViolations)
		}
	}
}

func TestCheckNotFound(t *testing.T) {
	// Given
	client := mockClientHead(http.StatusNotFound, nil)
	const dummyUsername = "dummy"

	// When
	available, err := s.IsAvailable(client)(dummyUsername)

	// Then
	if !(err == nil && available) {
		const template = "Twitter().IsAvailable(%q) == (%t, %v), but expected (true, <nil>)"
		t.Errorf(template, dummyUsername, available, err)
	}
}

func TestCheckOk(t *testing.T) {
	// Given
	client := mockClientHead(http.StatusOK, nil)
	const dummyUsername = "dummy"

	// When
	available, err := s.IsAvailable(client)(dummyUsername)

	// Then
	if err != nil || available {
		const template = "Twitter().IsAvailable(%q) == (%t, %v), but expected (false, <nil>)"
		t.Errorf(template, dummyUsername, available, err)
	}
}

func TestCheckOther(t *testing.T) {
	// Given
	const statusCode = 999 // anything other than 200 and 404
	client := mockClientHead(statusCode, nil)
	const dummyUsername = "dummy"

	// When
	_, err := s.IsAvailable(client)(dummyUsername) // irrelevant bool

	// Then
	if !sites.IsUnexpectedStatusCodeError(err) {
		const template = "got %v, but want an unexpected-status-code error"
		t.Errorf(template, err)
	}
}

func TestCheckNetworkError(t *testing.T) {
	// Given
	someError := errors.New("Oh no!")
	client := mockClientHead(0, someError)
	const dummyUsername = "dummy"

	// When
	_, err := s.IsAvailable(client)(dummyUsername) // irrelevant bool

	// Then
	if !sites.IsNetworkError(err) {
		const template = "got %v, but want network error"
		t.Errorf(template, err)
	}
}
