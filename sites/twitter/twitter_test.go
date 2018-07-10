package twitter_test

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/fortytw2/leaktest"
	"github.com/jubobs/whocanibe/mock"
	"github.com/jubobs/whocanibe/sites"
	"github.com/jubobs/whocanibe/sites/twitter"
)

var s = twitter.New()

func TestName(t *testing.T) {
	defer leaktest.Check(t)()
	const expected = "Twitter"
	actual := s.Name()
	if actual != expected {
		template := "got %q, want %q"
		t.Errorf(template, actual, expected)
	}
}

func TestHome(t *testing.T) {
	defer leaktest.Check(t)()
	const expected = "https://twitter.com"
	actual := s.Home()
	if actual != expected {
		template := "got %q, want %q"
		t.Errorf(template, actual, expected)
	}
}

func TestValidate(t *testing.T) {
	defer leaktest.Check(t)()
	noViolations := []sites.Violation{}
	cases := []struct {
		username   string
		violations []sites.Violation
	}{
		{
			"",
			[]sites.Violation{
				&sites.TooShort{
					Min:    1,
					Actual: 0,
				},
			},
		}, {
			"0",
			noViolations,
		}, {
			"exotic^chars",
			[]sites.Violation{
				&sites.IllegalChars{
					At:        []int{6},
					Whitelist: s.Whitelist(),
				},
			},
		}, {
			"underscores_ok",
			noViolations,
		}, {
			"twitter_no_ok",
			[]sites.Violation{
				&sites.IllegalSubstring{
					Pattern: s.IllegalPattern().String(),
					At:      []int{0, 7},
				},
			},
		}, {
			"not_ok_TwitteR",
			[]sites.Violation{
				&sites.IllegalSubstring{
					Pattern: s.IllegalPattern().String(),
					At:      []int{7, 14},
				},
			},
		}, {
			"admin_fine",
			noViolations,
		},
		{
			"longerthan15char",
			[]sites.Violation{
				&sites.TooLong{
					Max:    15,
					Actual: 16,
				},
			},
		}, {
			"exotic^chars_and_too_long",
			[]sites.Violation{
				&sites.IllegalChars{
					At:        []int{6},
					Whitelist: s.Whitelist(),
				},
				&sites.TooLong{
					Max:    15,
					Actual: 25,
				},
			},
		},
	}
	const template = "%q, got %#v, want %#v"
	for _, c := range cases {
		if vv := s.Validate(c.username); !reflect.DeepEqual(vv, c.violations) {
			t.Errorf(template, c.username, vv, c.violations)
		}
	}
}

func TestCheckNotFound(t *testing.T) {
	defer leaktest.Check(t)()
	// Given
	client := mock.Client(http.StatusNotFound, nil)
	const dummyUsername = "dummy"

	// When
	available, err := s.Check(client)(dummyUsername)

	// Then
	if !(available && err == nil) {
		const template = "%q, got (%t, %v), want (true, <nil>)"
		t.Errorf(template, dummyUsername, available, err)
	}
}

func TestCheckOk(t *testing.T) {
	defer leaktest.Check(t)()
	// Given
	client := mock.Client(http.StatusOK, nil)
	const dummyUsername = "dummy"

	// When
	available, err := s.Check(client)(dummyUsername)

	// Then
	if !(!available && err == nil) {
		const template = "%q, got (%t, %v), want (false, <nil>)"
		t.Errorf(template, dummyUsername, available, err)
	}
}

func TestCheckOther(t *testing.T) {
	defer leaktest.Check(t)()
	// Given
	const statusCode = 999 // anything other than 200 and 404
	client := mock.Client(statusCode, nil)
	const dummyUsername = "dummy"

	// When
	_, err := s.Check(client)(dummyUsername) // irrelevant bool

	// Then
	if actual, ok := err.(*sites.UnexpectedStatusCodeError); !ok {
		const template = "got %v, want %v"
		expected := &sites.UnexpectedStatusCodeError{StatusCode: statusCode}
		t.Errorf(template, actual, expected)
	}
}

func TestCheckNetworkError(t *testing.T) {
	defer leaktest.Check(t)()
	// Given
	someError := errors.New("Oh no!")
	client := mock.Client(0, someError)
	const dummyUsername = "dummy"

	// When
	_, err := s.Check(client)(dummyUsername) // irrelevant bool

	// Then
	if actual, ok := err.(*sites.NetworkError); !ok {
		const template = "got %v, want %v"
		expected := &sites.NetworkError{Cause: someError}
		t.Errorf(template, actual, expected)
	}
}
