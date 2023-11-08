package clFile

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// 创建文件夹如果不存在
func CreateDirIFNotExists(_path string) {
	_, err := os.Open(_path)
	if os.IsNotExist(err) {
		os.MkdirAll(_path, 0766)
	}
}



// 删除文件
func DelFile(_path string) {
	err := os.RemoveAll(_path)
	if err != nil {
		fmt.Printf("删除失败: %v", err)
	}
}


// 读入文件
func ReadFile(_filename string, _createIfNotExists bool) string {
	content, err := ioutil.ReadFile(_filename)
	if os.IsNotExist(err) {
		if _createIfNotExists {
			pFile, err := os.Create(_filename)
			if err == nil {
				pFile.Close()
			}
		}
		return ""
	}
	return string(content)
}


// 文件追加
func AppendFile(_filename, _content string) {
	pFile, err := os.OpenFile(_filename, os.O_RDWR, os.ModePerm)
	if os.IsNotExist(err) {
		pFile, err = os.Create(_filename)
		if err != nil {
			return
		}
	}
	pFile.Seek(0, io.SeekEnd)
	pFile.Write([]byte(_content))
	pFile.Close()
}



// 获取文件名
func GetFileName(_path string) string {
	fileInfo, err := os.Stat(_path)
	if err != nil {
		fmt.Printf("获取文件: %v 名失败! 错误:%v", _path, err)
		return ""
	}
	return fileInfo.Name()
}


// 获取文件名
func GetFileSize(_path string) int64 {
	fileInfo, err := os.Stat(_path)
	if err != nil {
		fmt.Printf("获取文件: %v 名失败! 错误:%v", _path, err)
		return 0
	}
	return fileInfo.Size()
}



// 获取文件MD5值
func GetFileMD5(_path string) string {
	content, err := ioutil.ReadFile(_path)
	if err != nil {
		fmt.Printf("打开文件: %v 失败! 错误:%v", _path, err)
		return ""
	}

	h := md5.New()
	h.Write(content) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}


// 文件是否存在
func FileIsExists(_filePath string) bool {
	_, err := os.Stat( _filePath )
	if os.IsNotExist(err) {
		return false
	}
	return true
}
