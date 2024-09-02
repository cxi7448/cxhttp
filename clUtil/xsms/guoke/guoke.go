package guoke

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/cxi7448/cxhttp/clUtil/xhr"
	"net/url"
)

type GuoKe struct {
	Appkey    string
	SecretKey string
	Host      string
}

func NewGuoKe(appkey, secret string) *GuoKe {
	return &GuoKe{
		Appkey:    appkey,
		SecretKey: secret,
		Host:      "http://api.wftqm.com",
	}
}

func (this *GuoKe) Send(phone, content string) error {
	rUrl := fmt.Sprintf("%v/api/sms/mtsend", this.Host)
	req := xhr.NewXhr(rUrl)
	result := clJson.M{}
	data := clJson.M{
		"appkey":    this.Appkey,
		"secretkey": this.SecretKey,
		"phone":     phone,
		"content":   url.QueryEscape(content),
	}
	err := req.Post(data, &result)
	if err != nil {
		clLog.Error("短信发送错误:%v", err)
		return err
	}
	if result.Uint32("code") != 0 {
		return fmt.Errorf(result.Get("result"))
	}
	return nil
}
