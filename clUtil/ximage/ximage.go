package ximage

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clCommon"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"golang.org/x/image/webp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	ROOT_DIR = "./tmp/"
)

func init() {
	os.MkdirAll(ROOT_DIR, 0777)
}

const perSize = 1000 // 1kb

// 设置最大字节，超过最大字节的话，使用min_quality压缩 默认quality压缩
// max_size 单位: KB
// max_size 单位: KB
func ImageToWebpByMaxSize(input string, max_size float64, quality, min_quality int) (string, error) {
	path, err := ImageToWebp(input, quality)
	if err != nil {
		clLog.Error("压缩图片失败:%v", err)
		return "", err
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		clLog.Error("读取压缩后的图片失败：%v", err)
		return "", err
	}
	if float64(len(content)) >= max_size {
		defer os.RemoveAll(path)
		new_path, err := ImageToWebp(path, min_quality)
		if err != nil {
			clLog.Error("压缩图片失败:%v", err)
			return "", err
		}
		return new_path, nil
	}
	return path, nil
}

func ImageToWebp(input string, quality int) (string, error) {
	ext := input[strings.LastIndex(input, "."):]
	var output = fmt.Sprintf("%v%v_%v%v", ROOT_DIR, time.Now().UnixNano(), clCommon.Md5([]byte(input)), ext)
	imageType, err := GetImageType(input)
	if err != nil {
		clLog.Error("读取图片类型失败:%v", err)
		return "", err
	}
	clLog.Info("图片类型:%v-%v", input, imageType)
	if imageType == "image/gif" {
		_, err = clCommon.RunCommandNoConsole(command_gif2webp, input, "-quiet", "-q", fmt.Sprint(quality), "-o", output)
		return output, err
	}

	_, err = clCommon.RunCommandNoConsole(command_cwebp, input, "-quiet", "-q", fmt.Sprint(quality), "-o", output)
	return output, err
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
	fmt.Println(fileType)
	return fileType == "image/png"
}

func GetImageType(file string) (string, error) {
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
