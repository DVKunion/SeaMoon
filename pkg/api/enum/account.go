package enum

type AccountType int8

const (
	Admin AccountType = iota + 1 // 管理员 用于登陆后台的账户
	Xray
	Cloud
)

type AuthType int8

const (
	AdminAuth AuthType = iota + 1
)

const (
	// XrayAuthUserPass xray 账户类型
	XrayAuthUserPass  AuthType = iota + 5 // xray http / socks 代理账户
	XrayAuthEmailPass                     // xray shadowsocks / trojan 代理账户
	XrayAuthIdEncrypt                     // xray vmess / vless 代理账户
)

const (
	// CloudAuthKeyMod  云账户认证类型
	CloudAuthKeyMod AuthType = iota + 10 // AK / SK 模式
	CloudAuthCert                        // 证书模式
)

const (
	// CloudFCAuthEmpty 函数认证类型
	CloudFCAuthEmpty     AuthType = iota + 20 // 无认证
	CloudFCAuthSignature                      // FC Signature 认证, 这类认证好处在于认证失败 403 不计次数
	CloudFCAuthJwt                            // FC Jwt 认证。 需要 jwks
	CloudFCAuthParis                          // SCF 网管密钥对认证
)

func TransCloudFCAuthType(t string) AuthType {
	switch t {
	case "anonymous", "NONE":
		return CloudFCAuthEmpty
	case "function":
		return CloudFCAuthSignature
	case "jwt":
		return CloudFCAuthJwt
	}
	return CloudFCAuthEmpty
}
