package types

type AuthType int8

const (
	Basic AuthType = iota + 1
	Bearer
	Admin // 用于登陆
)
