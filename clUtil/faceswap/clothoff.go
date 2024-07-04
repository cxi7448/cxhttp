package faceswap

import (
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/cxi7448/cxhttp/clUtil/xhttp"
)

// 脱衣
type ClothOff struct {
	Api
	api_key string
	webhook string
}

var clothOff = &ClothOff{}

func InitClothOff(api_key, webhook string) {
	clothOff.api_key = api_key
	clothOff.webhook = webhook
}
func (this *ClothOff) Undress(filename, order string) error {
	xhr := xhttp.New("https://public-api.clothoff.io/undress")
	param := map[string]string{
		"type_gen":    "img2clo",
		"id_gen":      order,
		"webhook":     this.webhook,
		"age":         "20",
		"breast_size": "normal",
		"body_type":   "normal",
		"butt_size":   "normal",
	}
	result := clJson.M{}
	xhr.SetHeaders(map[string]string{
		"x-api-key": this.api_key,
		"accept":    "application/json",
	})
	err := xhr.PostForm("image", filename, param, &result)
	if err != nil {
		return err
	}
	clLog.Info("提交结果:%+v", result)
	return nil
}
