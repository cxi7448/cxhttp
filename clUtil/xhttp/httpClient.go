package xhttp

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpClient struct {
	URL       string
	Method    string
	Param     interface{}
	Headers   map[string]string
	request   *http.Request
	respone   *http.Response
	cookies   clJson.M
	cookieSet func(this *HttpClient, cookie clJson.M)
	cookieGet func(this *HttpClient) clJson.M
	skipHttps bool // 是否跳过证书检测
}

func New(url string) *HttpClient {
	httpClient := &HttpClient{
		URL:    url,
		Method: "GET",
		Param:  map[string]interface{}{},
		Headers: map[string]string{
			"Content-Type": "application-json",
		},
		cookies: clJson.M{},
	}
	return httpClient
}

func (this *HttpClient) SetSkipHttps(skipHttps bool) *HttpClient {
	this.skipHttps = skipHttps
	return this
}

func (this *HttpClient) SetHeaders(headers map[string]string) *HttpClient {
	for key, value := range headers {
		this.Headers[key] = value
	}
	return this
}

func (this *HttpClient) Get(result interface{}) error {
	this.SetMethod("GET")
	return this.do(result)
}

func (this *HttpClient) Post(data interface{}, result interface{}) error {
	this.SetMethod("POST")
	this.Param = data
	return this.do(result)
}
func (this *HttpClient) do(result interface{}) error {
	var reqBody *strings.Reader
	if this.Param != nil {
		bParam, err := json.Marshal(this.Param)
		if err != nil {
			return err
		}
		reqBody = strings.NewReader(string(bParam))
	}

	request, err := http.NewRequest(this.Method, this.URL, reqBody)
	if err != nil {
		return err
	}
	this.request = request
	if len(this.Headers) > 0 {
		for key, value := range this.Headers {
			this.request.Header.Set(key, value)
		}
	}

	// 自动读取cookie信息
	if this.cookieGet != nil {
		this.SetCookies(this.cookieGet(this))
	}

	if len(this.cookies) > 0 {
		for name, value := range this.cookies {
			cookie := &http.Cookie{
				Name:  name,
				Value: fmt.Sprint(value),
			}
			this.request.AddCookie(cookie)
		}
	}
	client := http.Client{}
	if this.skipHttps {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	resp, err := client.Do(this.request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		clLog.Error("ioutil.ReadAll失败:[%v]", err)
		return err
	}
	this.respone = resp
	this.saveCookie()
	err = json.Unmarshal(body, result)
	if err != nil {
		clLog.Error("Unmarshal失败:[%v]", err)
		clLog.Error("Unmarshal失败:[%v]", string(body))
		return err
	}
	return nil
}

func (this *HttpClient) saveCookie() {
	// 自动存储cookie 下次访问的时候，自动带上cookie
	for _, jar := range this.request.Cookies() {
		this.cookies[jar.Name] = jar.Value
	}

	for _, jar := range this.respone.Cookies() {
		this.cookies[jar.Name] = jar.Value
	}
	// 自动存储cookie
	if this.cookieSet != nil {
		this.cookieSet(this, this.cookies)
	}
}

func (this *HttpClient) SetMethod(method string) *HttpClient {
	this.Method = method
	return this
}

func (this *HttpClient) SetCookies(cookies clJson.M) *HttpClient {
	for name, value := range cookies {
		this.cookies[name] = value
	}
	return this
}

func (this *HttpClient) GetCookie() clJson.M {
	return this.cookies
}

// 自定义cookie函数  只对同一个实例有效
func (this *HttpClient) CookieFunc(set func(this *HttpClient, cookies clJson.M), get func(this *HttpClient) clJson.M) clJson.M {
	this.cookieSet = set
	this.cookieGet = get
	return this.cookies
}
