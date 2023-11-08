package clJson

import (
	"fmt"
	"testing"
)

func TestCreateBy(t *testing.T) {
	jsonstr := `{
		"ac": "getSystemInfo"
	}`
	jsonObj := New([]byte(jsonstr))
	if jsonObj == nil {
		fmt.Printf("jsonObj解析错误!\n")
		return
	}
	fmt.Printf(">> ac的值: %v\n", jsonObj.GetStr("ac"))
}

func TestJsonStream_GetArray(t *testing.T) {

	jsonObj := New([]byte(`{"a":[1,2,3,4,5,6]}`))
	if jsonObj == nil {
		fmt.Printf("jsonObj解析错误!\n")
		return
	}

	jsonArr := jsonObj.GetArray("a")
	jsonArr.Each(func(key int, value *JsonStream) bool {
		fmt.Printf("val: %v\n", value.ToStr())
		return true
	})

	fmt.Printf("数组: %+v\n", jsonArr.ToCustom())
}

type TestCode struct {
	Code  uint32 `json:"code"`
	Param string `json:"param"`
}

func TestJsonArray_ToCustom(t *testing.T) {

	//fmt.Printf("Test: %v\n", strconv.QuoteToASCII("服务器繁忙"))
	jsonObj := New([]byte(`{"code":1,"param":"\u670d\u52a1\u5668\u7e41\u5fd9"}`))
	//if jsonObj == nil {
	//	return
	//}
	fmt.Printf("字符串1: %v\n", jsonObj.ToStr())
	//fmt.Printf("字符串2: %v\n", string(jsonObj2))
	//
	fmt.Printf("jsonObj: %+v\n", jsonObj.ToMap().ToCustom())
}
