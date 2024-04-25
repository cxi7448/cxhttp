package txyun

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/cxi7448/cxhttp/clUtil/xoss"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

/*
*
腾讯云
*/
type XTxyun struct {
	xoss.Config
	Client *cos.Client
	err    error
}

func New(config xoss.Config) *XTxyun {
	result := &XTxyun{
		Config: config,
	}
	u, _ := url.Parse(result.Domain)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  result.AccessKey, // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
			SecretKey: result.SecretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
		},
	})
	result.Client = c
	return result
}

func (this *XTxyun) UploadFile(objectName, localPath string, tryCount int) error {
	if this.Client == nil {
		return fmt.Errorf("创建Client失败!")
	}
	if tryCount <= 0 {
		return fmt.Errorf("腾讯云上传失败！")
	}
	_, err := this.Client.Object.PutFromFile(context.Background(), objectName, localPath, nil)
	if err != nil {
		clLog.Error("腾讯云上传错误:%v", err)
		return this.UploadFile(objectName, localPath, tryCount-1)
	}
	return err
}

func (this *XTxyun) UploadContent(objectName string, content []byte, tryCount int) error {
	if this.Client == nil {
		return fmt.Errorf("创建Client失败!")
	}
	if tryCount <= 0 {
		return fmt.Errorf("腾讯云上传失败！")
	}
	_, err := this.Client.Object.Put(context.Background(), objectName, bytes.NewReader(content), nil)
	if err != nil {
		clLog.Error("腾讯云上传错误:%v", err)
		return this.UploadContent(objectName, content, tryCount-1)
	}
	return err
}
