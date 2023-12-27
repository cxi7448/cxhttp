package xhash

import (
	"github.com/speps/go-hashids/v2"
)

var (
	salt     = "salt"
	alphabet = "qwertyuiopasdfghjklzxcvbnm"
	minLen   = 6
)

// 解析user_id
func Decode(str string) uint64 {
	if len(str) < minLen {
		return 0
	}
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLen
	hd.Alphabet = alphabet
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return 0
	}
	s, e := h.DecodeWithError(str)
	if e != nil {
		return 0
	}
	return uint64(s[0])
}

func Encode(uid uint64) string {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLen
	hd.Alphabet = alphabet
	h, _ := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{int(uid)})
	return e
}
