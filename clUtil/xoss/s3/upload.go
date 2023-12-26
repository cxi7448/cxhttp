package s3

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clCrypt"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"io/ioutil"
)

var (
	AES_KEY = []byte("jp9dzn3wg15m3l31wti538qbu38hemsd")
	AES_IV  = []byte("o0y70hehd961vfec")
	PREFIX  = "aesencipher/"
)

// 上传加密图片
// filepath string: 图片的本地地址
// objectName string: 图片存储桶的地址
func UploadEncodeImage(filepath string, objectName string) (error, string) {
	_content, err := ioutil.ReadFile(filepath)
	if err != nil {
		clLog.Error("文件[%v]打开失败:%v", filepath, err)
		return err, ""
	}
	encrypted, err := clCrypt.AesCFBEncrypt(_content, AES_KEY, AES_IV)
	if err != nil {
		clLog.Error("文件[%v]加密失败:%v", filepath, err)
		return err, ""
	}
	new_object := fmt.Sprintf("%v%v", PREFIX, objectName)
	s3 := NewImage()
	return s3.UploadContent(encrypted, new_object, 3), new_object
}

// 上传一般图片
// filepath string: 图片的本地地址
// objectName string: 图片存储桶的地址
func UploadImage(filepath string, objectName string) error {
	s3 := NewImage()
	return s3.UploadFile(filepath, objectName, 3)
}

/*
*
上传加密跟一般图片
filepath string: 图片的本地地址
objectName string: 图片存储桶的地址
*/
func UploadImageBoth(filepath string, objectName string) (error, string) {
	err := UploadImage(filepath, objectName)
	if err != nil {
		clLog.Error("上传原图失败:%v", err)
		return err, ""
	}
	err, newObjectName := UploadEncodeImage(filepath, objectName)
	if err != nil {
		clLog.Error("上传加密图片失败")
	}
	return err, newObjectName
}

// 上传加密内容
func UploadEncodeContent(content []byte, objectName string) (error, string) {
	s3 := NewImage()
	err := s3.UploadContent(content, objectName, 3)
	if err != nil {
		return err, ""
	}

	encrypted, err := clCrypt.AesCFBEncrypt(content, AES_KEY, AES_IV)
	if err != nil {
		clLog.Error("文件[%v]加密失败:%v", objectName, err)
		return err, ""
	}
	new_object := fmt.Sprintf("%v%v", PREFIX, objectName)
	return s3.UploadContent(encrypted, new_object, 3), new_object
}

// 上传视频文件
func UploadVideoFile(filepath, objectName string) error {
	s3 := NewVideo()
	return s3.UploadFile(filepath, objectName, 3)
}

// 上传一般文件
func UploadFile(filepath, objectName string) error {
	s3 := NewImage()
	return s3.UploadFile(filepath, objectName, 3)
}

// 上传一般文件
func UploadContent(content []byte, objectName string) error {
	s3 := NewImage()
	return s3.UploadContent(content, objectName, 3)
}
