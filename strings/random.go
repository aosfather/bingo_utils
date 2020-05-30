package strings

import (
	"math/rand"
	"time"
)

var (
	_CHARS   []byte = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	_NUMBERS []byte = []byte("0123456789")
)

//RandomStr 随机生成字符串
func RandomStr(length int) string {
	return randomstring(length, _CHARS)
}

func RandomNumber(length int) string {
	return randomstring(length, _NUMBERS)
}

func randomstring(length int, array []byte) string {
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	charlen := len(array)
	for i := 0; i < length; i++ {
		result = append(result, array[r.Intn(charlen)])
	}
	return string(result)
}
