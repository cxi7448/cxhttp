# clHttpClient 包

> 这是一个功能非常全面的Http请求客户端，可以满足各种场景下的需求.

要创建一个Http客户端非常简单，只需要提供一个Url即可，如:
```go
var httpClient = clHttpClient.NewClient("http://www.baidu.com")
```

当然，一般我们需要模拟请求各种场景下, 模拟User-Agent将成为不可缺少的一步：
```go
// 设置为360SE浏览器
httpClient.SetUA( clHttpClient.UA360SE )
```
我们内置提供了非常多种浏览器的UA可供选择，如下:
```go
const (
	UAClHttpClient = "clHttpClient/1.0"
	UASafariMac = "Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50"
	UASafariWin = "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50"
	UAIE90 = "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0"
	UAIE80 = "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0)"
	UAIE70 = "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0)"
	UAIE60 = " Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)"
	UAFireFoxMac = " Mozilla/5.0 (Macintosh; Intel Mac OS X 10.6; rv:2.0.1) Gecko/20100101 Firefox/4.0.1"
	UAFireFoxWin = "Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1"
	UAOperaMac = "Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11"
	UAOperaWin = "Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11"
	UAChromeMac = " Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11"
	UASougou = "?Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; SE 2.X MetaSr 1.0; SE 2.X MetaSr 1.0; .NET CLR 2.0.50727; SE 2.X MetaSr 1.0)"
	UA360SE = " Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)"
	UASafariIOS = "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5"
	UASafariIpodTouch = "Mozilla/5.0 (iPod; U; CPU iPhone OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5"
	UASafariIpad = "Mozilla/5.0 (iPad; U; CPU OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5"
	UAAndroidN1 = "Mozilla/5.0 (Linux; U; Android 2.3.7; en-us; Nexus One Build/FRF91) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"
	UAAndroidQQ = "MQQBrowser/26 Mozilla/5.0 (Linux; U; Android 2.3.7; zh-cn; MB200 Build/GRJ22; CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"
	UAOperaAndroid = "Opera/9.80 (Android 2.3.4; Linux; Opera Mobi/build-1107180945; U; en-GB) Presto/2.8.149 Version/11.10"
	UABlackBerry = "Mozilla/5.0 (BlackBerry; U; BlackBerry 9800; en) AppleWebKit/534.1+ (KHTML, like Gecko) Version/6.0.0.337 Mobile Safari/534.1+"
)
```

设置好了UA我们可以通过如下方式设置请求方式, 默认请求方式为POST
```go
httpClient.SetMethod("POST")
```

如果需要设置超时时间，则: (单位秒)
```go
// 设置10秒请求超时
httpClient.SetTimeout(10)
```

请求数据格式可用如下方式设置:
```go
httpClient.SetContentType(clHttpClient.ContentParam)
```
除此之外我们还提供了多种数据格式:
```go
const (
	ContentTypeForm = 0			// 正常form提交 [默认]
	ContentParam = 1			// 参数提交(只允许GET)
	ContentJson = 2				// 通过json提交
	ContentXml = 3				// 通过xml提交
	ContentXWWWFormUrl = 4		// 通过x-www-from-urlencode方式提交
)
```

如果我们需要设置代理服务器，则可以:
```go
httpClient.SetProxy("http://www.qq.com")
```

如果我们需要自定义头部:
```go
httpClient.AddHeader("referrer", "http://www.baidu.com")
```

最重要的提交参数:
```go
httpClient.AddParam("key1", "value1")
```
所有类型的数据格式均使用此方法提交参数，如果是json类型数据格式，则会自动转化为json字符串后再提交，xml也类似

当然总有一些情况此方法无法满足，没有关系，也可以通过如下方式自定义请求体
```go
httpClient.SetBody("hello world")
```
需要注意的是，上方的方式如果被设置了，则`AddParam`方法设置的所有数据均会失效!

如此一来，我们的准备工作就完成了，现在开始正式对服务器发起请求吧:
```go
	resp, err := httpClient.Do()
	if err != nil {
		fmt.Printf("请求错误: %v\n", err)
		return
	}

	fmt.Printf("请求状态码: %v\n", resp.StatusCode)
	fmt.Printf("请求回应体: %v\n", resp.Body)
	fmt.Printf("请求体长度: %v\n", resp.ContentLength)
	fmt.Printf("请求头部: %+v\n", resp.Header)
```
得到的resp结构体中包含四组数据，分别是状态码，回应数据，回应数据长度，回应头部信息等, 可分别进行操作。
如果回传的是json结构，可以通过一个方法快速的将json转化为Json结构体:

```go
type StructResp struct {
	Code uint32 `json:"code"`
	Msg string `json:"msg"`
	Data string `json:"data"`
}

var respObj = StructResp{}
resp.ToObject(&respObj)
```