package ximage

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clFile"
	"os"
	"runtime"
	"strings"
)

var command = ""
var command_cwebp = ""
var command_gif2webp = ""

const (
	bin_root = "./bin"
)

func init() {
	// 自动生成执行文件
	var paths = []string{}
	var filename_cwebp = ""
	var filename_gif2webp = ""
	fmt.Println(runtime.GOOS)
	if runtime.GOOS == "window" {
		paths = strings.Split(os.Getenv("PATH"), ";")
		filename_cwebp = "cwebp.exe"
		filename_gif2webp = "gif2webp.exe"
	} else {
		filename_cwebp = "cwebp"
		filename_gif2webp = "gif2webp"
		paths = strings.Split(os.Getenv("PATH"), ":")
	}
	// 扫描当前目录是否存在执行文件
	//pwd, _ := os.Getwd()
	//command_cwebp = fmt.Sprintf("%v/%v", pwd, filename_cwebp)
	//command_gif2webp = fmt.Sprintf("%v/%v", pwd, filename_gif2webp)
	if len(paths) > 0 {
		// 环境变量优先级最高
		for _, path := range paths {
			cwebp_cmd := fmt.Sprintf("%v/%v", path, filename_cwebp)
			if clFile.IsFile(cwebp_cmd) {
				command_cwebp = cwebp_cmd
			}

			gif2webp_cmd := fmt.Sprintf("%v/%v", path, filename_gif2webp)
			if clFile.IsFile(gif2webp_cmd) {
				command_gif2webp = gif2webp_cmd
			}
		}
	}
	if command_cwebp == "" {
		// 自动下载补齐 cwebp
		link := fmt.Sprintf("")
		clFile.Download(link, fmt.Sprintf("%vcwebp", bin_root))
	}
	if command_gif2webp == "" {
		// 自动下载补齐 gif2webp
	}
}
