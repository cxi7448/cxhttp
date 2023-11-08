package clConfig

import (
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/cxi7448/cxhttp/clUtil/clSuperMap"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var mConfig map[string]*clSuperMap.SuperMap
var aConfig map[string][]string

func init() {
	mConfig = make(map[string]*clSuperMap.SuperMap)
	aConfig = make(map[string][]string)
}

// 通过文件加载配置
// 调用多次后面的会覆盖前面
// @Param _filename 配置文件所在路径
// @Param _overWrite 是否覆盖之前读取的内容， true为覆盖，反之不覆盖
func LoadFromFile(_filename string, _overWrite bool) error {
	buffer, err := ioutil.ReadFile(_filename)
	if err != nil {
		return err
	}

	section := "global"
	if _overWrite {
		mConfig = make(map[string]*clSuperMap.SuperMap)
		mConfig[section] = clSuperMap.NewSuperMap()
		aConfig = make(map[string][]string)
		aConfig[section] = make([]string, 0)
	}

	var confLines = strings.Split(string(buffer), "\n")
	for _, line := range confLines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "--") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimPrefix(strings.TrimSuffix(line, "]"), "[")
			mConfig[section] = clSuperMap.NewSuperMap()
			aConfig[section] = make([]string, 0)
			continue
		}

		idx := strings.Index(line, "=")
		if idx <= 0 {
			// 可能是数组
			aConfig[section] = append(aConfig[section], line)
			continue
		}
		clLog.Info("section: %v key: %v, val: %v", section, line[:idx], line[idx+1:])
		mConfig[section].Add(line[:idx], line[idx+1:])
	}
	return nil
}

// 获取Uint32配置
func GetUint32(_key string, _def uint32) uint32 {
	if os.Getenv(_key) != "" {
		v, e := strconv.ParseInt(os.Getenv(_key), 10, 32)
		if e != nil {
			return _def
		}
		return uint32(v)
	}
	idx := strings.Index(_key, ".")
	if idx <= 0 {
		return _def
	}
	configItem, exists := mConfig[_key[:idx]]
	if !exists {
		return _def
	}
	return configItem.GetUInt32(_key[idx+1:], _def)
}

// 获取Uint64配置
func GetUint64(_key string, _def uint64) uint64 {
	if os.Getenv(_key) != "" {
		v, e := strconv.ParseInt(os.Getenv(_key), 10, 64)
		if e != nil {
			return _def
		}
		return uint64(v)
	}
	idx := strings.Index(_key, ".")
	if idx <= 0 {
		return _def
	}
	configItem, exists := mConfig[_key[:idx]]
	if !exists {
		return _def
	}
	return configItem.GetUInt64(_key[idx+1:], _def)
}

// 获取int32配置
func GetInt32(_key string, _def int32) int32 {
	if os.Getenv(_key) != "" {
		v, e := strconv.ParseInt(os.Getenv(_key), 10, 32)
		if e != nil {
			return _def
		}
		return int32(v)
	}
	idx := strings.Index(_key, ".")
	if idx <= 0 {
		return _def
	}
	configItem, exists := mConfig[_key[:idx]]
	if !exists {
		return _def
	}
	return configItem.GetInt32(_key[idx+1:], _def)
}

// 获取int64配置
func GetInt64(_key string, _def int64) int64 {
	if os.Getenv(_key) != "" {
		v, e := strconv.ParseInt(os.Getenv(_key), 10, 64)
		if e != nil {
			return _def
		}
		return int64(v)
	}
	idx := strings.Index(_key, ".")
	if idx <= 0 {
		return _def
	}
	configItem, exists := mConfig[_key[:idx]]
	if !exists {
		return _def
	}
	return configItem.GetInt64(_key[idx+1:], _def)
}

// 获取bool配置
func GetBool(_key string, _def bool) bool {
	if os.Getenv(_key) != "" {
		v, e := strconv.ParseBool(os.Getenv(_key))
		if e != nil {
			return _def
		}
		return v
	}
	idx := strings.Index(_key, ".")
	if idx <= 0 {
		return _def
	}
	configItem, exists := mConfig[_key[:idx]]
	if !exists {
		return _def
	}
	return configItem.GetBool(_key[idx+1:], _def)
}

// 获取Str配置
func GetStr(_key string, _def string) string {
	if os.Getenv(_key) != "" {
		return os.Getenv(_key)
	}
	idx := strings.Index(_key, ".")
	if idx <= 0 {
		return _def
	}
	configItem, exists := mConfig[_key[:idx]]
	if !exists {

		return _def
	}
	return configItem.GetStr(_key[idx+1:], _def)
}

// 获取Float32配置
func GetFloat32(_key string, _def float32) float32 {
	if os.Getenv(_key) != "" {
		v, e := strconv.ParseFloat(os.Getenv(_key), 32)
		if e != nil {
			return _def
		}
		return float32(v)
	}
	idx := strings.Index(_key, ".")
	if idx <= 0 {
		return _def
	}
	configItem, exists := mConfig[_key[:idx]]
	if !exists {
		return _def
	}
	return configItem.GetFloat32(_key[idx+1:], _def)
}

// 获取Float32配置
func GetFloat64(_key string, _def float64) float64 {
	if os.Getenv(_key) != "" {
		v, e := strconv.ParseFloat(os.Getenv(_key), 64)
		if e != nil {
			return _def
		}
		return v
	}
	idx := strings.Index(_key, ".")
	if idx <= 0 {
		return _def
	}
	configItem, exists := mConfig[_key[:idx]]
	if !exists {
		return _def
	}
	return configItem.GetFloat64(_key[idx+1:], _def)
}

// 获取Arr配置
func GetStrArray(_key string) []string {
	configItem, exists := aConfig[_key]
	if !exists {
		return []string{}
	}
	return configItem
}

// 获取Arr配置
func GetMap(_key string) map[string]string {
	configItem, exists := mConfig[_key]
	if !exists {
		return map[string]string{}
	}
	return configItem.GetMap()
}
