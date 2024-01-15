package network

import "errors"

var (
	ErrBadVersion  = errors.New("Bad version")
	ErrBadFormat   = errors.New("Bad format")
	ErrBadAddrType = errors.New("Bad address type")
	ErrShortBuffer = errors.New("Short buffer")
	ErrBadMethod   = errors.New("Bad method")
	ErrAuthFailure = errors.New("Auth failure")
)
