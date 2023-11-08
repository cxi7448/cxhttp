package clCache

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clCommon"
	"github.com/cxi7448/cxhttp/clUtil/clFile"
	"io/ioutil"
	"strings"
	"sync"
)

type FileCache struct {
	fileName   string
	data       map[string]string
	dataLocker sync.RWMutex
}

// 新建一个缓存对象
func New(_fileName string) *FileCache {
	return &FileCache{
		fileName: _fileName,
		data:     make(map[string]string),
	}
}

// 保存到磁盘
func (this *FileCache) SaveTo() {
	var fileBuffer = strings.Builder{}
	for k, v := range this.data {
		fileBuffer.WriteString(fmt.Sprintf("%v=%v\n", k, v))
	}
	ioutil.WriteFile(this.fileName, []byte(fileBuffer.String()), 0644)
}

// 移除key
func (this *FileCache) DelKey(_key string) {
	var fullMap = this.GetFullMap()

	this.dataLocker.Lock()
	defer this.dataLocker.Unlock()
	delete(fullMap, _key)
	this.data = fullMap
	this.SaveTo()
}

// 移除keys
func (this *FileCache) DelKeys(_keys []string) {
	var fullMap = this.GetFullMap()

	this.dataLocker.Lock()
	defer this.dataLocker.Unlock()
	for _, k := range _keys {
		delete(fullMap, k)
	}
	this.data = fullMap
	this.SaveTo()
}

// 清空
func (this *FileCache) Clear() {
	this.dataLocker.Lock()
	defer this.dataLocker.Unlock()

	clFile.DelFile(this.fileName)
}

// 设置值
func (this *FileCache) SetItem(_key string, _val interface{}) {
	var fullMap = this.GetFullMap()

	this.dataLocker.Lock()
	defer this.dataLocker.Unlock()
	fullMap[_key] = fmt.Sprintf("%v", _val)
	this.data = fullMap
	this.SaveTo()
}

// 获取整个map
func (this *FileCache) GetFullMap() map[string]string {
	this.dataLocker.RLock()
	defer this.dataLocker.RUnlock()

	buffer, err := ioutil.ReadFile(this.fileName)
	if err != nil {
		return map[string]string{}
	}

	var dataMap = make(map[string]string)
	bufferLines := strings.Split(string(buffer), "\n")
	for _, val := range bufferLines {
		var idx = strings.Index(val, "=")
		if idx < 0 {
			continue
		}
		dataMap[val[:idx]] = val[idx+1:]
	}
	return dataMap
}

// 获取值
func (this *FileCache) GetStr(_key string) string {
	this.dataLocker.RLock()
	defer this.dataLocker.RUnlock()

	buffer, err := ioutil.ReadFile(this.fileName)
	if err != nil {
		return ""
	}

	bufferLines := strings.Split(string(buffer), "\n")
	for _, val := range bufferLines {
		var idx = strings.Index(val, "=")
		if idx < 0 {
			continue
		}
		if val[:idx] == _key {
			return val[idx+1:]
		}
	}
	return ""
}

// 获取整数型
func (this *FileCache) GetInt32(_key string) int32 {
	var str = this.GetStr(_key)
	if str == "" {
		return 0
	}
	return clCommon.Int32(str)
}

// 获取整数型
func (this *FileCache) GetUint32(_key string) uint32 {
	var str = this.GetStr(_key)
	if str == "" {
		return 0
	}
	return clCommon.Uint32(str)
}

// 获取整数型
func (this *FileCache) GetInt64(_key string) int64 {
	var str = this.GetStr(_key)
	if str == "" {
		return 0
	}
	return clCommon.Int64(str)
}

// 获取整数型
func (this *FileCache) GetUint64(_key string) uint64 {
	var str = this.GetStr(_key)
	if str == "" {
		return 0
	}
	return clCommon.Uint64(str)
}

// 获取浮点数
func (this *FileCache) GetFloat64(_key string) float64 {
	var str = this.GetStr(_key)
	if str == "" {
		return 0
	}
	return clCommon.Float64(str)
}

func (this *FileCache) GetFloat32(_key string) float32 {
	var str = this.GetStr(_key)
	if str == "" {
		return 0
	}
	return clCommon.Float32(str)
}

// 获取布尔
func (this *FileCache) GetBool(_key string) bool {
	var str = this.GetStr(_key)
	if str == "" {
		return false
	}
	return clCommon.Bool(str)
}

// 获取列表
func (this *FileCache) GetSplit(_key string, _sep string) []string {
	var str = this.GetStr(_key)
	if str == "" {
		return []string{}
	}
	return strings.Split(str, _sep)
}
