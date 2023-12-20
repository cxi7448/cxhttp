package xffmpeg

import (
	"bytes"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clCommon"
	"math"
	"regexp"
	"strings"
)

/*
*
4K: 3840×2160像素
普清720P分辨率是1280×720
超清1080P分辨率是1920×1080
640*360是流畅分辨率
960*540是标清分辨率
720乘480
1280*720是高清分辨率
1920*1080的为1080P分辨率。
*/
const (
	FB_LC_W = 640 // 流畅
	FB_LC_H = 360

	//FB_BQ_W = 960 //
	//FB_BQ_H = 540

	FB_GQ_W = 720 // 高清
	FB_GQ_H = 480

	FB_CQ_W = 1280 // 超清
	FB_CQ_H = 720

	FB_LG_W = 1980 // 蓝光
	FB_LG_H = 1080

	FB_4K_W = 3840 // 蓝光4K
	FB_4K_H = 2160
)

type VideoInfo struct {
	Filepath    string
	Name        string
	Duration    string
	length      float64
	FenBianLv   string
	Width       uint32
	Has480      bool
	Height      uint32
	IsLandscape bool // 是否横屏
	FB          *FB
	Fps         float64
}
type FB struct {
	Width       uint32
	Md5         string
	Height      uint32
	IsLandscape bool
	Self        bool
	Mp4Filepath string // 转换后的地址
	M3u8Path    string
}

func (this *FB) Resolution() string {
	return fmt.Sprint(this.Height)
}

func (this *FB) Name() string {
	return fmt.Sprint(this.Height)
}
func (this *FB) Scale() string {
	if this.IsLandscape {
		return fmt.Sprintf("%v:%v", this.Width, this.Height)
	} else {
		return fmt.Sprintf("%v:%v", this.Height, this.Width)
	}
}
func newFB(width, height uint32, IsLandscape bool, name string) *FB {
	fb := &FB{
		Width:       width,
		Height:      height,
		IsLandscape: IsLandscape,
	}
	fb.Md5 = clCommon.Md5([]byte(fmt.Sprintf("%v_%v", fb.Resolution(), name)))
	return fb
}
func newLC(IsLandscape bool, name string) *FB {
	return newFB(FB_LC_W, FB_LC_H, IsLandscape, name)
}

//
//func newBQ(IsLandscape bool) FB {
//	return newFB(FB_BQ_W, FB_BQ_H, IsLandscape)
//}

func newGQ(IsLandscape bool, name string) *FB {
	return newFB(FB_GQ_W, FB_GQ_H, IsLandscape, name)
}

func newCQ(IsLandscape bool, name string) *FB {
	return newFB(FB_CQ_W, FB_CQ_H, IsLandscape, name)
}

func newLG(IsLandscape bool, name string) *FB {
	return newFB(FB_LG_W, FB_LG_H, IsLandscape, name)
}

func new4K(IsLandscape bool, name string) *FB {
	return newFB(FB_4K_W, FB_4K_H, IsLandscape, name)
}

func (this *FFmpeg) GetVideoInfo() *VideoInfo {
	if this.videoInfo != nil {
		return this.videoInfo
	}
	filepath := this.input
	index := strings.LastIndex(filepath, "/")
	var name string
	if index >= 0 {
		name = filepath[index+1:]
	} else {
		name = filepath
	}
	out, _ := this.Run()
	if len(out) == 0 || !(bytes.Contains(out, []byte("Duration:")) && bytes.Contains(out, []byte("Stream")) && bytes.Contains(out, []byte("Video:"))) {
		return nil
	}
	video := VideoInfo{
		Filepath: filepath,
		Name:     name,
	}
	for _, line := range bytes.Split(out, []byte("\n")) {
		row := bytes.TrimSpace(line)
		if bytes.HasPrefix(row, []byte("Duration:")) {
			cells := bytes.Split(row, []byte(","))
			video.Duration = strings.TrimSpace(string(cells[0])[9:])
		} else if bytes.HasPrefix(row, []byte("Stream")) && bytes.Contains(row, []byte("Video:")) && bytes.Contains(row, []byte("default")) {
			reg, _ := regexp.Compile(`,\s*[0-9]+x[0-9]+[^0-9]`)
			// , 1920x1080 [SAR 1:1 DAR 16:9],
			fenbianlv := reg.Find(row)
			if len(fenbianlv) > 0 {
				fenbianlv = bytes.Trim(fenbianlv, ",")
				video.FenBianLv = string(bytes.TrimSpace(fenbianlv))
				wh := strings.Split(video.FenBianLv, "x")
				video.Width = clCommon.Uint32(wh[0])
				video.Height = clCommon.Uint32(wh[1])
				video.IsLandscape = video.Width > video.Height
				var width = video.Width
				var height = video.Height
				if !video.IsLandscape {
					width = video.Height
					height = video.Width
				}
				video.FB = newFB(width, height, video.IsLandscape, "")
			}

			// 提取fps
			reg_fps, _ := regexp.Compile(`,\s*[0-9.]+\s*fps\s*,`)
			fps_row := reg_fps.Find(row)
			if len(fps_row) > 0 {
				fps_row = bytes.Trim(fps_row, ",")
				fps := string(bytes.TrimSpace(fps_row))
				video.Fps = clCommon.Float64(fps[:len(fps)-3])
			}
		}
	}
	video.FB.Md5 = clCommon.Md5([]byte(video.Name))
	this.videoInfo = &video
	return &video
}

func (this *VideoInfo) GetLength() float64 {
	if this.length > 0 {
		return this.length
	}
	rows := strings.Split(this.Duration, ":")
	h := clCommon.Float64(rows[0]) * 3600
	m := clCommon.Float64(rows[1]) * 60
	s := clCommon.Float64(rows[2])
	return math.Ceil(h + m + s)
}

func (this *VideoInfo) Is(fbl uint32) bool {
	var width = this.Height
	if !this.IsLandscape {
		width = this.Width
	}
	return width == fbl
}

// 找到合适的分辨率区间
func (this *VideoInfo) FetchFB(folder string) []*FB {
	var result = []*FB{}
	var width = this.Height // 默认横屏
	var height = this.Width
	if !this.IsLandscape {
		// 竖屏
		width = this.Width
		height = this.Height
	}
	if width >= FB_GQ_H {
		fb := newGQ(this.IsLandscape, this.Name)
		fb.Self = width == FB_GQ_H
		real_width := FB_GQ_H / float64(width) * float64(height)
		fb.Width = uint32(math.Ceil(real_width))
		result = append(result, fb)
	}
	if width >= FB_CQ_H {
		fb := newCQ(this.IsLandscape, this.Name)
		fb.Self = width == FB_CQ_H
		real_width := FB_CQ_H / float64(width) * float64(height)
		fb.Width = uint32(math.Ceil(real_width))
		result = append(result, fb)
	}
	if width >= FB_LG_H {
		fb := newLG(this.IsLandscape, this.Name)
		fb.Self = width == FB_LG_H
		real_width := FB_LG_H / float64(width) * float64(height)
		fb.Width = uint32(math.Ceil(real_width))
		result = append(result, fb)
	}
	if len(result) == 0 {
		fb := newGQ(this.IsLandscape, this.Name)
		fb.Self = width == FB_GQ_H
		real_width := FB_GQ_H / float64(width) * float64(height)
		fb.Width = uint32(math.Ceil(real_width))
		result = append(result, fb)
	}

	// 适配目录
	for key, fb := range result {
		if fb.Self {
			result[key].Mp4Filepath = this.Filepath
		} else {
			output := fmt.Sprintf("%v/%v_%v", folder, fb.Resolution(), this.Name)
			result[key].Mp4Filepath = output
		}
		tmp_folder := fmt.Sprintf("%v/%v", folder, clCommon.Md5([]byte(fb.Mp4Filepath)))
		result[key].M3u8Path = tmp_folder
	}
	return result
}
