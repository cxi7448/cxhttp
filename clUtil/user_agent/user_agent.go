package user_agent

import (
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"regexp"
	"strings"
)

const (
	UNKNOWN    = "unknown"
	OS_ANDROID = "Android"
	OS_IPHONE  = "iPhone"
	OS_IPAD    = "iPad"
	OS_WP      = "window phone"
)

type UserAgent struct {
	UA            string
	Model         string // 手机型号
	Version       string // 版本
	OS            string // 系统
	isAndroid     *bool  // 是否是安卓系统
	isIPhone      *bool  // 是否是苹果系统
	isIPad        *bool  // ipad
	isWP          *bool  // window系统
	isBuildHuawei *bool  // 是否是Build/HUAWEI
}

func New(user_agent string) *UserAgent {
	info := &UserAgent{
		UA: user_agent,
	}
	info.parse()
	clLog.Info("Model:[%v] Version:[%v] OS: [%v]", info.Model, info.Version, info.OS)
	return info
}

func (this *UserAgent) parse() {
	this.OS = UNKNOWN
	this.Version = UNKNOWN
	this.Model = UNKNOWN
	if this.IsAndroid() {
		this.OS = OS_ANDROID
		version := regexp.MustCompile(`Android[^;]+;`).FindString(this.UA)
		if version != "" {
			this.Version = version[0 : len(version)-1]
		}
		bindex := strings.LastIndex(this.UA, "Build/")
		if bindex > -1 {
			str_tmp := this.UA[0:bindex]
			this.Model = str_tmp[strings.LastIndex(str_tmp, ";")+1 : bindex]
		}
		if this.IsBuildHuawei() {
			res := regexp.MustCompile(`(?i)build/huawei\S+\)`).FindString(this.UA)
			if res != "" {
				this.Model = res[12 : len(res)-1]
			}
		}
	} else {
		if this.IsIphone() {
			this.OS = OS_IPHONE
		} else if this.IsIPad() {
			this.OS = OS_IPAD
		} else if this.IsWP() {
			this.OS = OS_WP
		}
		//if(u.indexOf("iPhone") > -1){
		//	dev_info.dev_os = "iOS"
		//	const version = u.match(/iPhone OS .*?(?= )/)
		//	console.log(version)
		//	if(version && version.length > 0){
		//		dev_info.dev_os_ver = version[0]
		//	}
		//	dev_info.dev_os_model = "iPhone"
		//}else if (u.indexOf("iPad") > -1) {
		//	dev_info.dev_os = "iOS"
		//	const version = u.match(/CPU OS .*?(?= )/)
		//	if(version && version.length > 0){
		//		dev_info.dev_os_ver = version[0]
		//	}
		//	dev_info.dev_os_model = "iPad"
		//} else if (u.indexOf("Windows Phone") > -1){
		//	dev_info.dev_os = "WP"
		//}
	}
	this.Model = strings.TrimSpace(this.Model)
	this.Version = strings.TrimSpace(this.Version)
}

func (this *UserAgent) IsBuildHuawei() bool {
	if this.isBuildHuawei != nil {
		return *this.isBuildHuawei
	}
	isBuildHuawei := regexp.MustCompile(`(?i)build/huawei\S+\s?`).MatchString(this.UA)
	this.isBuildHuawei = &isBuildHuawei
	return *this.isBuildHuawei
}

func (this *UserAgent) IsWP() bool {
	if this.isWP != nil {
		return *this.isWP
	}
	isWP := regexp.MustCompile(`Windows Phone`).MatchString(this.UA)
	this.isWP = &isWP
	return *this.isWP
}
func (this *UserAgent) IsIPad() bool {
	if this.isIPad != nil {
		return *this.isIPad
	}
	isIPad := regexp.MustCompile(`iPad`).MatchString(this.UA)
	this.isIPad = &isIPad
	return *this.isIPad
}

func (this *UserAgent) IsAndroid() bool {
	if this.isAndroid != nil {
		return *this.isAndroid
	}
	isAndroid := regexp.MustCompile(`(Android|Linux)`).MatchString(this.UA)
	this.isAndroid = &isAndroid
	return *this.isAndroid
}

func (this *UserAgent) IsIphone() bool {
	if this.isIPhone != nil {
		return *this.isIPhone
	}
	isIPhone := regexp.MustCompile(`iPhone`).MatchString(this.UA)
	this.isIPhone = &isIPhone
	return *this.isIPhone
}
