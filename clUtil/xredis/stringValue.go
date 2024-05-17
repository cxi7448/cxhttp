package xredis

import (
	"encoding/json"
	"github.com/cxi7448/cxhttp/clCommon"
	"time"
)

type StringValue struct {
	Value  string
	Expire int64
}

func (this *StringValue) ToFloat64() float64 {
	res := this.ToString()
	return clCommon.Float64(res)
}
func (this *StringValue) ToFloat32() float32 {
	res := this.ToString()
	return clCommon.Float32(res)
}
func (this *StringValue) ToInt32() int32 {
	res := this.ToString()
	return clCommon.Int32(res)
}

func (this *StringValue) ToInt() int {
	res := this.ToString()
	return clCommon.Int(res)
}
func (this *StringValue) ToInt64() int64 {
	res := this.ToString()
	return clCommon.Int64(res)
}

func (this *StringValue) ToUint32() uint32 {
	res := this.ToString()
	return clCommon.Uint32(res)
}

func (this *StringValue) ToUint() uint {
	res := this.ToString()
	return uint(clCommon.Uint32(res))
}
func (this *StringValue) ToUint64() uint64 {
	res := this.ToString()
	return clCommon.Uint64(res)
}

func (this *StringValue) ToString() string {
	return this.Value
}

func (this *StringValue) ToObj(value interface{}) error {
	if this.Expire > 0 && this.Expire < time.Now().Unix() {
		return nil
	}
	return json.Unmarshal([]byte(this.Value), value)
}

// 判断是否过期
func (this *StringValue) IsExpire() bool {
	return this.Expire > 0 && this.Expire < time.Now().Unix()
}

var defaultStringValue = &StringValue{}
