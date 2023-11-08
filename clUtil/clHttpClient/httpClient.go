package clHttpClient

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// client
type ClHttpClient struct {
	url         string                 // 请求地址
	proxy       string                 // 设置代理
	method      string                 // 请求方式
	timeout     uint32                 // 设置超时时间
	encodetype  uint32                 // 加密方式
	aesKey      string                 // AES加密key
	simpleKey   string                 // 一般加密的key
	body        string                 // 请求body
	params      map[string]interface{} // 参数列表
	header      map[string]string      // 请求头
	contentType uint32                 // 请求文档类型
	cert        *CertConfig            // 证书路径
	ua          string                 // 设备类型

}

type CertConfig struct {
	CertFilePath string // 证书文件路径
	KeyFilePath  string // 密钥文件路径
}

const (
	ContentTypeForm    = 0 // 正常form提交
	ContentParam       = 1 // 参数提交(只允许GET)
	ContentJson        = 2 // 通过json提交
	ContentXml         = 3 // 通过xml提交
	ContentXWWWFormUrl = 4 // 通过x-www-from-urlencode方式提交
)

const (
	UAClHttpClient    = "clHttpClient/1.0"
	UASafariMac       = "Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50"
	UASafariWin       = "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50"
	UAIE90            = "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0"
	UAIE80            = "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0)"
	UAIE70            = "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0)"
	UAIE60            = " Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)"
	UAFireFoxMac      = " Mozilla/5.0 (Macintosh; Intel Mac OS X 10.6; rv:2.0.1) Gecko/20100101 Firefox/4.0.1"
	UAFireFoxWin      = "Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1"
	UAOperaMac        = "Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11"
	UAOperaWin        = "Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11"
	UAChromeMac       = " Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11"
	UASougou          = "?Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; SE 2.X MetaSr 1.0; SE 2.X MetaSr 1.0; .NET CLR 2.0.50727; SE 2.X MetaSr 1.0)"
	UA360SE           = " Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)"
	UASafariIOS       = "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5"
	UASafariIpodTouch = "Mozilla/5.0 (iPod; U; CPU iPhone OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5"
	UASafariIpad      = "Mozilla/5.0 (iPad; U; CPU OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5"
	UAAndroidN1       = "Mozilla/5.0 (Linux; U; Android 2.3.7; en-us; Nexus One Build/FRF91) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"
	UAAndroidQQ       = "MQQBrowser/26 Mozilla/5.0 (Linux; U; Android 2.3.7; zh-cn; MB200 Build/GRJ22; CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"
	UAOperaAndroid    = "Opera/9.80 (Android 2.3.4; Linux; Opera Mobi/build-1107180945; U; en-GB) Presto/2.8.149 Version/11.10"
	UABlackBerry      = "Mozilla/5.0 (BlackBerry; U; BlackBerry 9800; en) AppleWebKit/534.1+ (KHTML, like Gecko) Version/6.0.0.337 Mobile Safari/534.1+"
)

// 获取一个新的对象
func NewClient(_url string) *ClHttpClient {

	client := ClHttpClient{
		url:         _url,
		method:      "POST",
		timeout:     30,
		params:      make(map[string]interface{}),
		header:      make(map[string]string),
		cert:        nil,
		body:        "",
		contentType: ContentTypeForm,
	}

	return &client
}

// 设置代理
func (this *ClHttpClient) SetProxy(_proxy string) {
	this.proxy = _proxy
}

// 设置超时时间
func (this *ClHttpClient) SetTimeout(_timeout uint32) {
	this.timeout = _timeout
}

// 设置方式
func (this *ClHttpClient) SetMethod(_method string) {
	this.method = _method
}

// 设置UA
func (this *ClHttpClient) SetUA(_ua string) {
	this.ua = _ua
}

// 设置证书路径
func (this *ClHttpClient) SetCert(_certPath string, _keyPath string) {
	this.cert = &CertConfig{
		CertFilePath: _certPath,
		KeyFilePath:  _keyPath,
	}
}

// 添加参数
func (this *ClHttpClient) AddParam(_key string, _val interface{}) {
	this.params[_key] = _val
}

// 设置请求体
func (this *ClHttpClient) SetBody(_body string) {
	this.body = _body
}

// 设置请求类型
func (this *ClHttpClient) SetContentType(_type uint32) {
	//if _type == ContentJson {
	//	this.method = "POST"
	//}
	this.contentType = _type
}

// 添加头
func (this *ClHttpClient) AddHeader(_key string, _val string) {

	this.header[_key] = _val
}

// 返回最终请求地址
func (this *ClHttpClient) Try() string {
	return this.BuildParamList()
}

// 构建参数
func (this *ClHttpClient) BuildParamList() string {
	if this.contentType != ContentTypeForm && this.contentType != ContentParam {
		return this.url
	}
	// 参数拼接
	paramStr := strings.Builder{}
	for PKey, PVal := range this.params {
		if paramStr.Len() > 0 {
			paramStr.WriteString("&")
		}
		paramStr.WriteString(fmt.Sprintf("%v=%v", PKey, PVal))
	}

	var httpUrl = this.url

	if paramStr.Len() == 0 {
		return httpUrl
	}

	return this.url + "?" + paramStr.String()
}

// 开始请求
func (this *ClHttpClient) Do() (*Response, error) {

	var client *http.Client
	var proxyUrl *url.URL
	var err error

	if this.proxy != "" {

		proxyUrl, err = url.Parse(this.proxy)
		if err != nil {
			return nil, err
		}
		client = &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(this.timeout))
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second * time.Duration(this.timeout)))
					return conn, nil
				},
				ResponseHeaderTimeout: time.Second * time.Duration(this.timeout),
				Proxy:                 http.ProxyURL(proxyUrl),
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			},
		}
	} else {

		tlsConfig := &tls.Config{}
		if this.cert != nil {
			Cacrt, caErr := ioutil.ReadFile(this.cert.CertFilePath)
			if caErr != nil {
				return nil, caErr
			}
			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM(Cacrt)

			cliCrt, keyErr := tls.LoadX509KeyPair(this.cert.CertFilePath, this.cert.KeyFilePath)
			if keyErr != nil {
				return nil, caErr
			}

			tlsConfig.RootCAs = pool
			tlsConfig.InsecureSkipVerify = true
			tlsConfig.Certificates = []tls.Certificate{cliCrt}
		}

		client = &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(this.timeout))
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second * time.Duration(this.timeout)))
					return conn, nil
				},
				ResponseHeaderTimeout: time.Second * time.Duration(this.timeout),
				TLSClientConfig:       tlsConfig,
			},
		}
	}

	var http_url = ""
	var body io.Reader = nil
	if this.body == "" {
		if this.method != "GET" {
			http_url = this.url
			if this.contentType == ContentTypeForm {
				var r = http.Request{}
				r.ParseForm()
				bodyStr := ""
				for key, val := range this.params {
					r.Form.Add(key, fmt.Sprintf("%v", val))
				}
				bodyStr = strings.TrimSpace(r.Form.Encode())
				body = strings.NewReader(bodyStr)
			} else if this.contentType == ContentJson {
				var jsonObj = clJson.M{}
				for key, val := range this.params {
					jsonObj[key] = val
				}
				body = strings.NewReader(clJson.CreateBy(jsonObj).ToStr())
			} else if this.contentType == ContentXml {

				xmlStr := strings.Builder{}
				xmlStr.WriteString("<xml>")
				for k, v := range this.params {
					xmlStr.WriteString(fmt.Sprintf("<%v>%v</%v>", k, v, k))
				}
				xmlStr.WriteString("</xml>")

				body = strings.NewReader(xmlStr.String())
			} else if this.contentType == ContentXWWWFormUrl {
				var urlData = url.Values{}
				for k, v := range this.params {
					urlData.Add(k, fmt.Sprintf("%v", v))
				}
				body = strings.NewReader(urlData.Encode())
			}
		} else {
			http_url = this.BuildParamList()
		}
	} else {
		http_url = this.BuildParamList()
		body = strings.NewReader(this.body)
	}

	req, err := http.NewRequest(this.method, http_url, body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("HttpProxy: %v 请求地址: %v 错误:%v", this.proxy, http_url, err))
	}

	// 添加头
	for HKey, HVal := range this.header {
		req.Header.Add(HKey, HVal)
	}

	if this.contentType == ContentTypeForm {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if this.contentType == ContentJson {
		req.Header.Set("Content-Type", "application/json")
	} else if this.contentType == ContentXml {
		req.Header.Set("Content-Type", "text/xml")
	} else if this.contentType == ContentXWWWFormUrl {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Set("Content-Type", "application/form-data")
	}
	if this.ua == "" {
		this.ua = UA360SE
	}
	req.Header.Set("User-Agent", this.ua)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("HttpProxy: %v 请求地址: %v 错误:%v", this.proxy, http_url, err))
	}
	resp := NewResponse(res)
	res.Body.Close()

	if resp == nil {
		return nil, errors.New(fmt.Sprintf("HttpProxy: %v 请求地址: %v 错误:%v", this.proxy, http_url, err))
	}
	return resp, nil
}

// 获取当前地址
func (this *ClHttpClient) GetUrl() string {
	return this.url
}
