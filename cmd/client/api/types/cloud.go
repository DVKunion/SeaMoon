package types

type CloudType int8

const (
	ALiYun CloudType = iota + 1
	TencentYun
	HuaweiYun
	BaiduYun
	Sealos
)

type CloudStatus int8

const (
	CREATE CloudStatus = iota + 1
	SUCCESS
	FAILED
	FORBIDDEN
	AUTH_ERROR
	SYNC_ERROR
)
