package xcompress

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/cxi7448/cxhttp/clUtil/xcompress/resize"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	ROOT_DIR     = "./tmp/"
	COMPRESS_DIR = "compress/"
)

func init() {
	os.MkdirAll(ROOT_DIR+COMPRESS_DIR, 0777)
}

func getFiletype(file string) (string, error) {
	f, err := os.Open(file)
	header := make([]byte, 512)
	var f2 = &os.File{}
	f2 = f
	_, err = f2.Read(header)
	if err != nil {
		clLog.Error("读取文件[%v]头失败:%v", file, err)
		return "", err
	}
	filetype := http.DetectContentType(header)
	return filetype, nil
}

/*
*
filepath:图片文件路径
mutiple:1-10
*/
func CompressMultiple(filepath string, _multiple float32) (string, error) {
	file, _ := os.Open(filepath)
	img, _, _ := image.DecodeConfig(file)
	if _multiple > 10 {
		_multiple = 10
	}
	times := float32((10 / _multiple) / 10)
	new_width := int(times * float32(img.Width))
	new_height := int(times * float32(img.Height))
	return CompressImage(filepath, uint(new_width), uint(new_height))
}

func CompressImage(file string, width, height uint) (string, error) {
	var filepath string
	var err error
	index := strings.LastIndex(file, ".")
	if index < 0 {
		clLog.Error("非法文件[%v]!", file)
		return filepath, err
	}
	filetype, err := getFiletype(file)
	if err != nil {
		return filepath, err
	}
	f, err1 := os.Open(file)
	if err1 != nil {
		clLog.Error("打开文件[%v]失败! %v", file, err1)
		return filepath, err1
	}
	defer f.Close()
	switch filetype {
	case "image/jpeg":
		filepath, err = compressJPEG(f, width, height)
		break
	case "image/jpg":
		filepath, err = compressJPEG(f, width, height)
		break
	case "image/png":
		filepath, err = compressPNG(f, width, height)
		break
	//case "gif":
	//filepath, err = compressGIF(f, width, height)
	//return filepath, fmt.Errorf("暂不支持GIF压缩!")
	default:
		return file, nil
		//filepath, err = compressJPEG(f, width, height)
	}
	if err != nil {
		clLog.Error("图片压缩失败:%v", err)
		return filepath, err
	}
	return filepath, nil
}

func compressGIF(f io.Reader, width uint, height uint) (string, error) {
	g, err := gif.DecodeAll(f)
	if err != nil {
		clLog.Error("读取gif配置失败:%v", err)
		return "", err
	}
	var new_g = &gif.GIF{
		//Image:     g.Image,
		Delay:     g.Delay,
		LoopCount: g.LoopCount,
		Disposal:  g.Disposal,
	}
	var new_paletted = []*image.Paletted{}
	var i = 1
	for _, img := range g.Image {
		//gif.Decode()
		_img := resize.Resize(width, height, img, resize.Lanczos3)
		filepath := fmt.Sprintf("%v%v%v_compress.gif", ROOT_DIR, COMPRESS_DIR, i)
		out, err := os.Create(filepath)
		if err != nil {
			clLog.Error("创建临时文件失败:%v", err)
			return filepath, err
		}
		gif.Encode(out, _img, nil)
		//f, err := os.Open(filepath)
		//if err != nil {
		//	clLog.Error("文件[%v]打开出错:%v", filepath, err)
		//	return "", err
		//}
		png_img, err := gif.Decode(out)
		if err != nil {
			clLog.Error("文件[%v]打开出错2:%v", filepath, err)
			return "", err
		}
		i++
		//p := image.NewPaletted(image.Rect(0, 0, int(width), int(height)), palette.Plan9)
		p := image.NewPaletted(image.Rect(0, 0, int(width), int(height)), img.Palette)
		draw.Draw(p, p.Bounds(), png_img, image.ZP, draw.Src)
		new_paletted = append(new_paletted, p)
	}
	new_g.Image = new_paletted
	filepath := fmt.Sprintf("%v%v%v_compress.gif", ROOT_DIR, COMPRESS_DIR, time.Now().UnixNano())
	out, err := os.Create(filepath)
	if err != nil {
		clLog.Error("创建临时文件失败:%v", err)
		return "", err
	}
	err = gif.EncodeAll(out, new_g)
	return filepath, err
}

//
//func compressJPG(f io.Reader, width uint, height uint) {
//	compressJPEG(f, width, height)
//}

func compressJPEG(f io.Reader, width uint, height uint) (string, error) {
	var img image.Image
	var err error
	var filepath string
	img, err = jpeg.Decode(f)
	if err != nil {
		clLog.Error("转换img对象失败! %v", err)
		return filepath, err
	}
	new_img := resize.Resize(width, height, img, resize.Lanczos3)
	filepath = fmt.Sprintf("%v%v%v_compress.jpg", ROOT_DIR, COMPRESS_DIR, time.Now().UnixNano())
	out, err := os.Create(filepath)
	if err != nil {
		clLog.Error("创建临时文件失败:%v", err)
		return filepath, err
	}
	err = jpeg.Encode(out, new_img, nil)
	if err != nil {
		clLog.Error("生成新的压缩文件失败:%v", err)
		return filepath, err
	}
	return filepath, err
}

func compressPNG(f io.Reader, width uint, height uint) (string, error) {
	var img image.Image
	var err error
	var filepath = fmt.Sprintf("%v%v%v_compress.png", ROOT_DIR, COMPRESS_DIR, time.Now().UnixNano())
	img, err = png.Decode(f)
	if err != nil {
		clLog.Error("转换img对象失败! %v", err)
		return filepath, err
	}
	new_img := resize.Resize(width, height, img, resize.Lanczos3)
	out, err := os.Create(filepath)
	if err != nil {
		clLog.Error("创建临时文件失败:%v", err)
		return filepath, err
	}
	err = png.Encode(out, new_img)
	if err != nil {
		clLog.Error("生成新的压缩文件失败:%v", err)
		return filepath, err
	}
	return filepath, err
}
