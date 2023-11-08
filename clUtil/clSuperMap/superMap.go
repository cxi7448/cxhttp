package clSuperMap

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// 对Map进行扩展，实现多种数据类型的存储和转换
const ValueNotExists = "THE_VALUE_IS_NOT_EXISTS"

// 创建一个超级Map
func NewSuperMap() *SuperMap {
	return &SuperMap{
		data:   make(map[string]string),
		locker: sync.RWMutex{},
	}
}


// 使用map生成一个超级map
func NewSuperMapByTree(_map map[string]interface{}) *SuperMap {
	var superMap = &SuperMap{
		data:   make(map[string]string),
		locker: sync.RWMutex{},
	}
	for key, val := range _map {
		superMap.Add(key, fmt.Sprintf("%v", val))
	}
	return superMap
}

// 使用map生成一个超级map
func NewSuperMapByMap(_map map[string]string) *SuperMap {
	var superMap = &SuperMap{
		data:   make(map[string]string),
		locker: sync.RWMutex{},
	}
	for key, val := range _map {
		superMap.Add(key, val)
	}
	return superMap
}

// 添加一个值
func (this *SuperMap) Add(_key string, _val interface{}) {
	this.locker.Lock()
	defer this.locker.Unlock()

	this.data[_key] = fmt.Sprintf("%v", _val)
}


// 判断某个key是否存在
func (this *SuperMap) IsExists(_key string) bool {
	this.locker.Lock()
	defer this.locker.Unlock()
	_, exists := this.data[_key]
	return exists
}


// 获取String类型
func (this *SuperMap) GetStr(_key string, _default string) string {
	this.locker.RLock()
	defer this.locker.RUnlock()

	val, exists := this.data[_key]
	if !exists {
		return _default
	}
	return val
}


// 获取Int类型
func (this *SuperMap) GetInt(_key string, _default int) int {
	var val = this.GetStr(_key, ValueNotExists)
	if val == ValueNotExists {
		return _default
	}
	b, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return _default
	}
	return int(b)
}


// 获取Int32类型
func (this *SuperMap) GetInt32(_key string, _default int32) int32 {
	var val = this.GetStr(_key, ValueNotExists)
	if val == ValueNotExists {
		return _default
	}
	b, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return _default
	}
	return int32(b)
}


// 获取Int64类型
func (this *SuperMap) GetInt64(_key string, _default int64) int64 {
	var val = this.GetStr(_key, ValueNotExists)
	if val == ValueNotExists {
		return _default
	}
	b, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return _default
	}
	return int64(b)
}


// 获取Uint类型
func (this *SuperMap) GetUInt(_key string, _default uint) uint {
	var val = this.GetStr(_key, ValueNotExists)
	if val == ValueNotExists {
		return _default
	}
	b, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return _default
	}
	return uint(b)
}


// 获取Uint32类型
func (this *SuperMap) GetUInt32(_key string, _default uint32) uint32 {
	var val = this.GetStr(_key, ValueNotExists)
	if val == ValueNotExists {
		return _default
	}
	b, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return _default
	}
	return uint32(b)
}


// 获取Uint64类型
func (this *SuperMap) GetUInt64(_key string, _default uint64) uint64 {
	var val = this.GetStr(_key, ValueNotExists)
	if val == ValueNotExists {
		return _default
	}
	b, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return _default
	}
	return uint64(b)
}

// 获取Float32
func (this *SuperMap) GetFloat32(_key string, _default float32) float32 {
	var val = this.GetStr(_key, ValueNotExists)
	if val == ValueNotExists {
		return _default
	}
	b, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return _default
	}
	return float32(b)
}


// 获取Float64
func (this *SuperMap) GetFloat64(_key string, _default float64) float64 {
	var val = this.GetStr(_key, ValueNotExists)
	if val == ValueNotExists {
		return _default
	}
	b, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return _default
	}
	return b
}


// 获取Bool类型
func (this *SuperMap) GetBool(_key string, _default bool) bool {
	var val = this.GetStr(_key, ValueNotExists)
	if val == ValueNotExists {
		return _default
	}
	b, err := strconv.ParseBool( val )
	if err != nil {
		return _default
	}
	return b
}


// 使用指定字符串进行分割
func (this *SuperMap) SplitBy(_key string, _sep string) []string {
	var val = this.GetStr(_key, ValueNotExists)
	if val == ValueNotExists {
		return []string{}
	}
	return strings.Split(val, _sep)
}



// 使用某个正则进行测试, 如果测试成功返回true，测试失败或者key不存在，返回false
func (this *SuperMap) RegMatch(_key string, _exp string) bool {
	var val = this.GetStr(_key, ValueNotExists)
	if val == ValueNotExists {
		return false
	}
	match, err := regexp.Match(_exp, []byte(val))
	if err != nil {
		return false
	}
	return match
}


// 获取整个Map
func (this *SuperMap) GetMap() map[string]string {
	this.locker.RLock()
	defer this.locker.RUnlock()

	return this.data
}


// 获取按顺序排列好的key
func (this *SuperMap) GetSortKeys() []string {
	this.locker.RLock()
	defer this.locker.RUnlock()

	var keyList = make([]string, 0)
	sort.Strings(keyList)
	return keyList
}



// 获取按Key的ASCII码顺序遍历Map
func (this *SuperMap) ForeachByKeySort(_func func(k, v string) bool ) {
	this.locker.RLock()
	defer this.locker.RUnlock()

	if _func == nil {
		return
	}

	var keyList = make([]string, 0)
	sort.Strings(keyList)

	for _, k := range keyList {
		isContinue := _func( k, this.GetStr(k, "") )
		if !isContinue {
			break
		}
	}
	return
}