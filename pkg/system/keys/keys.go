package keys

var key []byte

func SetGlobalKey(k []byte) {
	key = k
}

func GetGlobalKey() []byte {
	return key
}
