package ximage

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clFile"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"os"
	"runtime"
)

var command = ""
var command_cwebp = ""
var command_gif2webp = ""

const (
	host   = "https://github.com/cxi7448/cxhttp/raw/v1.0.6"
	folder = "bin"
)

func init() {
	pwd, _ := os.Getwd()
	// 自动生成执行文件
	var filename_cwebp = ""
	var filename_gif2webp = ""
	bin_root := fmt.Sprintf("%v/%v", pwd, folder)
	os.MkdirAll(bin_root, 0700)
	server := "linux"
	if runtime.GOOS == "windows" {
		filename_cwebp = "cwebp.exe"
		filename_gif2webp = "gif2webp.exe"
		server = "window"
	} else {
		if runtime.GOOS == "darwin" {
			server = "mac"
		}
		filename_cwebp = "cwebp"
		filename_gif2webp = "gif2webp"
	}
	cwebp_cmd := fmt.Sprintf("%v/%v", bin_root, filename_cwebp)
	gif2webp_cmd := fmt.Sprintf("%v/%v", bin_root, filename_gif2webp)
	if clFile.IsFile(cwebp_cmd) {
		command_cwebp = cwebp_cmd
	} else {
		link := fmt.Sprintf("%v/clUtil/ximage/%v/%v", host, server, filename_cwebp)
		err := clFile.Download(link, cwebp_cmd)
		if err != nil {
			clLog.Error("cwebp[%v]-[%v]-[%v]下载失败:%v", link, server, runtime.GOOS, err)
		} else {
			command_cwebp = fmt.Sprintf("%v/%v", bin_root, filename_cwebp)
			os.Chmod(command_cwebp, 0755)
		}
	}
	if clFile.IsFile(gif2webp_cmd) {
		command_gif2webp = gif2webp_cmd
	} else {
		// 自动下载补齐 gif2webp
		link := fmt.Sprintf("%v/clUtil/ximage/%v/%v", host, server, filename_gif2webp)
		err := clFile.Download(link, gif2webp_cmd)
		if err != nil {
			clLog.Error("cwebp[%v]-[%v]-[%v]下载失败:%v", link, server, runtime.GOOS, err)
		} else {
			command_gif2webp = fmt.Sprintf("%v/%v", bin_root, filename_gif2webp)
			os.Chmod(command_gif2webp, 0755)
		}
	}
	clLog.Info("ximage初始化完毕！")
}
