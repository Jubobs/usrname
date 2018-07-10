package usrname

import (
	"net/http"

	"github.com/pkg/errors"
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
		err1 := errors.Wrap(err, "usrname: client failed")
		return http.StatusInternalServerError, err1
	}
	defer res.Body.Close() // TODO: is the body nil if the method is HEAD?
	return res.StatusCode, err
}
