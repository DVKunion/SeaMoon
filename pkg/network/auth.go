package network

import (
	"crypto/subtle"
	"encoding/base64"
	"strings"
)

// StrEQ returns whether s1 and s2 are equal
func StrEQ(s1, s2 string) bool {
	return subtle.ConstantTimeCompare([]byte(s1), []byte(s2)) == 1
}

// StrInSlice return whether str in slice
func StrInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// VerifyByMap returns an verifier that verify by an username-password map
func VerifyByMap(users map[string]string) func(string, string) bool {
	return func(username, password string) bool {
		pw, ok := users[username]
		if !ok {
			return false
		}
		return StrEQ(pw, password)
	}
}

// VerifyByHtpasswd returns a verifier that verify by a htpasswd file
//func VerifyByHtpasswd(users string) func(string, string) bool {
//	f, err := htpasswd.New(users, htpasswd.DefaultSystems, nil)
//	if err != nil {
//		xlog.Error("Load htpasswd file failed", "err", err)
//	}
//	return func(username, password string) bool {
//		return f.Match(username, password)
//	}
//}

func HttpBasicAuth(auth string, verify func(string, string) bool) bool {
	prefix := "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return false
	}
	auth = strings.Trim(auth[len(prefix):], " ")
	dc, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return false
	}
	groups := strings.Split(string(dc), ":")
	if len(groups) != 2 {
		return false
	}
	return verify(groups[0], groups[1])
}
