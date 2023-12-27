package xhash

import (
	"github.com/speps/go-hashids/v2"
)

var (
	salt     = "jnfdbasgoew"
	alphabet = "QWERTYUIOPASDFGHJKLZXCVBNM123456789"
	minLen   = 6
)

func Init(_salt, _alphabet string, _minLen int) {
	salt = _salt
	alphabet = _alphabet
	minLen = _minLen
}

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
