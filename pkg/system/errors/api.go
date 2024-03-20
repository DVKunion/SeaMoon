package errors

import (
	"github.com/pkg/errors"

	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

type APIError struct {
	e error
	s string
}

type APIErrorMsg string

func (a APIError) Error() string {
	xlog.Error("API", a.s, "err", a.e)
	return a.s
}

func ApiError(msg APIErrorMsg, e ...error) APIError {
	res := APIError{
		s: string(msg),
		e: errors.New(""),
	}
	if len(e) > 0 {
		for _, es := range e {
			res.e = errors.Wrap(res.e, es.Error())
		}
	}
	return res
}
