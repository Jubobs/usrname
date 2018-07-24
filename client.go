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

// Client implementations must close the Body of the Response (if non-nil)
// before returning it.
type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type simpleClient struct{}

func (*simpleClient) Do(req *http.Request) (*http.Response, error) {
	client := http.DefaultClient
	client.Timeout = timeout
	res, err := client.Do(req)
	if err != nil {
		err := errors.Wrap(err, "usrname: client failed")
		return nil, err
	}
	defer res.Body.Close()
	return res, nil
}
