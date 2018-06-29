package sites

import (
	"fmt"
)

type networkError struct {
	cause error
}

func (err *networkError) Error() string {
	const template = "sites: network error: %v"
	return fmt.Sprintf(template, err.cause)
}

func (err *networkError) Cause() error {
	return err.cause
}

func IsNetworkError(err error) bool {
	_, ok := err.(*networkError)
	return ok
}

type unexpectedStatusCodeError struct {
	statusCode int
}

func (err *unexpectedStatusCodeError) Error() string {
	const template = "sites: unexpected status code: %d"
	return fmt.Sprintf(template, err.statusCode)
}

func IsUnexpectedStatusCodeError(err error) bool {
	_, ok := err.(*unexpectedStatusCodeError)
	return ok
}

func (err *unexpectedStatusCodeError) StatusCode() int {
	return err.statusCode
}
