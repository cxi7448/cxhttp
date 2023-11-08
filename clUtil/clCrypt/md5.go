package clCrypt

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(str []byte) string {
	h := md5.New()
	h.Write(str) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}