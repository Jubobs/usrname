package sites_test

// see https://github.com/golang/go/issues/24895 for why "_test" in filename

import (
	"github.com/jubobs/username-checker/sites"
	"net/url"
)

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
