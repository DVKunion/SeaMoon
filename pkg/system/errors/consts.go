package errors

// API 相关错误
const (
	ApiCommonError   = "api error"
	ApiParamsError   = "api request params error"
	ApiParamsExist   = "api request params already exist"
	ApiParamsRequire = "api require params missing"
	ApiServiceError  = "api request service error"
	ApiAuthError     = "api request auth error"
	ApiAuthRequire   = "api require auth missing"
	ApiAuthLimit     = "api request auth limit"
)

// SERVICE 相关错误
const (
	ServiceProtocolNotSupportError = "service accept proto not support"
	ServiceSocks5ReadMethodError   = "service socks5 read method proto error"
	ServiceSocks5WriteMethodError  = "service socks5 write method proto error"
	ServiceSocks5ReadCmdError      = "service socks5 read command proto error"
	ServiceSocks5DailError         = "service socks5 dial to remote error"
	ServiceSocks5ReplyError        = "service socks5 write to reply success error"
)

// NETWORK 相关错误
const (
	NetworkVersionError   = "network bad version error"
	NetworkAddrTypeError  = "network bad address type error"
	NetworkMethodError    = "network bad method error"
	NetworkTransportError = "network transport error"
)

// LISTENER 相关错误
const (
	ListenerAcceptError = "listener unexpect error"
	ListenerDailError   = "listener dail to remote error"
	ListenerLagError    = "listener latency calculation error"
)

// SIGNAL 相关错误
const (
	SignalGetObjError    = "signal try to get object from db error"
	SignalUpdateObjError = "signal try to update object from db error"
	SignalListenerError  = "signal listener unexpect error"
	SignalSpeedTestError = "signal try test speed error"
)
