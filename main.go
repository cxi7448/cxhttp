package main

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clGlobal"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/cxi7448/cxhttp/clUtil/xbuild"
	"os"
)

// HTTP服务默认使用端口号
const HTTPServerPort = 19999

func main() {
	//clLog.SetLogFlag(0)
	clGlobal.Init("./cl.conf")
	var module = "table"
	if len(os.Args) < 2 {
		fmt.Println(desc)
		return
	}
	module = os.Args[1]
	switch module {
	case "table":
		buildTable()
	case "rule":
		buildRule()
	case "api":
		buildApi()
	case "vue":
		buildVue()
	default:
		fmt.Println(desc)
	}
}

const desc = `使用帮助
	路由模式生成： xbuild rule postman导出.json  gomodule

	postman的Description(参数描述)关键字设置参数类型:
	1.整形,整数,int,long  		参数设置为整数
	2.必填  				 		参数设置为必填
	3.逗号隔开					参数设置为多个整数逗号隔开:  1,2,3
	4.手机号码 					参数设置为手机号码
	5.日期时间					参数设置为日期时间: 2024-01-01 01:10:10				
	6.日期						参数设置为日期(优先级比日期时间小): 2024-01-01
	7.浮点数						参数设置为浮点数(float)

	postman的Value(参数value)里面的值，将会被读取为参数的默认值
	
	示例:
	接口名字自动登陆
	POST: http://localhost:17777/api/UserLogin
	postman的QueryParams
	key				value				Description
	account								手机号码,必填
	text								
	vcode								整数,必填
	ids									逗号隔开
	type			0					整数,类型

	自动生成的rule
	// 自动登录/游客登陆
	rule.AddRule(rule.Rule{
		Request: "api",
		Name:    "UserLogin",
		Params: []rule.ParamInfo{
			{Name: "account",ParamType:rule.PTYPE_PHONE,Static:true,Default:""},
			{Name: "text",ParamType:rule.PTYPE_SAFE_STR,Static:false,Default:""},
			{Name: "vcode",ParamType:rule.PTYPE_INT,Static:true,Default:""},
			{Name: "ids",ParamType:rule.PTYPE_NUMBER_LIST,Static:false,Default:""},
			{Name: "type",ParamType:rule.PTYPE_INT,Static:false,Default:"0"},
		},
		CallBack: controller.ApiUserLogin,
	})
`

// 自动生成Index.vue
func buildVue() {
	if len(os.Args) < 3 {
		fmt.Println("请输入表名")
		return
	}
	table := os.Args[2]
	path := ""
	if len(os.Args) > 3 {
		path = os.Args[3]
	}
	lang := ""
	if len(os.Args) > 4 {
		lang = os.Args[4]
	}
	err := xbuild.BuildView(table, path, lang)
	if err != nil {
		fmt.Println("生成错误:", err)
	} else {
		fmt.Println("生成完毕")
	}
}

// 自动生成CURD
func buildApi() {
	if len(os.Args) < 3 {
		fmt.Println("请输入表名")
		return
	}
	table := os.Args[2]
	controller := ""
	model := ""
	router := ""
	if len(os.Args) > 3 {
		controller = os.Args[3]
	}
	if len(os.Args) > 4 {
		model = os.Args[4]
	}
	if len(os.Args) > 5 {
		router = os.Args[5]
	}
	err := xbuild.BuildCURD(table, controller, model, router)
	if err != nil {
		fmt.Println("生成错误:", err)
	} else {
		fmt.Println("生成完毕")
	}
}

func buildTable() {
	for _, arg := range os.Args[2:] {
		if arg == "" {
			continue
		}
		result, err := xbuild.BuildModel(arg)
		if err != nil {
			clLog.Error("错误:%v", err)
			continue
		}
		fmt.Printf("const Table = \"%v\"\n", arg)
		fmt.Println(result)
		fmt.Println()
	}
}

func buildRule() {
	if len(os.Args) < 3 {
		fmt.Println("请输入postman导出的json文件名字")
		return
	}
	if len(os.Args) < 4 {
		fmt.Println("请输入您的go mudole名字")
		return
	}
	jsonFile := os.Args[2]
	gomodule := os.Args[3]
	err := xbuild.BuildRule(jsonFile, gomodule)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("xbuild成功")
	}
}
