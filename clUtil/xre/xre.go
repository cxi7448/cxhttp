package xre

import "regexp"

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
