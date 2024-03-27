package errors

import (
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

type APIError struct {
	e error
	s string
}

type APIErrorMsg string

func (a APIError) Error() string {
	xlog.Error(a.s, "err", a.e)
	return a.s
}

func ApiError(msg APIErrorMsg, e error) APIError {
	return APIError{
		s: string(msg),
		e: e,
	}
}
