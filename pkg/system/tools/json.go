package tools

import (
	"encoding/json"
	"reflect"
)

func MarshalString(a any) string {
	b, err := json.Marshal(a)
	if err != nil {
		switch cast := reflect.TypeOf(a); cast.Kind() {
		case reflect.Array:
			return "[]"
		default:
			return ""
		}
	}
	return string(b)
}
