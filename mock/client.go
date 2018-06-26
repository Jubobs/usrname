package mock

import (
	"github.com/jubobs/username-checker/sites"
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
