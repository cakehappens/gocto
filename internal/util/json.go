package util

const JSONNull string = "null"

func IsJSONNull(data []byte) bool {
	return string(data) == JSONNull
}
