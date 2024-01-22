package types

type CloudType int8

const (
	ALiYun CloudType = iota + 1
	TencentYun
	HuaweiYun
	BaiduYun
	Sealos
)
