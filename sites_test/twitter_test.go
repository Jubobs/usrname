package sites_test

import (
	"github.com/jubobs/username-checker/sites"
	"net/http"
	"net/url"
	"testing"
)

var checker = sites.Twitter()

func TestTwitterName(t *testing.T) {
	expected := "Twitter"
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
	for _, c := range cases {
		err := checker.Validate(c.username)
		if (err == nil) != c.valid {
			template := "(Twitter().Validate(%q) == nil) == %t, want %t"
			t.Errorf(template, c.username, err == nil, c.valid)
		}
	}
}

type MockClient struct {
	GetStatusCodeFunc  func(url.URL) (int, error)
	HeadStatusCodeFunc func(url.URL) (int, error)
}

func (m MockClient) GetStatusCode(u url.URL) (int, error) {
	return m.GetStatusCodeFunc(u)
}

func (m MockClient) HeadStatusCode(u url.URL) (int, error) {
	return m.HeadStatusCodeFunc(u)
}

func mockClientHead(statusCode int, err error) sites.Client {
	return MockClient{
		HeadStatusCodeFunc: func(_ url.URL) (int, error) {
			return statusCode, err
		},
	}
}

func Test_that_Check_returns_nil_if_HeadStatusCode_returns_Not_Found_nil(t *testing.T) {
	client := mockClientHead(http.StatusNotFound, nil)
	dummyUsername := "dummy"

	var expected error = nil
	actual := checker.Check(client, dummyUsername)
	if actual != nil {
		t.Errorf("Twitter().Check() == %v, want %v", actual, expected)
	}
}
