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
	SignalStartProxy  = "signal start proxy listener success"
	SignalSpeedProxy  = "signal speed proxy success"
	SignalStopProxy   = "signal stop proxy success"
	SignalDeleteProxy = "signal delete proxy success"

	SignalSyncProvider   = "signal sync provider success"
	SignalDeleteProvider = "signal delete provider success"

	SignalDeployTunnel = "signal deploy tunnel success"
	SignalStopTunnel   = "signal stop tunnel success"
	SignalDeleteTunnel = "signal delete tunnel success"
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

// Signal 相关告警
const (
	SignalMissOperationWarn = "signal received unchanged status"
)

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

	// ServiceDBNeedParamsError error
	ServiceDBNeedParamsError   = "service operation params missing error"
	ServiceDBUpdateStatusError = "service update status error"
	ServiceDBUpdateFiledError  = "service update field error"
	ServiceDBDeleteError       = "service delete status error"
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
)

// SIGNAL 相关错误
const (
	SignalGetObjError       = "signal try to get object from db error"
	SignalUpdateObjError    = "signal try to update object from db error"
	SignalListenerError     = "signal listener unexpect error"
	SignalRecoverProxyError = "signal try recover active proxy error"
	SignalSpeedProxyError   = "signal try test speed error"
	SignalSyncProviderError = "signal try sync provider error"
	SignalDeployTunError    = "signal try deploy tunnel error"
	SignalStopTunError      = "signal try stop tunnel error"
	SignalDeleteTunError    = "signal try delete tunnel error"
)
