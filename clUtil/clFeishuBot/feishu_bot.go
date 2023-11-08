package clFeishuBot

import (
	"encoding/json"
	"github.com/cxi7448/cxhttp/clUtil/clHttpClient"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
)

func SendText(_botId string, _content string) {
	hc := clHttpClient.NewClient("https://open.feishu.cn/open-apis/bot/v2/hook/" + _botId)
	hc.SetMethod("POST")
	hc.SetContentType(clHttpClient.ContentJson)

	bytes, err := json.Marshal(clJson.M{
		"msg_type": "text",
		"content": clJson.M{
			"text": _content,
		},
	})
	if err != nil {
		clLog.Error("序列化对象错误: %v", err)
		return
	}
	hc.SetBody(string(bytes))
	hc.Do()
}

func SendPost(_botId string, _content clJson.M) {
	hc := clHttpClient.NewClient("https://open.feishu.cn/open-apis/bot/v2/hook/" + _botId)
	hc.SetContentType(clHttpClient.ContentJson)
	var body = clJson.CreateBy(clJson.M{
		"msg_type": "post",
		"content": clJson.M{
			"post": _content,
		},
	}).ToStr()
	hc.SetBody(body)
	_, err := hc.Do()
	if err != nil {
		return
	}
}
