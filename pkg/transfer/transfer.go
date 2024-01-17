package transfer

type Type string

const (
	HTTP   Type = "http"
	SOCKS5 Type = "socks5"
)

func (t Type) String() string {
	return string(t)
}
