package common

import (
	"bytes"
)

//----------------方法工具----------------
func IfUillStr(v1 string, v2 string) string {
	if len(v1) > 0 {
		return v1
	}
	return v2
}

func StringConcat(str ...string) string {
	if len(str) <= 0 {
		return ""
	}
	var buffer bytes.Buffer
	for _, s := range str {
		buffer.WriteString(s)
	}
	return buffer.String()
}
