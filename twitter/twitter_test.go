package twitter_test

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/fortytw2/leaktest"
	"github.com/jubobs/usrname"
	"github.com/jubobs/usrname/mockclient"
	"github.com/jubobs/usrname/twitter"
)

var checker = twitter.New()

func TestName(t *testing.T) {
	defer leaktest.Check(t)()
	const expected = "Twitter"
	actual := checker.Name()
	if actual != expected {
		template := "got %q, want %q"
		t.Errorf(template, actual, expected)
	}
}

func TestLink(t *testing.T) {
	defer leaktest.Check(t)()
	const username = "foobar"
	const expected = "https://twitter.com/" + username
	actual := checker.Link(username)
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
					Whitelist: checker.Whitelist(),
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
					Pattern: checker.IllegalPattern().String(),
					At:      []int{0, 7},
				},
			},
		}, {
			"twittermixedcase",
			"not_ok_TwitteR",
			[]usrname.Violation{
				&usrname.IllegalSubstring{
					Pattern: checker.IllegalPattern().String(),
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
					Whitelist: checker.Whitelist(),
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
			if vv := checker.Validate(c.username); !reflect.DeepEqual(vv, c.violations) {
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
			client:   mockclient.WithStatusCode(http.StatusNotFound),
			status:   usrname.Available,
		}, {
			label:    "ok",
			username: "dummy",
			client:   mockclient.WithStatusCode(http.StatusOK),
			status:   usrname.Unavailable,
		}, {
			label:    "other", // than 200, 302, 404
			username: "dummy",
			client:   mockclient.WithStatusCode(999),
			status:   usrname.UnknownStatus,
		}, {
			label:    "clienterror",
			username: "dummy",
			client:   mockclient.WithError(errors.New("Oh no!")),
			status:   usrname.UnknownStatus,
		}, {
			label:    "timeouterror",
			username: "dummy",
			client:   mockclient.WithError(&timeoutError{}),
			status:   usrname.UnknownStatus,
		}, {
			label:    "foundnolocation",
			username: "dummy",
			client:   mockclient.WithStatusCode(http.StatusFound),
			status:   usrname.UnknownStatus,
		}, {
			label:    "foundunexpectedlocation",
			username: "dummy",
			client: mockclient.WithStatusCodeAndHeader(
				http.StatusFound,
				"location",
				"http://unexpected",
			),
			status: usrname.UnknownStatus,
		}, {
			label:    "foundexpectedlocation",
			username: "dummy",
			client: mockclient.WithStatusCodeAndHeader(
				http.StatusFound,
				"location",
				"https://twitter.com/account/suspended",
			),
			status: usrname.Unavailable,
		},
	}

	const template = "Check(%q), got %q, want %q"
	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			res := checker.Check(c.client)(c.username)
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
