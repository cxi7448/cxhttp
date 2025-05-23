package rule

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clCommon"
	"github.com/cxi7448/cxhttp/clGlobal"
	"github.com/cxi7448/cxhttp/clResponse"
	"github.com/cxi7448/cxhttp/clUtil/clCrypt"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/cxi7448/cxhttp/clUtil/xre"
	"github.com/cxi7448/cxhttp/core/clAuth"
	"github.com/cxi7448/cxhttp/core/clCache"
	"github.com/cxi7448/cxhttp/jwt"
	"github.com/cxi7448/cxhttp/src/skylang"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type ServerParam struct {
	AcName        string
	RemoteIP      string     // 远程IP地址
	RequestURI    string     // 请求URI
	UriData       *HttpParam // uri上的参数列表
	Host          string     // 请求域名
	Method        string     // 请求方法
	Header        http.Header
	RequestURL    string     // 请求完整地址
	UA            string     // 目标设备信息
	UAType        uint32     // 目标设备类型
	Proctol       string     // 目标协议
	Port          string     // 端口
	Language      string     // 使用语言信息
	LangType      uint32     // 使用语言信息
	ContentType   string     // 提交的方式
	RawData       string     // 原始数据
	RawParam      *HttpParam // 原始的参数
	Encrypt       bool       // 是否需要加密
	IsJwt         bool       // 是否开启jwt
	IsForceEncode bool       // 是否开启强制加密模式
	AesKey        string     // 加密用的key
	Iv            string     // 加密用的iv
	Request       *http.Request
	Response      http.ResponseWriter
}

func (this *ServerParam) JwtLogin(_auth *clAuth.AuthInfo) error {
	token, err := jwt.GenToken(_auth)
	if err != nil {
		clLog.Error("生成jwt失败:%v", err)
		return err
	}
	this.Header.Set("Authorization", token) // 自动输入response headers中
	return nil
}

// 重定向
func (this *ServerParam) Location(url string) {
	this.Response.Header().Set("Location", url)
	this.Response.WriteHeader(http.StatusFound)
}

func (this *ServerParam) IsAndroid() bool {
	return xre.IsAndroid(this.Header.Get("User-Agent"))
}

func (this *ServerParam) IsIos() bool {
	return xre.IsIOS(this.Header.Get("User-Agent"))
}

func (this *ServerParam) GetAuth() *clAuth.AuthInfo {
	// 通过jwt处理登陆
	token := this.Header.Get("Authorization")
	if token == "" {
		return nil
	}
	c, err := jwt.ParseToken(token)
	if err != nil {
		return nil
	}
	if !c.IsExpire() {
		return nil
	}

	if !c.IsEffective() {
		return nil
	}

	if c.IsReflush() {
		token, err = c.ReflushToken() // 刷新
		if err != nil {
			clLog.Error("刷新token失败:%v", err)
		} else {
			this.Header.Set("Authorization", token) // 会自动输入response headers中
		}
	}
	return c.GetUser()
}

//@author cxhttp
//@lastUpdate 2019-08-10
//@comment 路由规则定义

// 路由结构  2019-08-10
type Rule struct {
	Request string      // 请求的名字
	Name    string      // 规则名称
	Params  []ParamInfo // 参数列表
	// 回调函数
	CallBack      func(_auth *clAuth.AuthInfo, _param *HttpParam, _server *ServerParam) string
	CacheExpire   int      // 缓存秒数, 负数为频率控制, 正数为缓存时间, 0为不缓存
	CacheType     int      // 缓存启动的时候才有效 0=全局缓存,1=根据IP缓存,2=根据用户缓存
	CacheKeyParam []string // 参与计算唯一性的参数名称列表
	Login         bool     // 是否登录才可以访问这个接口
	Method        string   // 请求方法, 为空则不限制请求方法, POST则为只允许POST请求
	RespContent   string   // 返回的结构体内容格式 默认是 text/json
}

// 路由列表
var ruleList map[string]Rule
var ruleLocker sync.RWMutex

// 请求方式
var requestList map[string]string

func init() {
	ruleList = make(map[string]Rule)
	requestList = make(map[string]string)
}

// 添加新的请求方式
func AddRequest(_request string, _acKey string) {
	requestList[_request] = _acKey
}

// 获取请求方式的ackey
func GetRequestAcKey(_request string) string {
	ackey, exists := requestList[_request]
	if !exists {
		return "ac"
	}
	return ackey
}

// @author cxhttp
// @lastUpdate 2019-08-10
// @comment 添加规则
// @param _rule 规则结构体
func AddRule(_rule Rule) {
	ruleLocker.Lock()
	defer ruleLocker.Unlock()
	// 检测是否重复
	for _, rule := range ruleList {
		if rule.Request == _rule.Request && rule.Name == _rule.Name {
			clLog.Error("路由[%v_%v]重复注册", _rule.Request, _rule.Name)
			os.Exit(1)
		}
	}
	if _rule.Name == "" {
		ruleList[_rule.Request] = _rule
	} else {
		ruleList[_rule.Request+"_"+_rule.Name] = _rule
	}
}

// @auth cxhttp
// @lastUpdate 2021-05-26
// @comment 构建缓存key
func BuildCacheKey(_params []string) string {

	var keys = strings.Join(_params, "_")
	return clCommon.Md5([]byte(keys))
}

// 删除Api缓存
func DelApiCache(_uri string, _acName string, _uid uint64, _params map[string]string) {
	ruleinfo, exists := ruleList[_uri+"_"+_acName]
	if !exists {
		clLog.Error("删除缓存失败! AC <%v_%v> 不存在!", _uri, _acName)
		return
	}
	paramsKeys := make([]string, 0)
	if ruleinfo.CacheType == 2 {
		paramsKeys = append(paramsKeys, fmt.Sprintf("uid=%v", _uid))
	}
	if ruleinfo.Params != nil {
		for _, pinfo := range ruleinfo.Params {
			value := _params[pinfo.Name]
			paramsKeys = append(paramsKeys, pinfo.Name+"="+value)
		}
	}
	cacheKey := BuildCacheKey(paramsKeys)
	clCache.DelCache(cacheKey)
}

// 删除全部Api缓存
func DelApiCacheAll(_uri string, _acName string) {
	clCache.DelCacheContains(_uri + "_" + _acName + "_")
}

func GetRuleInfo(_uri, ac string) *Rule {
	ruleLocker.RLock()
	defer ruleLocker.RUnlock()
	info, ok := ruleList[_uri+"_"+ac]
	if ok {
		return &info
	}
	info, ok = ruleList[_uri]
	if ok {
		return &info
	}
	return nil
}

// @author cxhttp
// @lastUpdate 2019-08-10
// @comment 调用规则
func CallRule(rq *http.Request, rw *http.ResponseWriter, _uri string, _param *HttpParam, _server *ServerParam, ac string) (string, string) {
	var acKey = GetRequestAcKey(_uri)
	// 通过AC获取到指定的路由
	var acName string
	if ac != "" {
		acName = ac
	} else {
		acName = _param.GetStr(acKey, "")
	}
	//ruleinfo, exists := ruleList[_uri+"_"+acName]
	ruleinfo := GetRuleInfo(_uri, acName)
	if ruleinfo == nil {
		if clGlobal.SkyConf.DebugRouter {
			clLog.Error("AC <%v_%v_%v> 不存在! IP: %v", _uri, acName, ac, _server.RemoteIP)
			clLog.Debug("%+v", ruleList)
		}
		respStr := clResponse.JCode(skylang.MSG_ERR_FAILED_INT, "模块不存在!", nil)
		if _server.Encrypt {
			respStr = clCrypt.AesCBCEncode(respStr, _server.AesKey, _server.Iv)
		}
		return respStr, "text/json"
	}

	if ruleinfo.RespContent == "" {
		ruleinfo.RespContent = "text/json"
	}

	var authInfo *clAuth.AuthInfo
	paramsKeys := make([]string, 0)
	paramsKeys = append(paramsKeys, _uri+"_"+acName)

	// 检查参数
	newParam := NewHttpParam(nil, nil)
	if ruleinfo.Params != nil {
		for _, pinfo := range ruleinfo.Params {
			if pinfo.Name == acKey {
				continue
			}
			value := _param.GetStr(pinfo.Name, "")
			if value == PARAM_CHECK_FAIED || value == "" {
				if pinfo.Static {
					// 严格模式
					var msg = "参数:" + pinfo.Name + "不合法!"
					if pinfo.Tips != "" {
						msg = pinfo.Tips
					}
					return clResponse.Error(msg), ruleinfo.RespContent
				} else {
					value = pinfo.Def()
				}
			} else {
				if !pinfo.CheckParam(value) {
					if pinfo.Static {
						// 严格模式
						var msg = "参数:" + pinfo.Name + "不合法!"
						if pinfo.Tips != "" {
							msg = pinfo.Tips
						}
						return clResponse.Error(msg), ruleinfo.RespContent
					} else {
						value = pinfo.Def()
					}
				}
			}
			newParam.Add(pinfo.Name, value)
		}
	} else {
		// 如果路由配置上参数列表为nil，那么就不过滤参数，所有参数都接收
		for key, val := range _param.values {
			newParam.Add(key, val)
		}
	}
	// 判断是否有需要自定义参与计算cacheKey的配置
	if len(ruleinfo.CacheKeyParam) > 0 {
		for _, val := range ruleinfo.CacheKeyParam {
			paramsKeys = append(paramsKeys, _param.GetStr(val, ""))
		}
	} else {
		for _, val := range _param.values {
			paramsKeys = append(paramsKeys, val)
		}
	}
	// 如果回调函数不存在
	if ruleinfo.CallBack == nil {
		if ruleinfo.RespContent != "" {
			return ruleinfo.RespContent, ruleinfo.RespContent
		}

		if clGlobal.SkyConf.DebugRouter {
			clLog.Error("AC[%v]回调函数为空!", acName)
		}
		respStr := clResponse.JCode(skylang.MSG_ERR_FAILED_INT, "模块不存在!", nil)
		if _server.Encrypt {
			respStr = clCrypt.AesCBCEncode(respStr, _server.AesKey, _server.Iv)
		}
		return respStr, "text/json"
	}

	// 需要登录
	if ruleinfo.Login {
		respStr, uInfo := DoAuthCheck(rq, acName, _server, _param)
		if uInfo == nil && respStr == "" {
			respStr = clResponse.NotLogin()
		}
		if respStr != "" {
			if _server.Encrypt {
				respStr = clCrypt.AesCBCEncode(respStr, _server.AesKey, _server.Iv)
			}
			return respStr, ruleinfo.RespContent
		}
		authInfo = uInfo
	}

	// 检查是否需要缓存
	var cacheKey = ""
	if ruleinfo.CacheExpire > 0 {
		// 根据用户缓存
		if ruleinfo.CacheType == 2 {
			paramsKeys = append(paramsKeys, fmt.Sprintf("uid=%v", authInfo.Uid))
		} else if ruleinfo.CacheType == 1 {
			// 根据IP缓存
			paramsKeys = append(paramsKeys, "ip="+_server.RemoteIP)
		}
		cacheKey = rq.RequestURI + "_" + acName + "_" + BuildCacheKey(paramsKeys)
		if _server.Encrypt { // 如果是加密的话，需要带上Iv
			cacheKey += _server.Iv
		}
		cacheStr := clCache.GetCache(cacheKey)
		if cacheStr != "" {
			return cacheStr, ruleinfo.RespContent
		}
	} else if ruleinfo.CacheExpire < 0 {
		// 根据用户缓存
		if ruleinfo.CacheType == 2 {
			if authInfo != nil {
				paramsKeys = append(paramsKeys, fmt.Sprintf("uid=%v", authInfo.Uid))
			}
		} else if ruleinfo.CacheType == 1 {
			// 根据IP缓存
			paramsKeys = append(paramsKeys, "ip="+_server.RemoteIP)
		}
		cacheKey = rq.RequestURI + "_" + acName + "_" + BuildCacheKey(paramsKeys) + "_NX"
		if _server.Encrypt { // 如果是加密的话，需要带上Iv
			cacheKey += _server.Iv
		}
		if !clCache.SetNX(cacheKey, uint32(-ruleinfo.CacheExpire)) {
			return clResponse.Failed(2000, "您操作的太频繁了!", nil), "text/json"
		}
	}

	// 调用前置函数，并返回结果
	var beforeParam = DoRequestBefore(_uri, &RequestBeforeParam{
		Request:    rq,
		AcName:     acName,
		ServerInfo: _server,
		UserInfo:   authInfo,
		Param:      _param,
		Rule:       ruleinfo,
	})
	if beforeParam.RejectResp != "" {
		return beforeParam.RejectResp, ruleinfo.RespContent
	}

	nowTime := time.Now()
	respStr := ruleinfo.CallBack(beforeParam.UserInfo, beforeParam.Param, beforeParam.ServerInfo)
	diffTime := time.Since(nowTime).Seconds()
	if diffTime > 5 {
		clLog.Error("接口:%v.%v 处理耗时(%0.2fs)过长!", _uri, acName, diffTime)
	}

	// 后处理函数
	afterResp := DoRequestAfter(_uri, &RequestAfterParam{
		Request:        rq,
		AcName:         acName,
		ServerInfo:     _server,
		UserInfo:       authInfo,
		Param:          _param,
		ResponseText:   respStr,
		ResponseWriter: rw,
	})

	if _server.IsJwt {
		// 处理jwt的Authorization
		authorization := rq.Header.Get("Authorization")
		if authorization != "" {
			_rw := *rw
			_rw.Header().Set("Authorization", authorization)
		}
	}

	respStr = afterResp.ResponseText

	// 需要加密
	if _server.Encrypt {
		respStr = clCrypt.AesCBCEncode(respStr, _server.AesKey, _server.Iv)
	}

	// 检查是否需要缓存
	if ruleinfo.CacheExpire > 0 {
		clCache.UpdateCacheSimple(cacheKey, respStr, uint32(ruleinfo.CacheExpire))
	}

	if clGlobal.SkyConf.DebugRouter {
		clLog.Debug("[%s][%s] REQUEST: %s / RESPONSE: %s", acName, _server.RemoteIP, strings.Join(paramsKeys, "&"), respStr)
	}

	return respStr, ruleinfo.RespContent
}
