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
}

func NewWith(PIXEL_ID, AccessToken, TestEventCode string) *FaceBook {
	return &FaceBook{AccessToken: AccessToken, PIXEL_ID: PIXEL_ID, TestEventCode: TestEventCode, user_data: clJson.M{}}
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

// CompleteRegistration 注册
// Purchase // 首充
// Search 搜索
// ViewContent 访问
func (this *FaceBook) Purchase(currency string, price float64) error {
	// map[events_received:1 fbtrace_id:A8J_fLIX0L4FvOULvM6X_y5 messages:[]]
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
				"event_name":    "Purchase",
				"event_time":    time.Now().Unix(),
				"action_source": this.website,
				"user_data":     this.user_data,
				"custom_data": clJson.M{
					"currency": currency,
					"value":    fmt.Sprint(price),
				},
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
	fmt.Printf("%+v \n", result)
	fmt.Println(xhr.Body)
	if result.Has("error") {
		getMap := result.GetMap("error")
		return fmt.Errorf(getMap.Get("error_user_msg"))
	}
	return nil
}

//					"fn":                  "名字", // 必须进行哈希处理。 推荐使用罗马字母 a-z 字符。仅限小写字母，且不可包含标点符号。若使用特殊字符，则须按 UTF-8 格式对文本进行编码
//					"ln":                  "姓",  // 必须进行哈希处理。 推荐使用罗马字母 a-z 字符。仅限小写字母，且不可包含标点符号。若使用特殊字符，则须按 UTF-8 格式对文本进行编码。
//					"db":                  "生日", // 必须进行哈希处理。 我们接受 YYYYMMDD 格式，其中涵盖各类月、日、年组合，带不带标点均可。
//					"ge":                  "性别", //必须进行哈希处理。我们接受以小写首字母表示性别的做法。示例：f 表示女性 m 表示男性
//					"ct":                  "城市", // 必须进行哈希处理。 推荐使用罗马字母字符 a 至 z。仅限小写字母，且不可包含标点符号、特殊字符和空格。若使用特殊字符，则须按 UTF-8 格式对文本进行编码。
//					"st":                  "",   // 必须进行哈希处理。使用 2 个字符的 ANSI 缩写代码，必须为小写字母。请使用小写字母对美国境外的州/省/自治区/直辖市名称作标准化处理，且不可包含标点符号、特殊字符和空格。
//					"external_id":         "",   //必须进行哈希处理。可以是广告主提供的任何唯一编号，如会员编号、用户编号和外部 Cookie 编号。您可以为给定事件发送一或多个外部编号。 如果是通过其他渠道发送外部编号，此编号的格式应与通过转化 API 发送时的格式相同。
//					"fbc":                 "",   // 点击编号
//					"fbp":                 "",   // 浏览器编号
//					"subscription_id":     "订阅编号",
//					"fb_login_id":         0, //"Facebook 登录编号",
//					"lead_id":             1, //"线索编号\n\n",
//					"anon_id":             "",
//					"madid":               "", // 您的移动广告客户编号、Android 设备中的广告编号或 Apple 设备中的广告 ID (IDFA)。
//					"page_id":             "", // 您的公共主页编号。指定与事件关联的公共主页编号。使用与智能助手关联的公共主页的 Facebook 公共主页编号。
//					"page_scoped_user_id": "", //指定与记录事件的 Messenger 智能助手关联的公共主页范围用户编号。使用提供给 Webhooks 的公共主页范围用户编号
//					"ctwa_clid":           "", // 点击 Meta 为 WhatsApp 直达广告生成的编号。
//					"ig_account_id":       "", //与商家关联的 Instagram 账户编号
//					"ig_sid":              "", // 根据 Instagram 范围用户编号 (IGSID) 识别与 Instagram 互动的用户。可以从此 Webhooks 获取 IGSID。
