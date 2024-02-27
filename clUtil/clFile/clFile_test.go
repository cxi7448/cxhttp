package clFile

import (
	"fmt"
	"testing"
)

func TestGetFileMD5(t *testing.T) {
	err := DownloadProcess("https://github.com/cxi7448/cxhttp/raw/v1.0.7/clUtil/xffmpeg/linux/ffmpeg", "ffmpeg")
	fmt.Println(err)
}
