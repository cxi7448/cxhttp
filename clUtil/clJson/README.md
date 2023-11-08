# clJson包
> 此扩展包建立初衷是为了更方便以及更高效的处理Json字符串,它也确实做到了.它主要实现了流方式解析json字符串为json对象。以及非常方便的由 `map[string]interface{}` 到json字符串之间的互相转化.


那么如何将一个json字符串转化为json对象呢?
```go
jsonStr := `{"a":1, "b":2}`
jsObj := clJson.New( []byte(jsonStr) )
```

使用`New`就能非常方便的将字符串转化为一个jsonObject, 同理，我们也可以通过如下方式获得一个空的json对象:
```go
jsObj := clJson.New( []byte("{}") )
```

好，得到了JsonObject我们就可以通过各种方式来操作它了:
```
// 获取a的值
jsObj.GetInt32("a")
```
除了Int32以外，还提供了非常多的数据类型可供选择:
```
// 获取字符串类型
jsObj.GetStr("a")
// 获取Float32类型
jsObj.GetFloat32("a")
// 获取Int32类型
jsObj.GetInt32("a")
// 获取Int64类型
jsObj.GetInt64("a")
// 获取Uint64类型
jsObj.GetUint64("a")
// 获取Uint32类型
jsObj.GetUint32("a")
// 获取Float64类型
jsObj.GetFloat64("a")
// 获取Bool类型
jsObj.GetBool("a")
```

如果是嵌套结构，获取下一级的值可以如下操作
```
var jsonStr = `{"a": {"b": 100}}
var jsonObj = clJson.New( []byte( jsonStr ) )
jsonObj.GetInt32("a", "b")
```
其他类型以此类推

假如要获取下级所有数据，可以如下:
```
var jsonStr = `{"a": {"b": 100, "c": 200}}
var jsonObj = clJson.New( []byte( jsonStr ) )
var jsonMap = jsonObj.GetMap("a")
```
得到的是一个jsonMap的对象，然后可以对他进行类似map的操作
```
// 获取字符串类型
jsonMap.GetStr("b")
// 获取Float32类型
jsonMap.GetFloat32("b")
// 获取Int32类型
jsonMap.GetInt32("b")
// 获取Int64类型
jsonMap.GetInt64("b")
// 获取Uint64类型
jsonMap.GetUint64("b")
// 获取Uint32类型
jsonMap.GetUint32("b")
// 获取Float64类型
jsonMap.GetFloat64("b")
// 获取Bool类型
jsonMap.GetBool("b")
// 转化为一般的 map[string]string 对象
jsonMap.ToCustom()
// 转化为一般的 map[string]interface{} 对象
jsonMap.ToTree()
// 转化为string
jsonMap.ToStr()
```

如果获取数组也可以这样:
```
// 获取数组类型
var jsonStr = `{"a": [ 100, 200, 300 ]}
var jsonObj = clJson.New( []byte( jsonStr ) )
var jsonArr = jsonObj.GetArray("a")

// 获取数组第一个元素
jsonArr.GetOffset(0)

// 判断数组是否为空
jsonArr.IsEmpty()

// 获取数组长度
jsonArr.GetLength()

// 转化为标准的[]string数组
jsonArr.ToCustom()

// 转化为[]*JsonStream
jsonArr.ToTree()

// 遍历jsonArr
jsonObj := New([]byte( `{"a":[1,2,3,4,5,6]}`))
if jsonObj == nil {
    fmt.Printf("jsonObj解析错误!\n")
    return
}

jsonArr := jsonObj.GetArray("a")
jsonArr.Each(func(key int, value *JsonStream) bool {
    fmt.Printf("val: %v\n", value.ToStr())
    // 返回true继续循环，返回false结束循环，相当于break
    return true
})
```