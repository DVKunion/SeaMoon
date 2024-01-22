package types

type ProxyType int8

const (
	AUTO ProxyType = iota + 1
	HTTP
	SOCKS5
)

type ProxyStatus int8

const (
	ACTIVE ProxyStatus = iota + 1
	INACTIVE
	ERROR
	COMPLETED
	INITIALIZING
	WAITING
)

type TunnelType int8

const (
	WebSocket TunnelType = iota + 1
	GRPC
)
