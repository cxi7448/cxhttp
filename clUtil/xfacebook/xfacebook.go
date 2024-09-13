package xfacebook

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clCommon"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	xhr2 "github.com/cxi7448/cxhttp/clUtil/xhr"
	"strings"
	"time"
)

const API_VERSION = "v20.0"
const HOST = "https://graph.facebook.com"

type FaceBook struct {
	AccessToken   string
	TestEventCode string
	PIXEL_ID      string
	ip            string
	email         string
	phone         string
	user_agent    string
	website       string
	country       string
	user_data     clJson.M
	custom_data   clJson.M
	event         string
}

func NewWith(PIXEL_ID, AccessToken, TestEventCode string) *FaceBook {
	return &FaceBook{AccessToken: AccessToken, PIXEL_ID: PIXEL_ID, TestEventCode: TestEventCode, user_data: clJson.M{}, custom_data: clJson.M{}, website: "website"}
}

func (this *FaceBook) SetTestEventCode(test_even_code string) *FaceBook {
	this.TestEventCode = test_even_code
	return this
}

func (this *FaceBook) SetWebsite(website string) *FaceBook {
	this.website = website
	return this
}

func (this *FaceBook) SetIP(value string) *FaceBook {
	this.ip = value
	return this
}

func (this *FaceBook) SetEmail(value string) *FaceBook {
	this.email = value
	return this
}

func (this *FaceBook) SetPhone(value string) *FaceBook {
	this.phone = value
	return this
}

func (this *FaceBook) SetUserAgent(value string) *FaceBook {
	this.user_agent = value
	return this
}

func (this *FaceBook) SetCountry(value string) *FaceBook {
	this.country = value
	return this
}

func (this *FaceBook) SetUserData(user_data clJson.M) *FaceBook {
	for key, val := range user_data {
		this.user_data[key] = val
	}
	return this
}

func (this *FaceBook) SetCustomData(data clJson.M) *FaceBook {
	return this
}

// CompleteRegistration 注册
// Purchase // 首充
// Search 搜索
// ViewContent 访问
func (this *FaceBook) Purchase(currency string, price float64) error {
	this.event = "Purchase"
	this.custom_data["currency"] = currency
	this.custom_data["value"] = fmt.Sprint(price)
	return this.Trigger()
}

func (this *FaceBook) Register() error {
	this.event = "CompleteRegistration"
	return this.Trigger()
}

func (this *FaceBook) Trigger() error {
	xhr := xhr2.NewXhr(fmt.Sprintf("%v/%v/%v/events?access_token=%v", HOST, API_VERSION, this.PIXEL_ID, this.AccessToken))
	if this.email != "" {
		this.user_data["em"] = []string{clCommon.Sha256(this.email)}
	}
	if this.phone != "" {
		this.user_data["ph"] = []string{clCommon.Sha256(this.phone)}
	}
	if this.country != "" {
		this.user_data["country"] = clCommon.Sha256(strings.ToLower(this.country))
	}
	if this.ip != "" {
		this.user_data["client_ip_address"] = this.ip
	}
	if this.user_agent != "" {
		this.user_data["client_user_agent"] = this.user_agent
	}
	data := clJson.M{
		"data": clJson.A{
			clJson.M{
				"event_name":    this.event,
				"event_time":    time.Now().Unix(),
				"action_source": this.website,
				"user_data":     this.user_data,
				"custom_data":   this.custom_data,
			},
		},
		"test_event_code": this.TestEventCode,
	}
	result := clJson.M{}
	xhr.SetJSON()
	err := xhr.Post(data, &result)
	if err != nil {
		return err
	}
	if result.Has("error") {
		getMap := result.GetMap("error")
		return fmt.Errorf(getMap.Get("error_user_msg"))
	}
	return nil
}
