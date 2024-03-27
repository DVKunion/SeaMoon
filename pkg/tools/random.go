package tools

import (
	"math/rand"
	"strings"
)

const randomList = "ASDFGHJKLZXCVBNMQWERTYUIOPasdfghjklzxcvbnmqwertyuiop1234567890"
const randomLetterList = "asdfghjklzxcvbnmqwertyuiop"

func GenerateRandomString(length int) string {
	var sb strings.Builder
	sb.Grow(length) // 预分配足够的空间

	for i := 0; i < length; i++ {
		c := randomList[rand.Intn(len(randomList))] // 随机选择一个字符
		sb.WriteByte(c)
	}
	return sb.String()
}

func GenerateRandomLetterString(length int) string {
	var sb strings.Builder
	sb.Grow(length) // 预分配足够的空间

	for i := 0; i < length; i++ {
		c := randomLetterList[rand.Intn(len(randomLetterList))] // 随机选择一个字符
		sb.WriteByte(c)
	}
	return sb.String()
}
