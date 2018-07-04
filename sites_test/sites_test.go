package sites_test

// see https://github.com/golang/go/issues/24895 for why "_test" in filename

import (
	"github.com/jubobs/username-checker/sites"
	"net/http"
)

type clientFunc func(*http.Request) (int, error)

func (f clientFunc) Send(req *http.Request) (int, error) {
	return f(req)
}

func mockClient(statusCode int, err error) sites.Client {
	send := func(_ *http.Request) (int, error) {
		return statusCode, err
	}
	return clientFunc(send)
}
