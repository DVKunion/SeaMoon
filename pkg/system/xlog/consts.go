package xlog

// API 相关日志
const (
	ApiServerStart = "api service start"
)

// SERVICE 相关日志
const (
	ServiceSocks5ConnectServer = "service socks5 server handle connect"
	ServiceSocks5Establish     = "service socks5 establish connect"
	ServiceSocks5DisConnect    = "service socks5 disconnect"
	ServiceTorConnectServer    = "service tor server handle connect"
	ServiceTorDisConnect       = "service tor disconnect"
)

// SIGNAL 相关日志

const (
	SignalListenStart = "signal start proxy listener"
	SignalListenStop  = "signal stop proxy listener"
)
