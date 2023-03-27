package text

import "math/rand"

var (
	allBytes    = []byte("abcdefghijklmnopqrstuvwxyz1234567890")
	allBytesLen = len(allBytes)
)

// RandString 获取指定长度的随机字符串
func RandString(n int) string {
	randBytes := make([]byte, n)
	for i := range randBytes {
		randBytes[i] = allBytes[rand.Intn(allBytesLen)]
	}
	return string(randBytes)
}

// GenerateUUID ...
func GenerateUUID() string {
	return RandString(16)
}
