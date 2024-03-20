package enum

type ProviderType int8

const (
	ProvTypeALiYun ProviderType = iota + 1
	ProvTypeTencentYun
	ProvTypeHuaweiYun
	ProvTypeBaiduYun
	ProvTypeSealos
)

type ProviderStatus int8

const (
	ProvStatusCrete ProviderStatus = iota + 1
	ProvStatusSuccess
	ProvStatusFailed
	ProvStatusForbidden
	ProvStatusAuthError
	ProvStatusSyncError
	ProvStatusArrearsError
)
