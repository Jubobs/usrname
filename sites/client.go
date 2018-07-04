package sites

import (
	"net/http"
)

func NewClient() Client {
	return &simpleClient{}
}

type Client interface {
	Send(*http.Request) (int, error)
}

type simpleClient struct{}

func (*simpleClient) Send(req *http.Request) (int, error) {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer res.Body.Close() // TODO: is the body nil if the method is HEAD?
	return res.StatusCode, err
}
