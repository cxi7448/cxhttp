package email

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/cxi7448/cxhttp/clUtil/xhr"
)

type NxClound struct {
	appid    string
	secret   string
	host     string
	template string // 短信模板
}

func NewNxClound(appid, secret string) *NxClound {
	var sNxClount = &NxClound{
		host:     "http://api2.nxcloud.com",
		template: "默认模板",
	}
	sNxClount.appid = appid
	sNxClount.secret = secret
	return sNxClount
}

func (this *NxClound) SetTemplate(template string) *NxClound {
	this.template = template
	return this
}

func (this *NxClound) Send(from, to string, templateData clJson.M) error {
	rUrl := fmt.Sprintf("%v/api/email/otp", this.host)
	req := xhr.NewXhr(rUrl)
	req.SetJSON()
	result := clJson.M{}
	data := clJson.M{
		"appKey":       this.appid,
		"secretKey":    this.secret,
		"from":         from,
		"to":           to,
		"templateName": this.template,
		"templateData": templateData,
	}
	err := req.Post(data, &result)
	if err != nil {
		clLog.Error("邮件发送错误:%v", err)
		return err
	}
	fmt.Printf("%+v\n", result)
	if result.Uint32("code") != 0 {
		return fmt.Errorf(result.Get("msg"))
	}
	return nil
}
