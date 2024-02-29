package xffmpeg

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clFile"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"os"
)

const downloadUrl = "https://github.com/cxi7448/cxhttp/raw/v1.0.7/clUtil/xffmpeg/linux/ffmpeg"

func DownloadFFMPEG() {
	pwd, _ := os.Getwd()
	file := fmt.Sprintf("%v/ffmpeg", pwd)
	if clFile.IsFile(file) {
		return
	}
	err := clFile.DownloadProcess(downloadUrl, file)
	if err != nil {
		clLog.Error("下载ffmpeg错误:%v", err)
	} else {
		clLog.Info("下载完毕!")
		os.Chmod(file, 0755)
	}
}
