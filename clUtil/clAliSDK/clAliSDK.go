package clAliSDK

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clCrypt"
	"net/url"
	"sort"
	"strings"
)

type Map map[string]interface{}

// 对参数进行编码
func URIEncode(_str string) string {
	StringToSign := url.PathEscape(_str)
	StringToSign = strings.ReplaceAll(StringToSign, "=", "%3D")
	StringToSign = strings.ReplaceAll(StringToSign, "&", "%26")
	StringToSign = strings.ReplaceAll(StringToSign, ":", "%253A")
	return StringToSign
}

// 对参数进行签名
func GlobalSigned(_param Map, _secretKey string) Map {
	var keySort = make([]string, 0)
	for key, _ := range _param {
		if key == "Signature" {
			continue
		}
		keySort = append(keySort, key)
	}

	// 排序
	sort.Strings(keySort)
	var paramStr = strings.Builder{}
	for _, keyName := range keySort {
		var val = _param[keyName]
		if paramStr.Len() > 0 {
			paramStr.WriteString("&")
		}
		paramStr.WriteString(fmt.Sprintf("%v=%v", keyName, val))
	}

	StringToSign := URIEncode(paramStr.String())
	StringToSign = "GET&%2F&" + StringToSign
	//clLog.Debug("string to sign: %v", StringToSign)
	mac := hmac.New(sha1.New, []byte(_secretKey+"&"))
	mac.Write([]byte(StringToSign))
	_param["Signature"] = URIEncode(clCrypt.Base64Encode(mac.Sum(nil)))
	return _param
}
