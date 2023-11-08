package clCrypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// RSAWithShaAndBase64 RSA 私钥加签
func RSAWithSha256AndBase64(privateKeyStr, content string) (sign string, err error) {

	cryptoHash := crypto.SHA256

	// 获取私钥
	privateKey, err := StringToRsaPrivateKey(privateKeyStr)
	if err != nil {
		return
	}

	// 获取哈希值
	hashBytes, err := HashBytes(cryptoHash, content)
	if err != nil {
		return
	}

	// 获取签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, cryptoHash, hashBytes)
	if err != nil {
		return
	}

	// base64编码
	sign = base64.StdEncoding.EncodeToString(signature)

	return
}



// hash值转化
func HashBytes(cryptoHash crypto.Hash, content string, append ...[]byte) (hashBytes []byte, err error) {
	hash := cryptoHash.New()
	_, err = hash.Write([]byte(content))

	if err != nil {
		return
	}

	if append == nil || len(append) == 0 {
		hashBytes = hash.Sum(nil)
	} else {
		hashBytes = hash.Sum(append[0])
	}

	return
}



// 私钥字符串转私钥结构
func StringToRsaPrivateKey(privateKeyStr string) (rsaPrivateKey *rsa.PrivateKey, err error) {
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		err = errors.New("private key pem decode error")
		return
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		err = errors.New("privateKey not *rsa.PrivateKey")
	}

	return
}