package ximage

import (
	"bytes"
	"fmt"
	"github.com/chai2010/webp"
	"github.com/cxi7448/cxhttp/clUtil/clCommon"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"image"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Ximage struct {
	Input     string
	Output    string
	ImageType string
	Buffer    []byte
	Err       error
	Ext       string
}

func New(input string) *Ximage {
	result := &Ximage{
		Input: input,
		Ext:   input[strings.LastIndex(input, "."):],
	}
	result.ReadFile(input)
	result.ImageType = result.GetImageType()
	return result
}

func (this *Ximage) ReadFile(input string) *Ximage {
	buffer, err := ioutil.ReadFile(input)
	if err != nil {
		this.Err = err
		return this
	}
	this.Buffer = buffer
	return this
}

func (this *Ximage) IsError() bool {
	return this.Err != nil
}
func (this *Ximage) ImageToWebp(quality int) error {
	if this.IsError() {
		clLog.Error("错误:%v", this.Err)
		return this.Err
	}
	this.Output = fmt.Sprintf("%v%v_%v%v", ROOT_DIR, time.Now().UnixNano(), clCommon.Md5([]byte(this.Input)), this.Ext)
	if this.IsWebp() {
		this.Output = this.Input
		this.Err = nil
		return nil
	}
	if this.IsGif() {
		_, err := clCommon.RunCommandNoConsole(command_gif2webp, this.Input, "-quiet", "-q", fmt.Sprint(quality), "-o", this.Output)
		this.Err = err
		return err
	} else {
		input, err := ioutil.ReadFile(this.Input)
		if err != nil {
			panic(err)
		}
		img, _, err := image.Decode(bytes.NewBuffer(input))
		if err != nil {
			this.Err = err
			return err
		}
		webpBytes, err := webp.EncodeRGB(img, float32(quality))
		if err != nil {
			this.Err = err
			return err
		}
		err = ioutil.WriteFile(this.Output, webpBytes, 0666)
		//_, err := clCommon.RunCommandNoConsole(command_cwebp, this.Input, "-quiet", "-q", fmt.Sprint(quality), "-o", this.Output)
		this.Err = err
		return err
	}
}

func (this *Ximage) GetImageType() string {
	if this.ImageType != "" {
		return this.ImageType
	}
	filetype := http.DetectContentType(this.Buffer[0:512])
	this.ImageType = filetype
	return this.ImageType
}

func (this *Ximage) IsWebp() bool {
	return this.GetImageType() == "image/webp"
}

func (this *Ximage) IsGif() bool {
	return this.GetImageType() == "image/gif"
}
func (this *Ximage) IsPng() bool {
	return this.GetImageType() == "image/png"
}
