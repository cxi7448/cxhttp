package bdyun

import (
	"bytes"
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/cxi7448/cxhttp/clUtil/xoss"
)

type XBdYun struct {
	xoss.Config
	Client *bos.Client
}

func NewWith(config xoss.Config) *XBdYun {
	result := &XBdYun{
		Config: config,
	}
	client, err := bos.NewClient(result.AccessKey, result.SecretKey, result.EndPoint)
	if err == nil {
		result.Client = client
	}
	return result
}

func (this *XBdYun) UploadContent(objectName string, content []byte, tryCount int) error {
	if this.Client == nil {
		return fmt.Errorf("初始化失败")
	}
	if tryCount <= 0 {
		return fmt.Errorf("访问百度BOS失败")
	}
	_, err := this.Client.PutObjectFromStream(this.Bucket, objectName, bytes.NewBuffer(content), nil)
	if err != nil {
		clLog.Error("上传百度云错误：%v", err)
		return this.UploadContent(objectName, content, tryCount-1)
	}
	return nil
}

func (this *XBdYun) UploadFile(objectName, localPath string, tryCount int) error {
	if this.Client == nil {
		return fmt.Errorf("初始化失败")
	}
	if tryCount <= 0 {
		return fmt.Errorf("访问百度BOS失败")
	}
	_, err := this.Client.PutObjectFromFile(this.Bucket, objectName, localPath, nil)
	if err != nil {
		clLog.Error("上传百度云错误：%v", err)
		return this.UploadFile(objectName, localPath, tryCount-1)
	}
	return nil
}
