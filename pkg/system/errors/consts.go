package errors

// API 相关错误
const (
	ApiCommonError   = "api error"
	ApiServeError    = "api serve error"
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
	// ServiceProtocolNotSupportError * local service
	ServiceProtocolNotSupportError = "service accept proto not support"
	ServiceSocks5ReadMethodError   = "service socks5 read method proto error"
	ServiceSocks5WriteMethodError  = "service socks5 write method proto error"
	ServiceSocks5ReadCmdError      = "service socks5 read command proto error"
	ServiceSocks5DailError         = "service socks5 dial to remote error"
	ServiceSocks5ReplyError        = "service socks5 write to reply success error"

	// ServiceError remote server
	ServiceError          = "service remote error"
	ServiceStatusError    = "service remote status error "
	ServiceTransportError = "service remote transport error"
	ServiceV2rayInitError = "service remote init v2ray error"
)

// NETWORK 相关错误
const (
	NetworkVersionError   = "network bad version error"
	NetworkAddrTypeError  = "network bad address type error"
	NetworkMethodError    = "network bad method error"
	NetworkTransportError = "network transport error"
)

// SDK 相关错误
const (
	SDKFCInfoError           = "sdk get function info error"
	SDKFCDetailError         = "sdk get function detail error"
	SDKTriggerError          = "sdk get function trigger error"
	SDKTriggerUnmarshalError = "sdk unmarshal function trigger error"
)

// LISTENER 相关错误
const (
	ListenerAcceptError = "listener unexpect error"
	ListenerDailError   = "listener dail to remote error"
	ListenerLagError    = "listener latency calculation error"
)

// SIGNAL 相关错误
const (
	SignalGetObjError       = "signal try to get object from db error"
	SignalUpdateObjError    = "signal try to update object from db error"
	SignalListenerError     = "signal listener unexpect error"
	SignalSpeedTestError    = "signal try test speed error"
	SignalRecoverProxyError = "signal try recover active proxy error"
)
