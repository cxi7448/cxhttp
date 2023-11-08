package clJson

import (
	"regexp"
	"strconv"
	"strings"
)

// 遍历数组
func (js *JsonMap) Each(hf func(key string, val *JsonStream)) {
	for key, val := range *js {
		hf(key, val.data)
	}
}

// 获取指定的对象
func (js *JsonMap) Find(key string) *JsonStream {
	if val, Ok := (*js)[key]; Ok {
		return val.data
	}
	return nil
}

// key是否存在
func (js *JsonMap) IsSet(_key string) bool {
	_, ok := (*js)[_key]
	return ok
}

// 重置为标准库的map[string] interface{} 去掉了反斜杠
func (js *JsonMap) ToCustom() map[string]string {
	var CustomMap = make(map[string]string)
	reg, _ := regexp.Compile(`(^\\)*\\\\(^\\)*`)
	for key, val := range *js {
		if val.data == nil {
			continue
		}
		res := reg.ReplaceAllString(val.data.ToStr(), "_replace_xiegang_") // 换掉 \\
		res = strings.ReplaceAll(res, "\\u", "_replace_chinese_")          // 换掉\n
		res = strings.ReplaceAll(res, "\\n", "_replace_huanhang_")         // 换掉\n
		res = strings.ReplaceAll(res, "\\", "")                            // 换掉剩余的\
		res = strings.ReplaceAll(res, "_replace_xiegang_", "\\")           // 还原斜杠
		res = strings.ReplaceAll(res, "_replace_huanhang_", "\\n")         // 还原换行
		res = strings.ReplaceAll(res, "_replace_chinese_", "\\u")          // 还原中文
		CustomMap[key] = res
	}
	return CustomMap
}

// 重置为标准库的map[string] string
func (js *JsonMap) Tostring() map[string]string {
	var CustomMap = make(map[string]string)
	for key, val := range *js {
		CustomMap[key] = val.data.ToStr()
	}
	return CustomMap
}

// 重置为标准库的map[string] interface{}
func (js *JsonMap) ToTree() map[string]interface{} {
	var CustomMap = make(map[string]interface{})
	for key, val := range *js {
		CustomMap[key] = val.StackParseTree()
	}
	return CustomMap
}

// 递归查找树结构
func (this *jsonItem) StackParseTree() interface{} {
	if this.data.dataType == JSON_TYPE_MAP {
		this_map := this.data.ToMap()
		var customMap = make(map[string]interface{})
		for key, val := range *this_map {
			customMap[key] = val.StackParseTree()
		}
		return customMap
	} else if this.data.dataType == JSON_TYPE_ARR {
		this_arr := this.data.ToArray()
		var customArr = make([]interface{}, 0)
		for _, val := range *this_arr {
			customArr = append(customArr, val.StackParseTree())
		}
		return customArr
	} else if this.data.dataType == JSON_TYPE_BOOL {
		b, e := strconv.ParseBool(this.data.ToStr())
		if e != nil {
			return false
		}
		return b
	} else if this.data.dataType == JSON_TYPE_INT {
		i, e := strconv.ParseFloat(this.data.ToStr(), 64)
		if e != nil {
			return 0
		}
		return i
	} else if this.data.dataType == JSON_TYPE_NULL {
		return nil
	}
	return this.data.ToStr()
}

// 获取指定下标并转换成string类型
// @param key string 下标
// @param def string 默认值
// @return string 返回指定的之
func (js *JsonMap) GetStr(key string, def string) string {
	if val, ok := (*js)[key]; ok {
		return val.data.ToStr()
	}
	return ""
}

// 获取指定下标并转换成前后去空的string类型
// @param key string 下标
// @param def string 默认值
// @return string 返回指定的之
func (js *JsonMap) GetStrTrim(key string, def string) string {
	if val, ok := (*js)[key]; ok {
		return strings.TrimSpace(val.data.ToStr())
	}
	return ""
}

// 获取指定下标并转换成uint32类型
func (js *JsonMap) GetUint32(key string, def uint32) uint32 {
	if val, ok := (*js)[key]; ok {
		b, e := strconv.ParseUint(val.data.ToStr(), 10, 32)
		if e != nil {
			return def
		}
		return uint32(b)
	}
	return def
}

// 获取指定下标并转换成uint32类型
func (js *JsonMap) GetUint64(key string, def uint64) uint64 {
	if val, ok := (*js)[key]; ok {
		b, e := strconv.ParseUint(val.data.ToStr(), 10, 64)
		if e != nil {
			return def
		}
		return b
	}
	return def
}

// 获取指定下标并转换成int32类型
func (js *JsonMap) GetInt32(key string, def int32) int32 {
	if val, ok := (*js)[key]; ok {
		b, e := strconv.ParseInt(val.data.ToStr(), 10, 32)
		if e != nil {
			return def
		}
		return int32(b)
	}
	return def
}

// 获取指定下标并转换成float32类型
func (js *JsonMap) GetFloat32(key string, def float32) float32 {
	if val, ok := (*js)[key]; ok {
		b, e := strconv.ParseFloat(val.data.ToStr(), 32)
		if e != nil {
			return def
		}
		return float32(b)
	}
	return def
}

// 获取指定下标并转换成float64类型
func (js *JsonMap) GetFloat64(key string, def float64) float64 {
	if val, ok := (*js)[key]; ok {
		b, e := strconv.ParseFloat(val.data.ToStr(), 64)
		if e != nil {
			return def
		}
		return b
	}
	return def
}

// 获取指定下标并转换成bool类型
func (js *JsonMap) GetBool(key string, def bool) bool {
	if val, ok := (*js)[key]; ok {
		b, e := strconv.ParseBool(val.data.ToStr())
		if e != nil {
			return def
		}
		return b
	}
	return def
}

func (js *JsonMap) DelKey(key string) *JsonMap {
	if _, ok := (*js)[key]; ok {
		delete(*js, key)
		return js
	}
	return nil
}

func (js *JsonMap) IsEmpty() bool {
	if len(*js) == 0 {
		return true
	}
	return false
}
