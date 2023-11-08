# ClSuperMap

> clSuperMap是为了强化 `map[string]string` 一些编程中常用的方法而产生的。它提供了一系列方便的方式维护`map[string]string` 并且它是线程安全的，内部集成异步锁进行读写保护。

主要的用法也非常简单, 要生成一个 `clSuperMap` 主要通过两种方式
```go
    smap := superMap.NewSuperMap()
```
或者也可以用一个现有的 `map[string]interface{}` 对象直接创建它:
```go
   imap = map[string] interface{}{
       "key1": 100,
       "key2": false,
   }
   smap := superMap.NewSuperMapByMap( imap )
```

得到一个 `clSuperMap` 对象后就可以非常容易的使用它, 假设要探测某个key是否存在，可以:
```go
    fmt.Printf("key是否存在: %v\n", imap.IsExists("key1"))
```

如果要添加也很简单:
```go
    smap.Add("key1", "value1")
```
要注意的是Add是可以支持添加任意类型的, 所以可以直接添加一个int, bool 或者float类型:
```go
    smap.Add("key2", 0)
    smap.Add("Key3", true)
    smap.Add("key4", 1.33)
```


还可以通过如下方法快捷的获取指定类型的值，并可以设置默认值，以便在获取不到的时候有个备用处理方案:
```go

// 获取字符串类型
smap.GetStr("key1", "def")
// 获取int类型
smap.GetInt("key1", 0)
// 获取int32类型
smap.GetInt32("key1", 0)
// 获取int64类型
smap.GetInt64("key1", 0)
// 获取uint类型
smap.GetUint("key1", 0)
// 获取uint32类型
smap.GetUint32("key1", 0)
// 获取uint64类型
smap.GetUint64("key1", 0)
// 获取Boolean类型
smap.GetBoolean("key1", 0)
// 获取Float32类型
smap.GetFloat32("key1", 0)
// 获取Float64类型
smap.GetFloat64("key1", 0)
// 获取一个字符串列表，经常有个场景，某个key中存储的是一个字符串的列表，用逗号或者竖线隔开。
// 以前的操作要进行两部:
// 1. 先读取出这个key的值
// 2. 然后使用strings.Split方法进行切割
// 为了提高效率，提供了如下方法，可以一步到位, 这样就可以一步将key1的值使用","分割后返回分割后的结果
smap.SplitBy("key1", ",")


// 检测key1的值是否符合指定正则
// 同样为了方便提供这个方法简化代码量
smap.RegMatch("key1", "/^[0-9]{1,3}$/")

// 如果要返回整个map可以使用如下方式：
// 需要注意的是返回的结果是map[string]string 而非 map[string]interface{}
smap.GetMap()

// 由于golang的hash map是无序的，如果需要返回一个按照Ascii码升序排列的key的数组
// 因为很多签名校验算法都需要它。所以你也可以很方便的使用如下方法进行获取
smap.GetSortKeys()
```

