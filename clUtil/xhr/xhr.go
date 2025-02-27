package xhr

import (
	"encoding/json"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type XHR struct {
	URL         string
	Method      string
	Body        string
	ContentType string
	Debug       bool
	data        interface{}
	headers     clJson.M
}

const ContentTypeJSON = "application/json"
const ContentTypeFORM = "application/x-www-form-urlencoded"

func NewXhr(rUrl string) *XHR {
	return &XHR{
		URL:         rUrl,
		Method:      "POST",
		ContentType: "",
	}
}

func (this *XHR) SetContentType(content_type string) *XHR {
	this.ContentType = content_type
	return this
}
func (this *XHR) SetDebug(debug bool) *XHR {
	this.Debug = debug
	return this
}
func (this *XHR) SetHeaders(headers clJson.M) *XHR {
	this.headers = headers
	return this
}

func (this *XHR) SetJSON() *XHR {
	this.SetContentType(ContentTypeJSON)
	return this
}
func (this *XHR) SetFORM() *XHR {
	this.SetContentType(ContentTypeFORM)
	return this
}

func (this *XHR) PostBody(data interface{}, value interface{}) error {
	this.SetJSON()
	var reqBody *strings.Reader
	var printData = ""
	this.data = data
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}
	reqBody = strings.NewReader(string(content))
	printData = string(content)
	if this.Debug {
		clLog.Info("POST:%v [%v] [%v]", this.URL, this.ContentType, printData)
	}
	request, err := http.NewRequest("POST", this.URL, reqBody)
	if err != nil {
		return err
	}
	request.Header.Set("content-type", this.ContentType)
	if len(this.headers) > 0 {
		for key, val := range this.headers {
			request.Header.Add(key, fmt.Sprint(val))
		}
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		clLog.Error("ioutil.ReadAll失败:[%v]", err)
		return err
	}
	this.Body = string(body)
	err = json.Unmarshal(body, value)
	if err != nil {
		clLog.Error("Unmarshal失败:[%v]", err)
		clLog.Error("Unmarshal失败:[%v]", string(body))
		return err
	}
	return nil
}

func (this *XHR) Post(data clJson.M, value interface{}) error {
	if this.ContentType == "" {
		this.SetFORM()
	}
	this.data = data
	var reqBody *strings.Reader
	var printData = ""
	if this.ContentType == ContentTypeJSON {
		content, err := json.Marshal(data)
		if err != nil {
			return err
		}
		reqBody = strings.NewReader(string(content))
		printData = string(content)
	} else {
		var form = url.Values{}
		for key, val := range data {
			form.Add(key, fmt.Sprint(val))
		}
		printData = form.Encode()
		reqBody = strings.NewReader(form.Encode())
	}
	if this.Debug {
		clLog.Info("POST:%v [%v] [%v]", this.URL, this.ContentType, printData)
	}
	request, err := http.NewRequest("POST", this.URL, reqBody)
	if err != nil {
		return err
	}
	request.Header.Set("content-type", this.ContentType)
	if len(this.headers) > 0 {
		for key, val := range this.headers {
			request.Header.Add(key, fmt.Sprint(val))
		}
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		clLog.Error("ioutil.ReadAll失败:[%v]", err)
		return err
	}
	this.Body = string(body)
	err = json.Unmarshal(body, value)
	if err != nil {
		clLog.Error("Unmarshal失败:[%v]", err)
		clLog.Error("Unmarshal失败:[%v]", string(body))
		return err
	}
	return nil
}

func (this *XHR) GetBody(_data ...clJson.M) ([]byte, error) {
	var reqUrl = this.URL
	data := clJson.M{}
	if len(_data) > 0 && len(_data[0]) > 0 {
		data = _data[0]
		var form = url.Values{}
		for key, val := range data {
			form.Add(key, fmt.Sprint(val))
		}
		if strings.Contains(reqUrl, "?") {
			reqUrl += form.Encode()
		} else {
			reqUrl += "?" + form.Encode()
		}
	}
	if this.Debug {
		clLog.Info("GET:%v", reqUrl)
	}
	request, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	if this.ContentType != "" {
		request.Header.Set("content-type", this.ContentType)
	}
	if len(this.headers) > 0 {
		for key, val := range this.headers {
			request.Header.Add(key, fmt.Sprint(val))
		}
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp != nil && resp.Body != nil {
			body, _ := ioutil.ReadAll(resp.Body)
			clLog.Error("错误:%v", string(body))
		}
		return nil, fmt.Errorf("错误:%v", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		clLog.Error("ioutil.ReadAll失败:[%v]", err)
		return nil, err
	}
	return body, nil
}

func (this *XHR) Get(value interface{}, _data ...clJson.M) error {
	var reqUrl = this.URL
	data := clJson.M{}
	if len(_data) > 0 && len(_data[0]) > 0 {
		data = _data[0]
		var form = url.Values{}
		for key, val := range data {
			form.Add(key, fmt.Sprint(val))
		}
		if strings.Contains(reqUrl, "?") {
			reqUrl += form.Encode()
		} else {
			reqUrl += "?" + form.Encode()
		}
	}
	if this.Debug {
		clLog.Info("GET:%v", reqUrl)
	}
	request, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return err
	}
	if this.ContentType != "" {
		request.Header.Set("content-type", this.ContentType)
	}
	if len(this.headers) > 0 {
		for key, val := range this.headers {
			request.Header.Add(key, fmt.Sprint(val))
		}
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp != nil && resp.Body != nil {
			body, _ := ioutil.ReadAll(resp.Body)
			clLog.Error("错误:%v", string(body))
		}
		return fmt.Errorf("错误:%v", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		clLog.Error("ioutil.ReadAll失败:[%v]", err)
		return err
	}
	this.Body = string(body)
	err = json.Unmarshal(body, value)
	if err != nil {
		clLog.Error("Unmarshal失败:[%v]", err)
		clLog.Error("Unmarshal失败:[%v]", string(body))
		return err
	}
	return nil
}
