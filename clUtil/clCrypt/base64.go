package clCrypt

import "encoding/base64"


// base64加密
func Base64Encode(_data []byte) string {
	return base64.StdEncoding.EncodeToString(_data)
}



// base64解密
func Base64Decode(_data string) []byte {
	res, err := base64.StdEncoding.DecodeString(_data)
	if err != nil {
		return []byte{}
	}
	return res
}