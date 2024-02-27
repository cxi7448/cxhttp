package clFile

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/cxi7448/cxhttp/clUtil/clCommon"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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
	_, err := os.Stat(_filePath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func IsFile(_filepath string) bool {
	f, err := os.Stat(_filepath)
	if os.IsNotExist(err) || f == nil {
		return false
	}
	if f.IsDir() {
		return false
	}
	return true
}

func Download(link, localPath string) error {
	res, err := http.Get(link)
	if err != nil {
		return err
	}
	f, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = io.Copy(f, res.Body)
	return err
}

func Copy(filePath, newPath string) error {
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(newPath, input, 0644)
	return err
}

func DownloadProcess(link, localPath string) error {
	tmpFile := localPath + ".tmp"
	target, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fileSize := clCommon.Int64(resp.Header.Get("Content-Length"))
	bar := pb.Full.Start64(fileSize)
	bar.SetWidth(120)                         //设置进度条宽度
	bar.SetRefreshRate(10 * time.Millisecond) //设置刷新速率
	defer bar.Finish()
	barReader := bar.NewProxyReader(resp.Body)
	if _, err := io.Copy(target, barReader); err != nil {
		target.Close()
		return err
	}
	target.Close()
	if err := os.Rename(tmpFile, localPath); err != nil {
		return err
	}
	return nil
}
