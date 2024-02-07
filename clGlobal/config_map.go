package clGlobal

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clCommon"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"sync"
)

var mConfigMap map[string]string
var mLocker sync.RWMutex

func init() {
	mConfigMap = make(map[string]string)
}

// 加载配置
func LoadConfig(_section, _key, _default string) string {
	mLocker.Lock()
	defer mLocker.Unlock()

	var temp string
	if conf == nil {
		clLog.Error("无法找到配置指针,请先执行Init指定配置文件")
		return ""
	}
	conf.GetStr(_section, _key, _default, &temp)
	mConfigMap[_section+"_"+_key] = temp
	return temp
}

// 获取配置
func GetConfig(_section, _key string, _default ...string) string {
	__default := ""
	if len(_default) > 0 {
		__default = _default[0]
	}
	mLocker.RLock()
	defer mLocker.RUnlock()

	val, exists := mConfigMap[_section+"_"+_key]
	if !exists {
		var temp string
		conf.GetStr(_section, _key, __default, &temp)
		mConfigMap[_section+"_"+_key] = temp
		return temp
	}
	return val
}

// 获取配置
func GetUint32(_section, _key string, _default ...uint32) uint32 {
	__default := uint32(0)
	if len(_default) > 0 {
		__default = _default[0]
	}
	val := GetConfig(_section, _key, fmt.Sprint(__default))
	if val == "" {
		return __default
	}
	return clCommon.Uint32(val)
}

// 强制指定配置
func SetConfig(_section, _key, _val string) {
	mLocker.Lock()
	defer mLocker.Unlock()

	mConfigMap[_section+"_"+_key] = _val
}
