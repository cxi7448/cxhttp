package clWeixin

import (
	"encoding/json"
	"github.com/cxi7448/cxhttp/clUtil/clHttpClient"
)

type Code2SessionResp struct {
	OpenId     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionId    string `json:"unionid"`     // 开放平台唯一标识
	ErrCode    uint32 `json:"errcode"`     // 错误码, 0=成功
	ErrMsg     string `json:"errmsg"`      // 错误信息

}

// 提供微信服务端接口开发的封装类
func Code2Session(_appId string, _secret string, _code string) (*Code2SessionResp, error) {
	client := clHttpClient.NewClient("https://api.weixin.qq.com/sns/jscode2session")
	client.AddParam("appid", _appId)
	client.AddParam("secret", _secret)
	client.AddParam("js_code", _code)
	client.AddParam("grant_type", "authorization_code")
	resp, err := client.Do()
	if err != nil {
		return nil, err
	}
	var respObj Code2SessionResp
	if err := json.Unmarshal([]byte(resp.Body), &respObj); err != nil {
		return nil, err
	}
	return &respObj, nil
}
