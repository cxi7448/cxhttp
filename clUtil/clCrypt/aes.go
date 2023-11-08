package clCrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"math/rand"
	"time"
)


// 生成16位随机字符串
func RandomBlock(_len int) []byte {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ/"
	bytes := []byte(str)
	result := []byte{}
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < _len; i++ {
		result = append(result, bytes[r.Int31n(int32(len(str)))])
	}
	return result
}



func pkcs7Unpadding(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, errors.New("尺寸非法")
	}
	if b == nil || len(b) == 0 {
		return nil, errors.New("尺寸非法")
	}


	if len(b)%blocksize != 0 {
		return nil, errors.New("尺寸非法")
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, errors.New("尺寸非法")
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, errors.New("尺寸非法")
		}
	}
	return b[:len(b)-n], nil
}


func pkcs7Padding(ciphertext []byte, blockSize int) ([]byte, error) {

	if blockSize <= 0 {
		return nil, errors.New("尺寸非法")
	}
	if ciphertext == nil || len(ciphertext) == 0 {
		return nil, errors.New("尺寸非法")
	}

	padding := blockSize - len(ciphertext) % blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...), nil
}



// AES解密数据
func AesCBCDecode(_buffer []byte, _key []byte, _iv []byte) string {

	_buffer = bytes.ReplaceAll(_buffer, []byte(" "), []byte("+"))

	dData := Base64Decode( string(_buffer) )
	bufferResp := make([]byte, 0)


	_cipher, err := aes.NewCipher( _key )
	if err != nil {
		return ""
	}

	blockMode := cipher.NewCBCDecrypter(_cipher, _iv[:_cipher.BlockSize()])

	origData := make([]byte, len(dData))
	blockMode.CryptBlocks(origData, dData)

	// 解析
	bufferResp, err = pkcs7Unpadding(origData, aes.BlockSize)
	if len(bufferResp) == 0 {
		return ""
	}

	return string(bufferResp)
}



// AES加密数据
func AesCBCEncode(_buffer string, _key string, iv string) string {
	_cipher, err := aes.NewCipher([]byte(_key))
	if err != nil {
		fmt.Printf("AES加密错误: %v\n", err)
		return ""
	}

	value := _buffer

	// 生成CBC加密对象
	blockMode := cipher.NewCBCEncrypter(_cipher, []byte(iv[:_cipher.BlockSize()]))

	// 填充字节
	dValue, _ := pkcs7Padding([]byte(value), aes.BlockSize)
	origData := make([]byte, len(dValue))

	// 加密
	blockMode.CryptBlocks(origData, dValue)
	return Base64Encode(origData)
}


func Aes256GCMDecode(_aesKey string, _ciphertext string, _nonce string, _associatedData string) string {
		// The key argument should be the AES key, either 16 or 32 bytes
		// to select AES-128 or AES-256.
		key := []byte(_aesKey)
		ciphertext := Base64Decode(_ciphertext)

		nonce := []byte(_nonce)

		block, err := aes.NewCipher(key)
		if err != nil {
			panic(err.Error())
		}

		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			panic(err.Error())
		}

		plaintext, err := aesgcm.Open(nil, nonce, ciphertext, []byte(_associatedData))
		if err != nil {
			panic(err.Error())
		}

		return string(plaintext)
}