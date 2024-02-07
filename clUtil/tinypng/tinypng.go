package tinypng

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clFile"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	CompressingUrl = "https://api.tinify.com/shrink"
)

var (
	//APIKEY = "pGbnN0m3CbkfBmnyMn1QzM44FBXYr0rk"
	//EMAIL  = "cxi7448@gmail.com"
	//APIKEY = "qFgqsbqzb5WN69XgW70FDM5DyKg9zjRD"
	APIKEY = "Ln7mMtQxYwsMWhMDzGtSSM9vm79sNwXQ"
	EMAIL  = "64839198aa@gmail.com"
)

type Tinypng struct {
	Input  Input  `json:"input"`
	Output Output `json:"output"`
}

type Output struct {
	Size   int64   `json:"size"`
	Type   string  `json:"type"`
	Width  int64   `json:"width"`
	Height int64   `json:"height"`
	Ratio  float64 `json:"ratio"`
	Url    string  `json:"url"`
}

type Input struct {
	Size int64  `json:"size"`
	Type string `json:"type"`
}

const (
	ROOT_DIR     = "./tmp/"
	COMPRESS_DIR = "compress/"
)

func init() {
	os.MkdirAll(ROOT_DIR+COMPRESS_DIR, 0777)
}

func Init(email, apikey string) {
	EMAIL = email
	APIKEY = apikey
}
func Compress(image_path string) (string, error) {
	ext := image_path[strings.LastIndex(image_path, "."):]
	new_path := fmt.Sprintf("%v%v%v_compress%v", ROOT_DIR, COMPRESS_DIR, time.Now().UnixNano(), ext)
	// 创建Request
	req, err := http.NewRequest(http.MethodPost, CompressingUrl, nil)
	if err != nil {
		clLog.Error("错误:%v", err)
		return "", err
	}
	// 将鉴权信息写入Request
	req.SetBasicAuth(EMAIL, APIKEY)
	// 将图片以二进制的形式写入Request
	data, err := ioutil.ReadFile(image_path)
	if err != nil {
		clLog.Error("错误:%v", err)
		return "", err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(data))
	// 发起请求
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		clLog.Error("错误:%v", err)
		return "", err
	}

	// 解析请求
	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		clLog.Error("解析Tinypng的返回失败:%v", err)
		return "", err
	}
	result := Tinypng{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println(string(data))
		clLog.Error("解析Tinypng的返回失败:%v", err)
		return "", err
	}
	if result.Output.Url == "" {
		fmt.Println(string(data))
		return "", fmt.Errorf("拿不到图片压缩地址")
	}
	err = clFile.Download(result.Output.Url, new_path)
	if err != nil {
		clLog.Error("下载Tinypng的压缩图片地址失败:%v", err)
	}
	return new_path, err
}
