package xlog

// API 相关日志
const (
	ApiServiceStart = "seamoon api service start"
	ApiServerStart  = "seamoon server start"
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
	SignalListenStart   = "signal start proxy listener"
	SignalListenStop    = "signal stop proxy listener"
	SignalListenRecover = "signal send recover proxy signal"
)

// DB 相关日志
const (
	DatabaseInit       = "database not found, start init ......"
	DatabaseConfigInit = "database not found config table, start init ......"
	DatabaseUserInit   = "database not found user table, start init ......"
)

// SDK 相关日志
const (
	SDKWaitingFCStatus = "sdk waiting for function create success"
)
