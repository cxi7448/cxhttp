package clWeixin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clHttpClient"
	"sync"
	"time"
)

var mAccessToken = ""
var mLastUptime = uint32(0)
var mLocker sync.RWMutex

type GetAccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   uint32 `json:"expires_in"`
	ErrCode     uint32 `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

func GetAccessToken(_appId string, _secret string) (string, error) {

	mLocker.Lock()
	defer mLocker.Unlock()

	if mLastUptime > uint32(time.Now().Unix()) {
		return mAccessToken, nil
	}

	client := clHttpClient.NewClient("https://api.weixin.qq.com/cgi-bin/token")
	client.AddParam("grant_type", "client_credential")
	client.AddParam("appid", _appId)
	client.AddParam("secret", _secret)
	resp, err := client.Do()
	if err != nil {
		return "", err
	}
	var respObj GetAccessTokenResp
	if err := json.Unmarshal([]byte(resp.Body), &respObj); err != nil {
		return "", err
	}
	if respObj.ErrCode != 0 {
		return "", errors.New(fmt.Sprintf("微信返回错误: %v - %v", respObj.ErrCode, respObj.ErrMsg))
	}
	mAccessToken = respObj.AccessToken
	mLastUptime = uint32(time.Now().Unix()) + respObj.ExpiresIn - 60
	return respObj.AccessToken, nil
}
