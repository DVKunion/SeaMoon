package enum

type FunctionStatus int8

const (
	FunctionInitializing FunctionStatus = iota + 1 // 初始化
	FunctionActive                                 // 可用
	FunctionInactive                               // 停用
	FunctionError                                  // 不可用
	FunctionWaiting                                // 异常
	FunctionDelete                                 // 删除
)

type FunctionType int8

const (
	FunctionTunnel FunctionType = iota + 1
)

type TunnelType string

const (
	TunnelTypeNULL = "unknown"
	TunnelTypeWST  = "websocket"
	TunnelTypeGRT  = "grpc"
)

var tpMaps = map[TunnelType]string{
	TunnelTypeNULL: "",
	TunnelTypeWST:  "ws",
	TunnelTypeGRT:  "grpc",
}

func TransTunnelType(t string) TunnelType {
	switch t {
	case "websocket":
		return TunnelTypeWST
	case "grpc":
		return TunnelTypeGRT
	default:
		return TunnelTypeNULL
	}
}
func (t TunnelType) String() string {
	return string(t)
}

func (t TunnelType) ToPtr() *string {
	return (*string)(&t)
}

func (t TunnelType) TranAddr(tls bool) string {
	proto := tpMaps[t]
	if tls {
		return proto + "s://"
	}
	return proto + "://"
}
