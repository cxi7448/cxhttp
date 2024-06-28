package faceswap

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
)

/*
*
src: 原图
face: 脸图
_type: API分类
*/
func Faceswap(src, face Img, _type ...string) (string, error) {
	clLog.Info("收到转换信息")
	clLog.Info("src:%+v", src)
	clLog.Info("face:%+v", face)
	api_type := TYPE_AKOOL
	if len(_type) > 0 {
		api_type = _type[0]
	}
	switch api_type {
	case TYPE_AKOOL:
		return akool.FaceSwap(src, face)
	default:
		return "", fmt.Errorf("未知API")
	}
}

/*
*
生成脸部数据opts
*/
func GenDetect(image string, _type ...string) (string, error) {
	clLog.Info("图片信息:%v", image)
	api_type := TYPE_AKOOL
	if len(_type) > 0 {
		api_type = _type[0]
	}
	switch api_type {
	case TYPE_AKOOL:
		return akool.Detect(image)
	default:
		return "", fmt.Errorf("未知API")
	}
}
