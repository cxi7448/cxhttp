package rule

import (
	"encoding/json"
	"fmt"
	"github.com/cxi7448/cxhttp/clCommon"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"strconv"
	"strings"
)

type HttpParam struct {
	values map[string]string
	path   Path
}

// @author cxhttp
// @lastUpdate 2019-08-10
// @comment 添加一个参数到参数列表中
// @param _key 参数的名称
// @param _val 参数的值
func (this *HttpParam) Add(_key, _val string) {
	this.values[_key] = _val
}

func NewHttpParam(_params map[string]string, path []string) *HttpParam {
	result := &HttpParam{
		values: make(map[string]string),
		path:   []string{},
	}
	if _params != nil {
		result.values = _params
	}
	if len(path) > 0 {
		result.path = path
	}
	return result
}

// @author cxhttp
// @lastUpdate 2019-08-10
// @comment 获取一个字符串型的参数
// @param _key 要获取的参数名称
// @param _default 如果key不存在, 默认返回什么
func (this *HttpParam) GetStr(_key string, _default string) string {
	val, exists := this.values[_key]
	if !exists {
		return _default
	}
	return val
}

// @author cxhttp
// @lastUpdate 2019-08-10
// @comment 获取uint32类型的参数
// @param _key 要获取的参数名称
// @param _default 如果key不存在, 默认返回什么
func (this *HttpParam) GetUint32(_key string, _default uint32) uint32 {

	val, exists := this.values[_key]
	if !exists {
		return _default
	}

	i, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return _default
	}
	return uint32(i)
}

// @author cxhttp
// @lastUpdate 2019-08-10
// @comment 获取uint64类型的参数
// @param _key 要获取的参数名称
// @param _default 如果key不存在, 默认返回什么
func (this *HttpParam) GetUint64(_key string, _default uint64) uint64 {

	val, exists := this.values[_key]
	if !exists {
		return _default
	}

	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return _default
	}
	return uint64(i)
}

// @author cxhttp
// @lastUpdate 2019-08-10
// @comment 获取int32类型的参数
// @param _key 要获取的参数名称
// @param _default 如果key不存在, 默认返回什么
func (this *HttpParam) GetInt32(_key string, _default int32) int32 {

	val, exists := this.values[_key]
	if !exists {
		return _default
	}

	i, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return _default
	}
	return int32(i)
}

// @author cxhttp
// @lastUpdate 2019-08-10
// @comment 获取int64类型的参数
// @param _key 要获取的参数名称
// @param _default 如果key不存在, 默认返回什么
func (this *HttpParam) GetInt64(_key string, _default int64) int64 {

	val, exists := this.values[_key]
	if !exists {
		return _default
	}

	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return _default
	}
	return int64(i)
}

// @author cxhttp
// @lastUpdate 2019-08-10
// @comment 获取32位浮点数
// @param _key 要获取的参数名称
// @param _default 如果key不存在，默认返回什么
func (this *HttpParam) GetFloat32(_key string, _default float32) float32 {

	val, exists := this.values[_key]
	if !exists {
		return _default
	}

	i, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return _default
	}
	return float32(i)
}

// @author cxhttp
// @lastUpdate 2019-08-10
// @comment 获取64位浮点数
// @param _key 要获取的参数名称
// @param _default 如果key不存在，默认返回什么
func (this *HttpParam) GetFloat64(_key string, _default float64) float64 {

	val, exists := this.values[_key]
	if !exists {
		return _default
	}

	i, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return _default
	}
	return float64(i)
}

// @author cxhttp
// @lastUpdate 2019-08-10
// @comment 获取浮点数类型
// @param _key 要获取的参数名称
// @param _default 如果不存在默认返回什么
func (this *HttpParam) GetBool(_key string, _default bool) bool {

	val, exists := this.values[_key]
	if !exists {
		return _default
	}

	switch strings.ToUpper(val) {
	case "OK", "ON", "YES", "TRUE", "Y", "T":
		return true
	}

	return false
}

// @author cxhttp
// @lastUpdate 2022-09-15
// @comment 将参数根据指定字符切割后返回
// @param _key 要获得的参数名称
// @param _sep 进行分割的符号
func (this *HttpParam) GetStrSplit(_key string, _sep string) []string {
	val, exists := this.values[_key]
	if !exists {
		return []string{}
	}

	return strings.Split(val, _sep)
}

// @author cxhttp
// @lastUpdate 2022-09-15
// @comment 获取整数列表
// @param _key 要获得的参数名称
func (this *HttpParam) GetUint32Split(_key string) []uint32 {
	val, exists := this.values[_key]
	if !exists || val == "" {
		return []uint32{}
	}

	strArr := strings.Split(val, ",")
	uint32Arr := make([]uint32, 0)
	for _, strItem := range strArr {
		uint32Arr = append(uint32Arr, clCommon.Uint32(strItem))
	}

	return uint32Arr
}

// @author cxhttp
// @lastUpdate 2022-09-15
// @comment 获取浮点数列表
// @param _key 要获得的参数名称
func (this *HttpParam) GetFloatSplit(_key string) []float64 {
	val, exists := this.values[_key]
	if !exists {
		return []float64{}
	}

	strArr := strings.Split(val, ",")
	float64Arr := make([]float64, 0)
	for _, strItem := range strArr {
		float64Arr = append(float64Arr, clCommon.Float64(strItem))
	}

	return float64Arr
}

func (this *HttpParam) GetObj(_key string, value interface{}) error {
	val, exists := this.values[_key]
	if !exists {
		return fmt.Errorf(fmt.Sprintf("KEY %v NOT EXISTS", _key))
	}
	return json.Unmarshal([]byte(val), value)
}

// @author cxhttp
// @lastUpdate 2023-02-23
// @comment 返回所有的参数
func (this *HttpParam) ToMap() map[string]string {
	return this.values
}

// @author cxhttp
// @lastUpdate 2023-02-23
// @comment 将所有参数以json字符串形式返回
func (this *HttpParam) ToJson() string {
	return clJson.CreateBy(this.values).ToStr()
}

func (this *HttpParam) Path() Path {
	return this.path
}
