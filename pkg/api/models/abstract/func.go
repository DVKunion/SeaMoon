package abstract

type FunctionConfig interface {
}

type TunnelConfig struct {
	// 函数配置
	Region   string  `json:"region"`    // 一个隧道只能是一个区域
	CPU      float32 `json:"cpu"`       // CPU 资源
	Memory   int32   `json:"memory"`    // 内存资源
	Instance int32   `json:"instance"`  // 最大实例处理数
	SSRCrypt string  `json:"ssr_crypt"` // ssr 加密方式
	SSRPass  string  `json:"ssr_pass"`  // ssr 密码
	V2rayUid string  `json:"v2ray_uid"` // v2ray_uid
	TLS      bool    `json:"tls"`       // 是否开启 TLS 传输, 开启后自动使用 wss  协议
}
