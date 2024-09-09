package abstract

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DVKunion/SeaMoon/pkg/system/tools"
)

func TestUserPassAuth(t *testing.T) {
	// check single auth
	a1 := &UserPassAuth{
		User: "admin",
		Pass: "admin",
	}

	assert.Equal(t, "{\"user\":\"admin\",\"pass\":\"admin\"}", tools.MarshalString(a1))

	// check single auth with level nil
	a2 := &UserPassAuth{
		User:  "admin",
		Pass:  "admin",
		Level: nil,
	}

	assert.Equal(t, "{\"user\":\"admin\",\"pass\":\"admin\"}", tools.MarshalString(a2))

	// check single auth with level
	a3 := &UserPassAuth{
		User:  "admin",
		Pass:  "admin",
		Level: tools.IntPtr(0),
	}

	assert.Equal(t, "{\"user\":\"admin\",\"pass\":\"admin\",\"level\":0}", tools.MarshalString(a3))

	// check auth list
	al1 := AuthList{}

	assert.Equal(t, "[]", tools.MarshalString(al1))

	al2 := AuthList{a1}
	assert.Equal(t, "[{\"user\":\"admin\",\"pass\":\"admin\"}]", tools.MarshalString(al2))

	al3 := AuthList{a1, a2, a3}
	assert.Equal(t, "[{\"user\":\"admin\",\"pass\":\"admin\"},{\"user\":\"admin\",\"pass\":\"admin\"},{\"user\":\"admin\",\"pass\":\"admin\",\"level\":0}]", tools.MarshalString(al3))
}

func TestEmailPassAuth(t *testing.T) {
	// check single auth
	// for shadowsocks / torjan inbound, sometimes may not have method field
	a1 := &EmailPassAuth{
		Email:    "admin@admin.com",
		Password: "admin",
		Level:    0,
	}

	assert.Equal(t, "{\"password\":\"admin\",\"level\":0,\"email\":\"admin@admin.com\"}", tools.MarshalString(a1))

	// check shadowsocks with method case
	a2 := &EmailPassAuth{
		Email:    "admin@admin.com",
		Password: "admin",
		Method:   "aes-256-gcm",
		Level:    0,
	}

	assert.Equal(t, "{\"password\":\"admin\",\"method\":\"aes-256-gcm\",\"level\":0,\"email\":\"admin@admin.com\"}", tools.MarshalString(a2))

	// check auth lists
	al1 := &AuthList{a1}
	assert.Equal(t, "[{\"password\":\"admin\",\"level\":0,\"email\":\"admin@admin.com\"}]", tools.MarshalString(al1))

	al2 := &AuthList{a1, a2}
	assert.Equal(t, "[{\"password\":\"admin\",\"level\":0,\"email\":\"admin@admin.com\"},{\"password\":\"admin\",\"method\":\"aes-256-gcm\",\"level\":0,\"email\":\"admin@admin.com\"}]", tools.MarshalString(al2))
}

func TestIdEncryptAuth(t *testing.T) {
	a1 := IdEncryptAuth{
		Id:    "1234-5678-9101-1121",
		Level: 0,
		Email: "admin@admin.com",
	}
	assert.Equal(t, "{\"id\":\"1234-5678-9101-1121\",\"level\":0,\"email\":\"admin@admin.com\"}", tools.MarshalString(a1))
}
