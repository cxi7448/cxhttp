package xcompress

import (
	"fmt"
	"testing"
)

func TestCompressImage(t *testing.T) {
	_multiple := float32(1)
	width := 1230
	width = int(float32((10/_multiple)/10) * float32(width))
	fmt.Println(width, float32(float32(10/_multiple)/10))
}
