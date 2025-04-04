package clResponse

import (
	"encoding/json"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
)

// 标准输出json
func JCode(_code uint32, _msg string, _data interface{}) string {
	resp, _ := json.Marshal(SkyResp{
		Code: _code,
		Msg:  _msg,
		Data: _data,
	})
	return string(resp)
}

// 成功
func Success(data ...interface{}) string {
	var _data interface{}
	if len(data) > 0 {
		_data = data[0]
	}
	resp, _ := json.Marshal(SkyResp{
		Code: 0,
		Msg:  "ok",
		Data: _data,
	})
	return string(resp)
}

func Error(message string, args ...interface{}) string {
	resp, _ := json.Marshal(SkyResp{
		Code: 1,
		Msg:  fmt.Sprintf(message, args...),
		Data: nil,
	})
	return string(resp)
}

// 发生错误
func Failed(_msg uint32, _param string, _data interface{}) string {
	resp, _ := json.Marshal(SkyResp{
		Code: _msg,
		Msg:  _param,
		Data: _data,
	})
	return string(resp)
}

// 发生错误
func JCodeByLang(_langType, _msg uint32, _data interface{}, _param ...interface{}) string {
	resp, _ := json.Marshal(SkyResp{
		Code: _msg,
		Msg:  GenStr(_langType, _msg, _param...),
		Data: _data,
	})
	return string(resp)
}

func JCodeI18n(code uint32, data ...interface{}) string {
	resp, _ := json.Marshal(SkyResp{
		Code: code,
		Data: data,
	})
	return string(resp)
}

// 返回自定义内容
func Diy(_diyContent string) string {
	return _diyContent
}

// 系统内部错误
// 如: 数据库连接不上, redis连接失败, 数据库语法错误等内部代码逻辑错误时返回
func SystemError() string {
	resp, _ := json.Marshal(SkyResp{
		Code: 40001,
		Msg:  "系统内部错误,请联系管理人员查看",
		Data: nil,
	})
	return string(resp)
}

// 玩家一些非法操作引起的错误
// 如: 一些非法提交导致的错误
func ServerError() string {
	resp, _ := json.Marshal(SkyResp{
		Code: 40002,
		Msg:  "服务器繁忙,请稍后再试",
		Data: nil,
	})
	return string(resp)
}

// 需要登录的接口
func NotLogin() string {
	resp, _ := json.Marshal(SkyResp{
		Code: 40000,
		Msg:  "您还未登录或者登录状态已经失效",
		Data: clJson.M{
			"type": 0,
		},
	})
	return string(resp)
}

// 需要登录的接口
func LogoutByKick() string {
	resp, _ := json.Marshal(SkyResp{
		Code: 40000,
		Msg:  "您的账号已经在其他设备登录",
		Data: clJson.M{
			"type": 1,
		},
	})
	return string(resp)
}

// 操作过于频繁
func TooQuickly() string {
	resp, _ := json.Marshal(SkyResp{
		Code: 40003,
		Msg:  "您的操作过于频繁,请稍后再试",
		Data: nil,
	})
	return string(resp)
}

func Download(title string, csv []byte) string {
	resp, _ := json.Marshal(SkyResp{
		Code: 0,
		Msg:  "",
		Data: clJson.M{
			"title":   title,
			"content": csv,
		},
	})
	return string(resp)
}

func ParseDownload(resp string) SkyRespDownload {
	result := SkyRespDownload{}
	json.Unmarshal([]byte(resp), &result)
	return result
}

// 列表
func List(list interface{}, total int32, extra ...clJson.M) string {
	data := clJson.M{
		"list":  list,
		"total": total,
	}
	if len(extra) > 0 {
		for key, val := range extra[0] {
			data[key] = val
		}
	}
	resp, _ := json.Marshal(SkyResp{
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
	return string(resp)
}
