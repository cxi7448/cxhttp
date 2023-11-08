package clWeixin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clCrypt"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"strings"
)

// 用户水印
type WxUserInfoWaterMark struct {
	AppId     string `json:"appid"`
	Timestamp uint32 `json:"timestamp"`
}

// 微信用户数据结构
type WxUserInfo struct {
	OpenId    string              `json:"openId"`
	Nickname  string              `json:"nickName"`
	Gender    uint32              `json:"gender"`
	City      string              `json:"city"`
	Province  string              `json:"province"`
	Country   string              `json:"country"`
	AvatarUrl string              `json:"avatarUrl"`
	UnionId   string              `json:"unionId"`
	WaterMark WxUserInfoWaterMark `json:"watermark"`
}

// 微信解密用户数据
func DecryptUserInfo(_buffer, _key, _iv string) (*WxUserInfo, error) {

	//_buffer = strings.ReplaceAll(_buffer, " ", "+")
	_iv = strings.ReplaceAll(_iv, " ", "+")
	_key = strings.ReplaceAll(_key, " ", "+")

	//dData := clCrypt.Base64Decode( _buffer )
	dIv := clCrypt.Base64Decode(_iv)
	dKey := clCrypt.Base64Decode(_key)

	bufferResp := clCrypt.AesCBCDecode([]byte(_buffer), dKey, dIv)
	clLog.Info("bufferResp: %v", bufferResp)
	var uInfo WxUserInfo
	err := json.Unmarshal([]byte(bufferResp), &uInfo)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("反序列化用户数据错误: %v", err))
	}

	return &uInfo, nil
}
