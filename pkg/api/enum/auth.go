package enum

type AuthType int8

const (
	AuthEmpty  AuthType = iota + 1 // 无认证
	AuthAdmin                      // 后台账户
	AuthBasic                      // HTTP Basic 认证
	AuthBearer                     // HTTP Bearer 认证

	AuthSignature // FC Signature 认证, 这类认证好处在于认证失败 403 不计次数
	AuthJwt       // FC Jwt 认证。 需要 jwks
	AuthParis     // SCF 网管密钥对认证
)

func TransAuthType(t string) AuthType {
	switch t {
	case "anonymous", "NONE":
		return AuthEmpty
	case "function":
		return AuthSignature
	case "jwt":
		return AuthJwt
	}
	return AuthEmpty
}
