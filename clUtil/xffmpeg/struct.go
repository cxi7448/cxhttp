package xffmpeg

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clFile"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type FFmpeg struct {
	args      []string   // 运行的参数列表
	input     string     // 输入
	output    string     // 输出
	err       error      // 错误内容
	cmd       string     // 命令
	cmd_file  string     //命令文件
	out       string     // 输出结果
	videoInfo *VideoInfo // 视频信息
}

func New() *FFmpeg {
	ffmpeg := &FFmpeg{
		args: []string{},
		cmd:  "", // 默认环境变量的值
	}
	var paths []string
	var filename string
	if runtime.GOOS == "window" {
		paths = strings.Split(os.Getenv("PATH"), ";")
		filename = "ffmpeg.exe"
	} else {
		filename = "ffmpeg"
		paths = strings.Split(os.Getenv("PATH"), ":")
	}
	// 扫描当前目录是否存在执行文件
	pwd, _ := os.Getwd()
	ffmpeg_cmd := fmt.Sprintf("%v/%v", pwd, filename)
	if clFile.IsFile(ffmpeg_cmd) {
		ffmpeg.Cmd(ffmpeg_cmd)
	} else {
		for _, path := range paths {
			ffmpeg_cmd = fmt.Sprintf("%v/ffmpeg", path)
			if clFile.IsFile(ffmpeg_cmd) {
				ffmpeg.Cmd(ffmpeg_cmd)
				break
			}
		}
	}
	return ffmpeg
}

// 设置输入
func (this *FFmpeg) Input(input string) *FFmpeg {
	this.input = input
	return this
}

// 设置命令
func (this *FFmpeg) Cmd(cmd string) *FFmpeg {
	this.cmd = cmd
	return this
}

// 设置输出
func (this *FFmpeg) Output(output string) *FFmpeg {
	this.output = output
	return this
}

func (this *FFmpeg) AddArgs(args ...string) *FFmpeg {
	this.args = append(this.args, args...)
	return this
}

func (this *FFmpeg) Run() ([]byte, error) {
	var args = []string{"-y", "-i", this.input}
	args = append(args, this.args...)
	if this.output != "" {
		args = append(args, this.output)
	}
	return this.runCommand(args...)
}

func (this *FFmpeg) runCommand(args ...string) ([]byte, error) {
	var out []byte
	var err error
	clLog.Info("执行命令:%v", fmt.Sprintf("%v %v", this.cmd, strings.Join(args, " ")))
	cmd := exec.Command(this.cmd, args...)
	out, err = cmd.CombinedOutput()
	this.out = string(out)
	this.err = err
	return out, err
}

func (this *FFmpeg) CmdOutput() string {
	return this.out
}

func (this *FFmpeg) Error() error {
	return this.err
}

func (this *FFmpeg) GetPreviewImage() ([]byte, error) {
	vInfo := this.GetVideoInfo()
	duration := vInfo.GetLength()
	var count = 100
	var interval = 100 / float64(1+count)
	var timemarks = []float64{}
	for i := 0; i < count; i++ {
		timemarks = append(timemarks, interval*(float64(i)+1))
	}
	per_time := fmt.Sprintf("%0.2f", duration*(timemarks[2]-timemarks[1])/100)
	this.args = []string{"-vf", fmt.Sprintf("fps=1/%v:round=zero:start_time=0.99,scale=160:-1,tile=%vx1", per_time, count), "-frames:v", "1"}
	return this.Run()
}
