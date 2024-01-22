package types

type AuthType int8

const (
	Empty  AuthType = iota + 1 // 无认证
	Admin                      // 后台账户
	Basic                      // HTTP Basic 认证
	Bearer                     // HTTP Bearer 认证

	Signature // FC Signature 认证, 这类认证好处在于认证失败 403 不计次数
	Jwt       // FC Jwt 认证。 需要 jwks
)

func TransAuthType(t string) AuthType {
	switch t {
	case "anonymous":
		return Empty
	case "function":
		return Signature
	case "jwt":
		return Jwt
	}
	return Empty
}
