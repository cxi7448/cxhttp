package clCrypt

import (
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/farmerx/gorsa"
	"strings"
)

// 简单的RSA加密
func RSAEncode(_data []byte, _publicKey string) (error, string) {
	_publicKey = WrapPublicKey(_publicKey)

	if err := gorsa.RSA.SetPublicKey(_publicKey); err != nil {
		return err, ""
	}

	signBuffer, err := gorsa.RSA.PubKeyENCTYPT(_data)
	if err != nil {
		return err, ""
	}

	// 最后base64加密
	return nil, Base64Encode(signBuffer)
}

// 简单的RSA解密
func RSADecode(_data string, _privateKey string) (error, []byte) {
	_privateKey = WrapPrivateKey(_privateKey)

	if err := gorsa.RSA.SetPrivateKey(_privateKey); err != nil {
		return err, nil
	}

	rawEncrypt := Base64Decode(_data)
	rawDecode, err := gorsa.RSA.PriKeyDECRYPT(rawEncrypt)
	if err != nil {
		clLog.Error("错误: %v", err)
		return err, nil
	}

	return nil, rawDecode
}

func WrapPublicKey(_keyContent string) string {
	keyStr := strings.Builder{}
	if !strings.HasPrefix(_keyContent, "-----BEGIN PUBLIC KEY-----") {
		keyStr.WriteString("-----BEGIN PUBLIC KEY-----\n")
	}

	for i := 0; i < len(_keyContent); i++ {
		keyStr.WriteByte(_keyContent[i])
	}

	if !strings.HasSuffix(_keyContent, "-----END PUBLIC KEY-----") {
		keyStr.WriteString("\n-----END PUBLIC KEY-----")
	}
	return keyStr.String()
}

func WrapPrivateKey(_keyContent string) string {
	keyStr := strings.Builder{}
	if !strings.HasPrefix(_keyContent, "-----BEGIN PRIVATE KEY-----") {
		keyStr.WriteString("-----BEGIN PUBLIC KEY-----\n")
	}
	keyStr.WriteString("-----BEGIN PRIVATE KEY-----\n")

	for i := 0; i < len(_keyContent); i++ {
		keyStr.WriteByte(_keyContent[i])
	}
	if !strings.HasSuffix(_keyContent, "-----END PRIVATE KEY-----") {
		keyStr.WriteString("\n-----END PRIVATE KEY-----")
	}
	return keyStr.String()
}
