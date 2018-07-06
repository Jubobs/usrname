package twitter_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/jubobs/username-checker/mock"
	"github.com/jubobs/username-checker/sites"
	"github.com/jubobs/username-checker/sites/twitter"
)

var s = twitter.New()

func TestName(t *testing.T) {
	const expected = "Twitter"
	actual := s.Name()
	if actual != expected {
		template := "twitter.New().Name() == %q, want %q"
		t.Errorf(template, actual, expected)
	}
}

func TestHome(t *testing.T) {
	const expected = "https://twitter.com"
	actual := s.Home()
	if actual != expected {
		template := "twitter.New().Home() == %q, want %q"
		t.Errorf(template, actual, expected)
	}
}

func TestCheckValid(t *testing.T) {
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
	const template = "(len(twitter.New().CheckValid(%q))) is %d, but expected %d"
	for _, c := range cases {
		if vs := s.CheckValid(c.username); len(vs) != c.noOfViolations {
			t.Errorf(template, c.username, len(vs), c.noOfViolations)
		}
	}
}

func TestCheckNotFound(t *testing.T) {
	// Given
	client := mock.Client(http.StatusNotFound, nil)
	const dummyUsername = "dummy"

	// When
	available, err := s.CheckAvailable(client)(dummyUsername)

	// Then
	if !(err == nil && available) {
		const template = "twitter.New().CheckAvailable(%q) == (%t, %v), but expected (true, <nil>)"
		t.Errorf(template, dummyUsername, available, err)
	}
}

func TestCheckOk(t *testing.T) {
	// Given
	client := mock.Client(http.StatusOK, nil)
	const dummyUsername = "dummy"

	// When
	available, err := s.CheckAvailable(client)(dummyUsername)

	// Then
	if err != nil || available {
		const template = "twitter.New().CheckAvailable(%q) == (%t, %v), but expected (false, <nil>)"
		t.Errorf(template, dummyUsername, available, err)
	}
}

func TestCheckOther(t *testing.T) {
	// Given
	const statusCode = 999 // anything other than 200 and 404
	client := mock.Client(statusCode, nil)
	const dummyUsername = "dummy"

	// When
	_, err := s.CheckAvailable(client)(dummyUsername) // irrelevant bool

	// Then
	if actual, ok := err.(*sites.UnexpectedStatusCodeError); !ok {
		const template = "got %v, but want %v"
		expected := &sites.UnexpectedStatusCodeError{statusCode}
		t.Errorf(template, actual, expected)
	}
}

func TestCheckNetworkError(t *testing.T) {
	// Given
	someError := errors.New("Oh no!")
	client := mock.Client(0, someError)
	const dummyUsername = "dummy"

	// When
	_, err := s.CheckAvailable(client)(dummyUsername) // irrelevant bool

	// Then
	if actual, ok := err.(*sites.NetworkError); !ok {
		const template = "got %v, but want %v"
		expected := &sites.NetworkError{someError}
		t.Errorf(template, actual, expected)
	}
}
