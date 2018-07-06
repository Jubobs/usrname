package mock

import (
	"net/http"

	"github.com/jubobs/username-checker/sites"
)

type ClientFunc func(*http.Request) (int, error)

func (f ClientFunc) Send(req *http.Request) (int, error) {
	return f(req)
}

func Client(statusCode int, err error) sites.Client {
	send := func(_ *http.Request) (int, error) {
		return statusCode, err
	}
	return ClientFunc(send)
}
