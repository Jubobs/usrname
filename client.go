package usrname

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
)

var timeout time.Duration

func init() {
	timeout = 1000 * time.Millisecond
}

func NewClient() Client {
	return &simpleClient{}
}

type Client interface {
	Send(*http.Request) (int, error)
}

type simpleClient struct{}

func (*simpleClient) Send(req *http.Request) (int, error) {
	client := http.DefaultClient
	client.Timeout = timeout
	res, err := client.Do(req)
	if err != nil {
		err1 := errors.Wrap(err, "usrname: client failed")
		return http.StatusInternalServerError, err1
	}
	defer res.Body.Close()
	return res.StatusCode, nil
}
