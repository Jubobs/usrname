package sites

import (
	"net/http"
	"net/url"
)

func NewClient() Client {
	return &simpleClient{}
}

type Client interface {
	HeadStatusCode(u url.URL) (int, error)
}

type simpleClient struct{}

func (*simpleClient) HeadStatusCode(u url.URL) (int, error) {
	res, err := http.Head(u.String())
	return res.StatusCode, err
}
