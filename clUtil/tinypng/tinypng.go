package tinypng

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cxi7448/cxhttp/clCommon"
	"github.com/cxi7448/cxhttp/clUtil/clFile"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
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
	APIKEY = "Ln7mMtQxYwsMWhMDzGtSSM9vm79sNwXQ"
	EMAIL  = "64839198aa@gmail.com"
)

type Tinypng struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Input   Input  `json:"input"`
	Output  Output `json:"output"`
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
	//if try_times <= 0 {
	//	return "", fmt.Errorf("tinypng找不到合适的账号压缩,请稍后再试!")
	//}
	// 将图片以二进制的形式写入Request
switchAccount:
	data, err := ioutil.ReadFile(image_path)
	if err != nil {
		clLog.Error("错误:%v", err)
		return "", err
	}
	ext := image_path[strings.LastIndex(image_path, "."):]
	new_path := fmt.Sprintf("%v%v%v_%v_compress%v", ROOT_DIR, COMPRESS_DIR, time.Now().UnixNano(), clCommon.Md5([]byte(image_path)), ext)
	// 创建Request
	req, err := http.NewRequest(http.MethodPost, CompressingUrl, nil)
	if err != nil {
		clLog.Error("错误:%v", err)
		return "", err
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(data))
	info := GetOne()
	if info == nil {
		return "", fmt.Errorf("tinypng找不到合适的账号压缩,请稍后再试!!")
	}
	//clLog.Info("拿到账号:%v-%v", info.Email, info.Apikey)
	// 将鉴权信息写入Request
	req.SetBasicAuth(info.Email, info.Apikey)
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
	if result.Error != "" {
		// 错误信息 异步处理报错的信息
		clLog.Error("切换账号[%v]-[%v]错误信息:%v", info.Email, info.Apikey, string(data))
		info.SetTimes(-1)
		goto switchAccount
	}
	if result.Output.Url == "" {
		fmt.Println(string(data))
		return "", fmt.Errorf("拿不到图片压缩地址")
	}
	info.SetTimes(1)
	if strings.Contains(ext, "gif") {
		err = clFile.Download(result.Output.Url, new_path)
		if err != nil {
			clLog.Error("下载Tinypng的压缩图片地址失败:%v", err)
		}
		return new_path, err
	}
	path, err := Transform(result.Output.Url)
	if err != nil {
		clLog.Error("转换更小的图片失败:%v", err)
		err = clFile.Download(result.Output.Url, new_path)
		if err != nil {
			clLog.Error("下载Tinypng的压缩图片地址失败:%v", err)
		}
		return new_path, err
	}
	return path, err
}

func Transform(output string) (string, error) {
transform:
	info := GetOne()
	if info == nil {
		return "", nil
	}
	data := clJson.M{
		"convert": clJson.M{
			"type": "*/*",
		},
		"transform": clJson.M{"background": "#000000"},
	}
	reqBody, _ := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPost, output, bytes.NewBuffer(reqBody))
	if err != nil {
		clLog.Error("转换更小分辨率失败:%v", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(info.Email, info.Apikey)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		clLog.Error("错误:%v", err)
		return "", err
	}
	// 解析请求
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		clLog.Error("解析Tinypng的返回失败:%v", err)
		return "", err
	}
	responseContentType := response.Header.Get("Content-Type")
	ext := ""
	if responseContentType == "image/webp" {
		ext = ".webp"
	} else if responseContentType == "image/png" {
		ext = ".png"
	} else if responseContentType == "image/jpeg" {
		ext = ".jpg"
	} else {
		// 错误的结果
		clLog.Error("转换失败:%v", string(responseData))
		goto transform
	}
	new_path := fmt.Sprintf("%v%v%v_%v_compress%v", ROOT_DIR, COMPRESS_DIR, time.Now().UnixNano(), clCommon.Md5([]byte(output)), ext)
	err = ioutil.WriteFile(new_path, responseData, 0644)
	if err != nil {
		clLog.Error("写入文件失败:%v", err)
	}
	return new_path, err
}
