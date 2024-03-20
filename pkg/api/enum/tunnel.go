package enum

type TunnelStatus int8

const (
	TunnelInitializing TunnelStatus = iota + 1 // 初始化
	TunnelActive                               // 可用
	TunnelInactive                             // 停用
	TunnelError                                // 不可用
	TunnelWaiting                              // 异常
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

func (t TunnelType) TranAddr(tls bool) string {
	proto := tpMaps[t]
	if tls {
		return proto + "s://"
	}
	return proto + "://"
}
