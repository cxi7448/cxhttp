package xcompress

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clFile"
	"github.com/cxi7448/cxhttp/clUtil/ximage"
	_ "image/png" // 导入PNG支持包
	"testing"
)

var filename = "xx.jpg"

func TestCompressMultiple(t *testing.T) {
	clFile.Download("https://mtimg.qkmkpz.com/preview/17383d98edd6422c2ee9e129a1051698.jpg", filename)
	//ImageToWebp("b.png", "b2.png")
}

func TestCompressImage(t *testing.T) {
	fmt.Println(ximage.ImageToWebp(filename, 50))
}
