package clFeishuBot

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clHttpClient"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"strings"
	"time"
)

var (
	accessToken string
	updateTime  uint32

	appId      string
	appSecret  string
	encryptKey string
	verifyKey  string
)

// 初始化信息
func Init(_appId string, _appSecret string, _encrypt string, _verify string) {
	appId = _appId
	appSecret = _appSecret
	encryptKey = _encrypt
	verifyKey = _verify
}

// 获取token
func GetAccessToken() string {
	if accessToken != "" && uint32(time.Now().Unix()) < updateTime+360 {
		return accessToken
	}

	hc := clHttpClient.NewClient("https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal")
	hc.SetContentType(clHttpClient.ContentJson)

	hc.SetBody(clJson.CreateBy(clJson.M{
		"app_id":     appId,
		"app_secret": appSecret,
	}).ToStr())

	resp, err := hc.Do()
	if err != nil {
		clLog.Error("Err: %v", err)
		return ""
	}

	var tokenInfo RespGetAccessToken
	if err := json.Unmarshal([]byte(resp.Body), &tokenInfo); err != nil {
		clLog.Error("反序列化数据结构错误: %v", err)
		return ""
	}
	accessToken = tokenInfo.AppAccessToken
	updateTime = uint32(time.Now().Unix())
	return tokenInfo.AppAccessToken
}

// 解密飞书信息
func DecryptFeishuData(encrypt string) (string, error) {
	buf, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", fmt.Errorf("base64StdEncode Error[%v]", err)
	}
	if len(buf) < aes.BlockSize {
		return "", errors.New("cipher  too short")
	}
	keyBs := sha256.Sum256([]byte(encryptKey))
	block, err := aes.NewCipher(keyBs[:sha256.Size])
	if err != nil {
		return "", fmt.Errorf("AESNewCipher Error[%v]", err)
	}
	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]
	// CBC mode always works in whole blocks.
	if len(buf)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buf, buf)
	n := strings.Index(string(buf), "{")
	if n == -1 {
		n = 0
	}
	m := strings.LastIndex(string(buf), "}")
	if m == -1 {
		m = len(buf) - 1
	}
	return string(buf[n : m+1]), nil
}
