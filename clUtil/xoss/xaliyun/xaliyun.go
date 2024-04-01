package xaliyun

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"sync"
)

type XAliyun struct {
	EndPoint  string
	Bucket    string
	AccessKey string
	SecretKey string
	Region    string
	Client    *oss.Client
	err       error
}

var (
	ENDPOINT  = ""
	BUCKET    = ""
	ACCESSKEY = ""
	SECRETKEY = ""
	REGION    = ""
)

func Init(aliyun XAliyun) {
	ENDPOINT = aliyun.EndPoint
	BUCKET = aliyun.Bucket
	ACCESSKEY = aliyun.AccessKey
	SECRETKEY = aliyun.SecretKey
	REGION = aliyun.Region
}

var sClient *XAliyun
var slocker sync.RWMutex

func New() *XAliyun {
	slocker.Lock()
	defer slocker.Unlock()
	if sClient != nil {
		return sClient
	}
	result := &XAliyun{
		EndPoint:  ENDPOINT,
		Bucket:    BUCKET,
		AccessKey: ACCESSKEY,
		SecretKey: SECRETKEY,
		Region:    REGION,
	}
	if sClient == nil {
		client, err := oss.New(result.EndPoint, result.AccessKey, result.SecretKey)
		result.err = err
		result.Client = client
	}
	sClient = result
	return result
}
func (this *XAliyun) UploadFile(localPath, objectName string, _tryCount int) error {
	if this.err != nil {
		return this.err
	}
	if _tryCount <= 0 {
		return fmt.Errorf("上传失败:%v", this.err)
	}
	bucket, err := this.Client.Bucket(this.Bucket)
	if err != nil {
		return err
	}
	err = bucket.PutObjectFromFile(objectName, localPath)
	if err != nil {
		clLog.Error("错误:%v", err)
		return this.UploadFile(localPath, objectName, _tryCount-1)
	}
	return nil
}

func (this *XAliyun) Exists(objectName string) (bool, error) {
	if this.err != nil {
		return false, this.err
	}
	bucket, err := this.Client.Bucket(this.Bucket)
	if err != nil {
		return false, err
	}
	return bucket.IsObjectExist(objectName)
}
