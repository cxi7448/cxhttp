package clWeixin

import (
	"github.com/cxi7448/cxhttp/clUtil/clHttpClient"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
)

// 发送订阅消息
func SendSubscribeMessage(_appId, _secret, _openId, _tempId string, data clJson.M) error {

	accessToken, err := GetAccessToken(_appId, _secret)
	if err != nil {
		return err
	}
	httpClient := clHttpClient.NewClient("https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=" + accessToken)

	//httpClient.AddParam("access_token", GetAccessToken())
	httpClient.AddParam("touser", _openId)
	httpClient.AddParam("template_id", _tempId)
	httpClient.AddParam("page", "/")
	httpClient.AddParam("data", data)
	httpClient.AddParam("miniprogram_state", "formal")

	httpClient.SetContentType(clHttpClient.ContentJson)
	resp, err := httpClient.Do()
	if err != nil {
		return err
	}
	clLog.Debug("发送订阅消息回应: %v", resp)
	return nil
}
