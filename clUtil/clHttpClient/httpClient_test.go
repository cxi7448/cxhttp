package clHttpClient

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {

	var httpClient = NewClient("http://www.baidu.com")
	httpClient.SetUA( UA360SE )
	httpClient.SetMethod("POST")
	httpClient.SetTimeout(10)
	httpClient.SetContentType(ContentParam)
	httpClient.SetProxy("http://www.qq.com")
	httpClient.AddHeader("referrer", "http://www.baidu.com")
	httpClient.AddParam("key1", "value1")
	httpClient.SetBody("hello world")
	resp, err := httpClient.Do()
	if err != nil {
		fmt.Printf("请求错误: %v\n", err)
		return
	}

	fmt.Printf("请求状态码: %v\n", resp.StatusCode)
	fmt.Printf("请求回应体: %v\n", resp.Body)
	fmt.Printf("请求体长度: %v\n", resp.ContentLength)
	fmt.Printf("请求头部: %+v\n", resp.Header)

	var respObj = StructResp{}
	resp.ToObject(&respObj)
}

type StructResp struct {
	Code uint32 `json:"code"`
	Msg string `json:"msg"`
	Data string `json:"data"`
}