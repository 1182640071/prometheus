package common

import (
	"bytes"
	"math/rand"
	"time"
)

const char = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandChar(size int) string {
	rand.NewSource(time.Now().UnixNano()) // 产生随机种子
	var s bytes.Buffer
	for i := 0; i < size; i ++ {
		s.WriteByte(char[rand.Int63() % int64(len(char))])
	}
	return s.String()
}