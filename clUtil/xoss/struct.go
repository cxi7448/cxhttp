package xoss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/minio/minio-go"
)

/*
*
s3跟阿里云上传插件
*/
type XClient struct {
	EndPoint     string
	AccessKey    string
	Bucket       string
	SecretKey    string
	Regin        string
	Domain       string
	Comment      string // 用途
	Type         uint   // 分类 0是s3
	s3Client     *minio.Client
	aliyunClient *oss.Client
}
