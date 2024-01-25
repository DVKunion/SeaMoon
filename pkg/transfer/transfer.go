package transfer

import "strings"

type Type string

const (
	HTTP   Type = "http"
	SOCKS5 Type = "socks5"
)

func (t Type) String() string {
	return string(t)
}

func (t Type) Path(p string) string {
	if strings.HasSuffix(p, "/") {
		return p + t.String()
	} else {
		return p + "/" + t.String()
	}
}
