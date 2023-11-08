package example

import (
	"github.com/cxi7448/cxhttp/clResponse"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/core/clAuth"
	"github.com/cxi7448/cxhttp/core/rule"
)

func ApiExample(_auth *clAuth.AuthInfo, _param *rule.HttpParam, _server *rule.ServerParam) string {

	//// 获取字符串列表数组
	//strArr := _param.GetStrSplit("str_list", ",")
	//
	//// 获取整数列表数组
	//numArr := _param.GetUint32Split("id_list")
	//
	//// 获取浮点数列表数组
	//posArr := _param.GetFloatSplit("pos_list")

	return clResponse.Success(_server)
}

// 上传范例
func ApiUploadExample(_auth *clAuth.AuthInfo, _param *rule.HttpParam, _server *rule.ServerParam) string {

	return clResponse.Success(clJson.M{
		"filename":  _param.GetStr("filename", ""),
		"fileExt":   _param.GetStr("fileExt", ""),
		"localPath": _param.GetStr("localPath", ""),
	})
}
