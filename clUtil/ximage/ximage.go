package ximage

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"golang.org/x/image/webp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"time"
)

const (
	ROOT_DIR = "./tmp/"
)

func init() {
	os.MkdirAll(ROOT_DIR, 0777)
}

const perSize = 1000 // 1kb

func ImageToWebp(input string, quality int) (string, error) {
	info := New(input)
	err := info.ImageToWebp(quality)
	if err != nil {
		clLog.Error("错误:%v", err)
	}
	return info.Output, info.Err
}

func ImageAdaptionToSize(input string, max_size float64, _quality, _min_quality int) (string, error) {
	info := ImageToWebpV2(input, _quality)
	if info.IsError() {
		clLog.Error("图片压缩失败：%v", info.Err)
		return "", info.Err
	}
	clLog.Info("大小:%v quality:%v", float64(len(info.Buffer))/1000, _quality)
	if info.IsWebp() {
		return info.Output, nil
	}
	if _quality <= _min_quality {
		return info.Output, nil
	}
	size := float64(len(info.Buffer)) / 1000
	if size <= max_size {
		return info.Output, nil
	}
	os.RemoveAll(info.Output)
	return ImageAdaptionToSize(input, max_size, _quality-5, _min_quality)
}
func ImageToWebpV2(input string, quality int) *Ximage {
	xi := New(input)
	err := xi.ImageToWebp(quality)
	if err != nil {
		clLog.Error("转换错误:%v", err)
	} else {
		// 更新字节大小
		xi.ReadFile(xi.Output)
	}
	return xi
}

func ImageToPng(localPath string) (string, error) {
	imageType, err := GetImageType(localPath)
	if err != nil {
		clLog.Error("读取图片类型失败:%v", err)
		return "", err
	}
	if imageType == "image/webp" {
		return WebpToPng(localPath)
	}
	return JpgToPng(localPath)
}

func WebpToPng(localPath string) (string, error) {
	// 打开原始JPG图片文件
	srcFile, err := os.Open(localPath)
	if err != nil {
		clLog.Error("读取JPG失败:%v", err)
		return "", err
	}
	defer srcFile.Close()

	// 读取JPG图片数据并创建Image对象
	img, err := webp.Decode(srcFile)
	if err != nil {
		clLog.Error("读取JPG失败:%v", err)
		return "", err
	}

	// 设置目标PNG图片路径及名称
	var outputPath = fmt.Sprintf("%v%v.png", ROOT_DIR, time.Now().UnixNano())
	// 保存为PNG格式
	dstFile, err := os.Create(outputPath)
	if err != nil {
		clLog.Error("创建PNG失败:%v", err)
		return "", err
	}
	defer dstFile.Close()

	err = png.Encode(dstFile, img)
	if err != nil {
		clLog.Error("无法保存为PNG格式:%v", err)
		return "", err
	}
	return outputPath, nil
}

func JpgToPng(localPath string) (string, error) {
	// 打开原始JPG图片文件
	srcFile, err := os.Open(localPath)
	if err != nil {
		clLog.Error("读取JPG失败:%v", err)
		return "", err
	}
	defer srcFile.Close()

	// 读取JPG图片数据并创建Image对象
	img, _, err := image.Decode(srcFile)
	if err != nil {
		clLog.Error("读取JPG失败:%v", err)
		return "", err
	}

	// 设置目标PNG图片路径及名称
	var outputPath = fmt.Sprintf("%v%v.png", ROOT_DIR, time.Now().UnixNano())
	// 保存为PNG格式
	dstFile, err := os.Create(outputPath)
	if err != nil {
		clLog.Error("创建PNG失败:%v", err)
		return "", err
	}
	defer dstFile.Close()

	err = png.Encode(dstFile, img)
	if err != nil {
		clLog.Error("无法保存为PNG格式:%v", err)
		return "", err
	}
	return outputPath, nil
}

func IsPng(localPath string) bool {
	fileType, err := GetImageType(localPath)
	if err != nil {
		return false
	}
	return fileType == "image/png"
}

func GetImageType(file string) (string, error) {
	info := New(file)
	return info.GetImageType(), nil
}
