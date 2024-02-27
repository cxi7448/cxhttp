package s3

import (
	"bytes"
	"fmt"
	"github.com/cxi7448/cxhttp/clCommon"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/minio/minio-go"
	"io/ioutil"
	"sync"
)

type S3 struct {
	EndPoint  string
	Bucket    string
	AccessKey string
	SecretKey string
	Region    string
	Domain    string
	Client    *minio.Client
	t         string // image,video, 其他
}

var sS3 = []S3{}
var sPool = make(map[string]*S3)
var sLocker sync.RWMutex

func SetVideo(s3 S3) {
	Set(s3, "VIDEO")
}
func SetImage(s3 S3) {
	Set(s3, "IMAGE")
}
func Set(s3 S3, _type string) {
	s3.t = _type
	isSet := false
	for key, val := range sS3 {
		if val.t == _type {
			sS3[key] = s3
			isSet = true
			break
		}
	}
	if !isSet {
		sS3 = append(sS3, s3)
	}
}
func get(_type string) *S3 {
	for _, s3 := range sS3 {
		if s3.t == _type {
			return &s3
		}
	}
	return &S3{t: _type}
}

// 图片
func NewImage() *S3 {
	return NewByType("IMAGE")
}

func NewByType(_type string) *S3 {
	return NewS3(get(_type))
}
func NewVideo() *S3 {
	return NewByType("VIDEO")
}
func NewS3(s3 *S3) *S3 {
	if s3.Bucket == "" {
		return &S3{}
	}
	sLocker.RLock()
	obj, ok := sPool[s3.genToken()]
	sLocker.RUnlock()
	if ok {
		return obj
	}
	var client *minio.Client
	var err error
	if s3.Region != "" {
		client, err = minio.NewWithRegion(s3.EndPoint, s3.AccessKey, s3.SecretKey, true, s3.Region)
		if err != nil {
			clLog.Error("初始化minio客户端错误: %v", err)
		}
	} else {
		client, err = minio.New(s3.EndPoint, s3.AccessKey, s3.SecretKey, true)
		if err != nil {
			clLog.Error("初始化minio客户端错误: %v", err)
		}
	}
	if client != nil {
		s3.Client = client
		sLocker.Lock()
		sPool[s3.genToken()] = s3
		sLocker.Unlock()
	}
	return s3
}
func (this *S3) UploadFile(filepath, objectName string, _tryCount int) error {
	if this.Client == nil {
		return fmt.Errorf("存储桶实例失败!")
	}
	_, err := this.Client.FPutObject(this.Bucket, objectName, filepath, minio.PutObjectOptions{
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	})
	if err != nil {
		clLog.Error("上传文件错误: %v", err)
		if _tryCount > 0 {
			return this.UploadFile(filepath, objectName, _tryCount-1)
		}
		return err
	}
	return nil
}
func (this *S3) UploadContent(content []byte, objectName string, _tryCount int) error {
	if this.Client == nil {
		return fmt.Errorf("存储桶实例失败!")
	}
	_, err := this.Client.PutObject(this.Bucket, objectName, bytes.NewBuffer(content), int64(len(content)), minio.PutObjectOptions{
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	})
	if err != nil {
		clLog.Error("上传文件错误: %v", err)
		if _tryCount > 0 {
			return this.UploadContent(content, objectName, _tryCount-1)
		}
		return err
	}
	return nil
}

func (this *S3) GetContent(objectName string, _tryCount int) ([]byte, error) {
	if this.Client == nil {
		return nil, fmt.Errorf("存储桶实例失败!")
	}
	obj, err := this.Client.GetObject(this.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		clLog.Error("读取%v内容失败:%v", objectName, err)
		if _tryCount > 0 {
			return this.GetContent(objectName, _tryCount-1)
		}
		return nil, err
	}
	defer obj.Close()
	content, err := ioutil.ReadAll(obj)
	if err != nil {
		clLog.Error("读取内容失败:%v", err)
		return nil, err
	}
	return content, nil
}

func (this *S3) Download(objectName, localPath string, _tryCount int) error {
	if this.Client == nil {
		return fmt.Errorf("存储桶实例失败!")
	}
	err := this.Client.FGetObject(this.Bucket, objectName, localPath, minio.GetObjectOptions{})
	if err != nil {
		if _tryCount > 0 {
			return this.Download(objectName, localPath, _tryCount-1)
		}
	}
	return err
}

// s3对象
func (this *S3) C() *minio.Client {
	return this.Client
}

func (this *S3) genToken() string {
	return clCommon.Md5([]byte(fmt.Sprintf("%v:%v:%v:%v:%v", this.EndPoint, this.AccessKey, this.SecretKey, this.Bucket, this.Region)))
}
