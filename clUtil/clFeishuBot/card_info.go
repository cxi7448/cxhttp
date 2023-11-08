package clFeishuBot

import "github.com/cxi7448/cxhttp/clUtil/clJson"

// 卡片消息
func BuildCardInfoBody(_elements ...interface{}) interface{} {
	return clJson.M{
		"config": clJson.M{
			"wide_screen_mode": true,
			"enable_forward":   false,
			"update_multi":     true,
		},
		"elements": _elements,
	}
}

// 创建按钮列表
func BuildActionsBy(_actions ...interface{}) interface{} {
	return clJson.M{
		"tag":     "action",
		"actions": _actions,
	}
}

// 生成一个文本内容的标题
func BuildCardInfoTitle(_content string) interface{} {
	return clJson.M{
		"tag": "div",
		"text": clJson.M{
			"tag":     "lark_md",
			"content": _content,
		},
	}
}

// 生成一条横线
func BuildCardInfoHR() interface{} {
	return clJson.M{
		"tag": "hr",
	}
}

// 生成一个文本内容Markdown
func BuildCardInfoMD(_content string) interface{} {
	return clJson.M{
		"tag":     "markdown",
		"content": _content,
	}
}

// 生成一个按钮
func BuildButton(_title string, _ac string, _param string) interface{} {
	return clJson.M{
		"tag": "button",
		"text": clJson.M{
			"tag":     "plain_text",
			"content": _title,
		},
		"type": "default",
		"value": clJson.M{
			"ac":    _ac,
			"param": _param,
		},
	}
}
