package mockclient

import (
	"net/http"

	"github.com/jubobs/usrname"
)

type clientFunc func(*http.Request) (*http.Response, error)

func (f clientFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func WithError(err error) usrname.Client {
	do := func(_ *http.Request) (*http.Response, error) {
		return nil, err
	}
	return clientFunc(do)
}

func WithStatusCode(sc int) usrname.Client {
	do := func(_ *http.Request) (*http.Response, error) {
		res := http.Response{
			StatusCode: sc,
		}
		return &res, nil
	}
	return clientFunc(do)
}

func WithStatusCodeAndHeader(sc int, h string, v string) usrname.Client {
	do := func(_ *http.Request) (*http.Response, error) {
		header := map[string][]string{
			h: []string{v},
		}
		res := http.Response{
			StatusCode: sc,
			Header:     header,
		}
		return &res, nil
	}
	return clientFunc(do)
}
