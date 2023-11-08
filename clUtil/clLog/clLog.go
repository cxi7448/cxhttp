package clLog

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clFile"
	"github.com/cxi7448/cxhttp/clUtil/clTime"
	"runtime"
	"strings"
	"time"
)

var logFlag uint32 = LogFlagAll
var logMode uint32 = LogModeConsole
var version string = ""
var showTime bool = true
var logFilePath string = ""

const (
	LogFlagInfo    = 1 << 0
	LogFlagDebug   = 1 << 1
	LogFlagWarning = 1 << 2
	LogFlagError   = 1 << 3
	LogFlagAll     = LogFlagInfo | LogFlagWarning | LogFlagDebug | LogFlagError
)

const (
	LogColorOrange = uint8(iota + 91)
	LogColorGreen
	LogColorYellow
	LogColorPurple
	LogColorMagenta
	LogColorBlue
	LogColorWhite
)

const (
	// 控制台模式，会带有颜色区分
	LogModeConsole = 0
	// 原生模式，不带任何颜色信息
	LogModeRaw = 1
)

// 设置日志输出模式
func SetLogMode(_flag uint32) {
	logMode = _flag
}

// 设置日志输出掩码
func SetLogFlag(_flag uint32) {
	logFlag = _flag
}

// 设置版本号
func SetVersion(_version string) {
	version = _version
}

// 设置时间
func SetShowTime(_show bool) {
	showTime = _show
}

// 设置日志文件名
func SetLogFileName(_name string) {
	logFilePath = _name
}

// 输出Info日志
func Info(_fmt string, _args ...interface{}) {
	if logFlag&LogFlagInfo == 0 {
		return
	}
	print(LogFlagInfo, _fmt, _args...)
}

// 输出Debug日志
func Debug(_fmt string, _args ...interface{}) {
	if logFlag&LogFlagDebug == 0 {
		return
	}
	print(LogFlagDebug, _fmt, _args...)
}

// 输出Warning日志
func Warning(_fmt string, _args ...interface{}) {
	if logFlag&LogFlagWarning == 0 {
		return
	}
	print(LogFlagWarning, _fmt, _args...)
}

// 输出Error日志
func Error(_fmt string, _args ...interface{}) {
	if logFlag&LogFlagError == 0 {
		return
	}
	print(LogFlagError, _fmt, _args...)
}

// 打印
func print(_type uint32, _fmt string, _args ...interface{}) {
	var timeStr = clTime.GetDateByFormat(uint32(time.Now().Unix()), "15:04:05")
	var content = _fmt
	if _args != nil && len(_args) > 0 {
		content = fmt.Sprintf(_fmt, _args...)
	}

	debugFile := ""
	_, file, line, ok := runtime.Caller(2)
	if ok {
		fileItems := strings.Split(file, "/")
		debugFile = fmt.Sprintf("%v:%v", fileItems[len(fileItems)-1], line)
	}

	logContent := ""
	if showTime {
		logContent += timeStr
	}

	if version != "" {
		logContent += " " + version
	}
	logContent += " " + debugFile
	color := LogColorWhite
	switch _type {
	case LogFlagInfo:
		logContent += "[Info]"
	case LogFlagWarning:
		logContent += "[Warn]"
		color = LogColorYellow
	case LogFlagError:
		logContent += "[Err]"
		color = LogColorOrange
	case LogFlagDebug:
		logContent += "[Debug]"
		color = LogColorMagenta
	}

	if logMode == LogModeConsole {
		fmt.Printf("\x1b[0;%dm%v\x1b[0m\n", color, logContent+" "+content)
		return
	}
	fmt.Println(logContent + " " + content)
	if logFilePath != "" {
		clFile.AppendFile(logFilePath, logContent+" "+content)
	}
}
