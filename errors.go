package usrname

import (
	"fmt"
)

const (
	nwErrTempl  = "usrname: network error: %v"
	uscErrTempl = "usrname: unexpected status code: %d"
)

type NetworkError struct {
	Cause error
}

func (err *NetworkError) Error() string {
	return fmt.Sprintf(nwErrTempl, err.Cause)
}

type UnexpectedStatusCodeError struct {
	StatusCode int
}

func (err *UnexpectedStatusCodeError) Error() string {
	return fmt.Sprintf(uscErrTempl, err.StatusCode)
}
