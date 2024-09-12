package faceswap

import (
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/clUtil/xhttp"
)

// 脱衣
type ClothOff struct {
	Api
	api_key string
	webhook string
	api_url string
}

var clothOff = &ClothOff{}

func InitClothOff(api_key, webhook, api_url string) {
	clothOff.api_key = api_key
	clothOff.webhook = webhook
	clothOff.api_url = api_url
	if clothOff.api_url == "" {
		clothOff.api_url = "https://public-api.clothoff.io/undress"
	}
}
func (this *ClothOff) Undress(filename, order string) error {
	xhr := xhttp.New(clothOff.api_url)
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
	return nil
}
