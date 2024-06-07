package base64_util

import "encoding/base64"

func BytesToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Base64ToBytes(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}
