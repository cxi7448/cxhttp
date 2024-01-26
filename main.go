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
	clLog.SetLogFlag(0)
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
	default:
		fmt.Println(desc)
	}

	//clGlobal.Init("cl.conf")
	//
	//rule_list.Init()
	//rule_list.InitSuperAPI()
	//
	//clAuth.SetAuthPrefix("U_INFO")
	//
	//httpserver.SetAESKey("5d41402abc4b2a76b9719d911017c592")
	//// 关闭上传功能
	//httpserver.SetEnableUploadFile(false)
	//// 关闭上传调试页
	//httpserver.SetEnableUploadTest(false)
	//
	//clLog.Info("正在启动服务，端口: %v", HTTPServerPort)
	//clLog.Info("可尝试使用 http://localhost:%v 访问", HTTPServerPort)
	//clAuth.SetGetUserByDB(func(_uid uint64) *clAuth.AuthInfo {
	//	return &clAuth.AuthInfo{
	//		Uid:        1,
	//		Token:      "1000",
	//		LastUptime: 0,
	//		IsLogin:    true,
	//		ExtraData:  nil,
	//	}
	//})
	//clGlobal.SkyConf.DebugRouter = true
	//httpserver.SetUploadFileSizeLimit(1024 * 1024 * 300)
	//
	//// 根据路由配置表生成api文档
	////rule.ApiGeneral("./apis", "apis", "/request")
	//
	//// 根据数据库中的配置生成模型
	////modelCreator.CreateAllModelFile("127.0.0.1", "root", "root", "testdb", "testModel")
	//
	//httpserver.StartServer(HTTPServerPort)
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

func buildTable() {
	for _, arg := range os.Args[2:] {
		if arg == "" {
			continue
		}
		result, err := xbuild.BuildModel(arg)
		if err != nil {
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
