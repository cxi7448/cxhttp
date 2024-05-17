package xredis

import (
	"encoding/json"
	"fmt"
	"github.com/cxi7448/cxhttp/clCommon"
	"time"
)

func (this *SetValue) ToFloat64() float64 {
	res := this.ToString()
	return clCommon.Float64(res)
}
func (this *SetValue) ToFloat32() float32 {
	res := this.ToString()
	return clCommon.Float32(res)
}
func (this *SetValue) ToInt32() int32 {
	res := this.ToString()
	return clCommon.Int32(res)
}

func (this *SetValue) ToInt() int {
	res := this.ToString()
	return clCommon.Int(res)
}
func (this *SetValue) ToInt64() int64 {
	res := this.ToString()
	return clCommon.Int64(res)
}

func (this *SetValue) ToUint32() uint32 {
	res := this.ToString()
	return clCommon.Uint32(res)
}

func (this *SetValue) ToUint() uint {
	res := this.ToString()
	return uint(clCommon.Uint32(res))
}
func (this *SetValue) ToUint64() uint64 {
	res := this.ToString()
	return clCommon.Uint64(res)
}

func (this *SetValue) ToString() string {
	if this.Expire > 0 && this.Expire < time.Now().Unix() {
		return ""
	}
	return fmt.Sprint(this.Value)
}

func (this *SetValue) ToObj(value interface{}) error {
	if this.Expire > 0 && this.Expire < time.Now().Unix() {
		return nil
	}
	content, err := json.Marshal(this.Value)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, value)
}

// 判断是否过期
func (this *SetValue) IsExpire() bool {
	return this.Expire > 0 && this.Expire < time.Now().Unix()
}

var defaultSetValue = &SetValue{
	Expire: 1,
}
