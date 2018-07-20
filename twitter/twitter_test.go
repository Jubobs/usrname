package twitter_test

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/fortytw2/leaktest"
	"github.com/jubobs/usrname"
	"github.com/jubobs/usrname/mock"
	"github.com/jubobs/usrname/twitter"
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

func TestLink(t *testing.T) {
	defer leaktest.Check(t)()
	const expected = "https://twitter.com/foobar"
	actual := s.Link("foobar")
	if actual != expected {
		template := "got %q, want %q"
		t.Errorf(template, actual, expected)
	}
}

func TestValidate(t *testing.T) {
	defer leaktest.Check(t)()
	noViolations := []usrname.Violation{}
	cases := []struct {
		label      string
		username   string
		violations []usrname.Violation
	}{
		{
			"empty",
			"",
			[]usrname.Violation{
				&usrname.TooShort{
					Min:    1,
					Actual: 0,
				},
			},
		}, {
			"onechar",
			"0",
			noViolations,
		}, {
			"exoticchars",
			"exotic^chars",
			[]usrname.Violation{
				&usrname.IllegalChars{
					At:        []int{6},
					Whitelist: s.Whitelist(),
				},
			},
		}, {
			"underscores",
			"__init__",
			noViolations,
		}, {
			"twitterlowercase",
			"twitter_no_ok",
			[]usrname.Violation{
				&usrname.IllegalSubstring{
					Pattern: s.IllegalPattern().String(),
					At:      []int{0, 7},
				},
			},
		}, {
			"twittermixedcase",
			"not_ok_TwitteR",
			[]usrname.Violation{
				&usrname.IllegalSubstring{
					Pattern: s.IllegalPattern().String(),
					At:      []int{7, 14},
				},
			},
		}, {
			"admin",
			"admin_fine",
			noViolations,
		},
		{
			"toolong",
			"0123456789012345",
			[]usrname.Violation{
				&usrname.TooLong{
					Max:    15,
					Actual: 16,
				},
			},
		}, {
			"exoticcharstoolong",
			"exotic^chars_and_too_long",
			[]usrname.Violation{
				&usrname.IllegalChars{
					At:        []int{6},
					Whitelist: s.Whitelist(),
				},
				&usrname.TooLong{
					Max:    15,
					Actual: 25,
				},
			},
		},
	}
	const template = "Validate(%q), got %#v, want %#v"
	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			if vv := s.Validate(c.username); !reflect.DeepEqual(vv, c.violations) {
				t.Errorf(template, c.username, vv, c.violations)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	defer leaktest.Check(t)()

	cases := []struct {
		label    string
		username string
		client   usrname.Client
		status   usrname.Status
	}{
		{
			label:    "invalid",
			username: "_obviously_invalid!",
			client:   nil,
			status:   usrname.Invalid,
		}, {
			label:    "notfound",
			username: "dummy",
			client:   mock.Client(http.StatusNotFound, nil),
			status:   usrname.Available,
		}, {
			label:    "ok",
			username: "dummy",
			client:   mock.Client(http.StatusOK, nil),
			status:   usrname.Unavailable,
		}, {
			label:    "other",
			username: "dummy",
			client:   mock.Client(999, nil), // anything other than 200 or 404
			status:   usrname.UnknownStatus,
		}, {
			label:    "clienterror",
			username: "dummy",
			client:   mock.Client(0, errors.New("Oh no!")),
			status:   usrname.UnknownStatus,
		}, {
			label:    "timeouterror",
			username: "dummy",
			client:   mock.Client(0, &timeoutError{}),
			status:   usrname.UnknownStatus,
		},
	}

	const template = "Check(%q), got %q, want %q"
	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			res := s.Check(c.client)(c.username)
			actual := res.Status
			expected := c.status
			if actual != expected {
				t.Errorf(template, c.username, actual, expected)
			}
		})
	}
}

type timeoutError struct {
	error
}

func (*timeoutError) Timeout() bool {
	return true
}
