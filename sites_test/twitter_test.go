package sites_test

import (
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
	}
	const template = "(len(Twitter().Validate(%q))) is %d, but expected %d"
	for _, c := range cases {
		if vs := checker.Validate(c.username); len(vs) != c.noOfViolations {
			t.Errorf(template, c.username, len(vs), c.noOfViolations)
		}
	}
}

func TestCheckNotFound(t *testing.T) {
	// Given
	client := mockClientHead(http.StatusNotFound, nil)
	const dummyUsername = "dummy"

	// When
	available, err := checker.IsAvailable(client)(dummyUsername)

	// Then
	if !(err == nil && available) {
		const template = "Twitter().IsAvailable(%q) == (%t, %v), but expected (true, <nil>)"
		t.Errorf(template, dummyUsername, available, err)
	}
}
