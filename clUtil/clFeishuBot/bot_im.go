package clFeishuBot

import (
	"github.com/cxi7448/cxhttp/clUtil/clHttpClient"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
)

// 发送文本消息到群组
func SendTextMessageToGroup(_openId string, _content string) {
	token := GetAccessToken()
	if token == "" {
		return
	}
	hc := clHttpClient.NewClient("https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=chat_id")
	hc.SetMethod("POST")
	hc.SetContentType(clHttpClient.ContentJson)
	hc.AddHeader("Authorization", "Bearer "+token)
	hc.SetBody(clJson.CreateBy(clJson.M{
		"receive_id": _openId,
		"msg_type":   "text",
		"content": clJson.CreateBy(clJson.M{
			"text": _content,
		}).ToStr(),
	}).ToStr())
	resp, err := hc.Do()
	if err != nil {
		clLog.Error("发送消息到群组错误: %v", err)
		return
	}

	clLog.Debug("发送消息到群组返回: %v", resp.Body)
}

// 发送文本消息到群组At某人
func SendTextMessageToGroupAtAll(_openId string, _content string) {
	token := GetAccessToken()
	if token == "" {
		return
	}
	hc := clHttpClient.NewClient("https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=chat_id")
	hc.SetMethod("POST")
	hc.SetContentType(clHttpClient.ContentJson)
	hc.AddHeader("Authorization", "Bearer "+token)
	hc.SetBody(clJson.CreateBy(clJson.M{
		"receive_id": _openId,
		"msg_type":   "text",
		"content": clJson.CreateBy(clJson.M{
			"text": "<at user_id=\"all\">所有人</at> " + _content,
		}).ToStr(),
	}).ToStr())
	resp, err := hc.Do()
	if err != nil {
		clLog.Error("发送消息到群组错误: %v", err)
		return
	}

	clLog.Debug("发送消息到群组返回: %v", resp.Body)
}

// 发送文本到用户
func SendTextMessageToUser(_openId string, _content string) {

	token := GetAccessToken()
	if token == "" {
		return
	}

	hc := clHttpClient.NewClient("https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=open_id")
	hc.SetMethod("POST")
	hc.SetContentType(clHttpClient.ContentJson)
	hc.AddHeader("Authorization", "Bearer "+token)
	hc.SetBody(clJson.CreateBy(clJson.M{
		"receive_id": _openId,
		"msg_type":   "text",
		"content": clJson.CreateBy(clJson.M{
			"text": _content,
		}).ToStr(),
	}).ToStr())
	resp, err := hc.Do()
	if err != nil {
		clLog.Error("发送消息到用户错误: %v", err)
		return
	}

	clLog.Debug("发送消息到用户返回: %v", resp.Body)
}

// 获取机器人所在群组列表
func GetRobotInGroupList() {
	var token = GetAccessToken()
	if token == "" {
		clLog.Error("无法获取机器人所在群组列表, access_token为空!")
		return
	}
	hc := clHttpClient.NewClient("https://open.feishu.cn/open-apis/im/v1/chats")
	hc.SetMethod("GET")
	hc.AddHeader("Authorization", "Bearer "+token)
	resp, err := hc.Do()
	if err != nil {
		clLog.Error("获取机器人所在群组信息错误: %v", err)
		return
	}

	clLog.Debug("获取机器人所在群组信息: %v", resp.Body)
}

// 发送消息到群组中
func SendMessageToGroup(_openId string, _content interface{}) {
	token := GetAccessToken()
	if token == "" {
		clLog.Error("发送消息到群组失败: access_token为空!")
		return
	}
	hc := clHttpClient.NewClient("https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=chat_id")
	hc.SetMethod("POST")
	hc.SetContentType(clHttpClient.ContentJson)
	hc.AddHeader("Authorization", "Bearer "+token)

	hc.SetBody(clJson.CreateBy(clJson.M{
		"receive_id": _openId,
		"msg_type":   "interactive",
		"content":    clJson.CreateBy(_content).ToStr(),
	}).ToStr())
	resp, err := hc.Do()
	if err != nil {
		clLog.Error("发送消息到群组错误: %v", err)
		return
	}

	clLog.Debug("发送消息到群组返回: %v", resp.Body)
}

// 发送消息给指定的人
func SendMessageToPerson(_chatId string, _userId string, _card interface{}) {
	token := GetAccessToken()
	if token == "" {
		clLog.Error("发送消息到群组失败: access_token为空!")
		return
	}

	hc := clHttpClient.NewClient("https://open.feishu.cn/open-apis/ephemeral/v1/send")
	hc.SetContentType(clHttpClient.ContentJson)
	hc.SetMethod("POST")
	hc.AddHeader("Authorization", "Bearer "+token)
	hc.SetBody(clJson.CreateBy(
		clJson.M{
			"chat_id":  _chatId,
			"open_id":  _userId,
			"msg_type": "interactive",
			"card":     _card,
		},
	).ToStr())

	res, err := hc.Do()
	if err != nil {
		clLog.Error("发送消息给指定的人错误: %v", err)
		return
	}

	clLog.Info("发送消息给指定的人返回: %v", res.Body)
}

// 更新消息
func UpdateMessageTo(_messageId string, _content interface{}) {

	token := GetAccessToken()
	if token == "" {
		return
	}

	var url = "https://open.feishu.cn/open-apis/im/v1/messages/" + _messageId
	hc := clHttpClient.NewClient(url)
	hc.SetMethod("PATCH")
	hc.SetContentType(clHttpClient.ContentJson)
	hc.AddHeader("Authorization", "Bearer "+token)
	hc.SetBody(clJson.CreateBy(clJson.M{
		"content": clJson.CreateBy(_content).ToStr(),
	}).ToStr())
	resp, err := hc.Do()
	if err != nil {
		clLog.Error("更新消息(%v)错误: %v", err)
		return
	}

	clLog.Debug("更新消息(%v)成功返回: %v", _messageId, resp.Body)
}
