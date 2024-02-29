package types

type ProxyStatus int8

const (
	INITIALIZING ProxyStatus = iota + 1
	ACTIVE
	INACTIVE
	ERROR
	WAITING
)
