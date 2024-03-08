package xre

import (
	"regexp"
	"strings"
)

func IsTinyInt(str string) bool {
	match, _ := regexp.Match(`^(\-)?[0-9]{1,3}$`, []byte(str))
	return match
}

func IsInt(str string) bool {
	match, _ := regexp.Match(`^(\-)?[0-9]{1,10}$`, []byte(str))
	return match
}

func IsLong(str string) bool {
	match, _ := regexp.Match(`^(\-)?[0-9]{1,20}$`, []byte(str))
	return match
}

func IsFloat(str string) bool {
	match, _ := regexp.Match(`^(\-)?[0-9]{1,20}(\.[0-9]{1,10})?$`, []byte(str))
	return match
}
func IsDate(_param string) bool {
	match, _ := regexp.Match(`^[0-9]{4}\-[01][0-9]\-([012][0-9]|[3][01])$`, []byte(_param))
	return match
}

func IsTime(_param string) bool {
	match, _ := regexp.Match(`^[0-2][0-9]\:[0-5][0-9]\:[0-5][0-9]$`, []byte(_param))
	return match
}
func IsDateTime(_param string) bool {
	match, _ := regexp.Match(`^[0-9]{4}\-[01][0-9]\-([012][0-9]|[3][01])\s[0-2][0-9]\:[0-5][0-9]\:[0-5][0-9]$`, []byte(_param))
	return match
}

func IsIp(_param string) bool {
	match, _ := regexp.Match(`^[0-9]{1,3}(\.[0-9]{1,3}){3}$`, []byte(_param))
	return match
}
func IsNumberList(_param string) bool {
	match, _ := regexp.Match(`^[0-9]{1,20}(\,[0-9]{1,20})*$`, []byte(_param))
	return match
}
func IsPhone(_param string) bool {
	match, _ := regexp.Match(`^(13|14|15|16|17|18|19)[0-9]{9}$`, []byte(_param))
	return match
}

func IsChinese(ch rune) bool {
	// Unicode范围内的汉字编码区间
	return ch >= '\u4E00' && ch <= '\u9FFF' || ch >= '\u3400' && ch <= '\u4DBF'
}

func IsEn(ch rune) bool {
	return (ch >= 65 && ch <= 90) || (ch >= 97 && ch <= 122)
}

func IsNickname(_nickname string) bool {
	temp := []rune(_nickname)
	Len := len(temp)
	if Len < 2 {
		return false
	}
	chLen := 0
	enLen := 0
	for _, str := range _nickname {
		if IsChinese(str) {
			chLen += 1
		} else if IsEn(str) {
			enLen += 1
		} else {
			return false
		}
	}
	if chLen > 8 || enLen > 16 {
		// 最大8中文  最大16个字母
		return false
	}
	if chLen == Len && Len < 2 {
		// 纯中文最低两位
		return false
	}
	if Len == 2 && chLen == 1 && enLen == 1 {
		return false
	}
	if enLen == Len && Len < 4 {
		// 纯英文最低4位
		return false
	}
	return true
}

func IsIOS(userAgent string) bool {
	return strings.Contains(strings.ToLower(userAgent), "iphone")
}

func IsAndroid(userAgent string) bool {
	return strings.Contains(strings.ToLower(userAgent), "android")
}
