package sites

import (
	"net/http"
	"net/url"
)

func NewClient() Client {
	return &simpleClient{}
}

type Client interface {
	GetStatusCode(u url.URL) (int, error)
	HeadStatusCode(u url.URL) (int, error)
}

type simpleClient struct{}

func (*simpleClient) HeadStatusCode(u url.URL) (int, error) {
	res, err := http.Head(u.String())
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return res.StatusCode, err
}

func (*simpleClient) GetStatusCode(u url.URL) (int, error) {
	res, err := http.Get(u.String())
	if err != nil {
		return http.StatusInternalServerError, err
	}
	res.Body.Close()
	return res.StatusCode, err
}
