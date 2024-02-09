package xcompress

import (
	"fmt"
	"testing"
)

func TestCompressImage(t *testing.T) {
	path, err := JpgToPng("a.jpg")
	fmt.Println(path)
	fmt.Println(err)
}
