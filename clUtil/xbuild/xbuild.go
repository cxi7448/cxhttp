package xbuild

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cxi7448/cxhttp/clGlobal"
	"github.com/cxi7448/cxhttp/clUtil/clCommon"
	"github.com/cxi7448/cxhttp/clUtil/clFile"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const rootDir = "src"

func BuildRule(json_path, gomodule string) error {
	if gomodule == "" {
		return fmt.Errorf("请输入go module")
	}
	clLog.Info("开始读取postman导出的JSON文件:%v", json_path)
	content, err := ioutil.ReadFile(json_path)
	if err != nil {
		return fmt.Errorf("请输入合法的postman的jsonfile:%v", err)
	}
	rows := Item{}
	err = json.Unmarshal(content, &rows)
	if err != nil {
		return fmt.Errorf("请输入合法的postman的jsonfile:%v", err)
	}
	items := rows.GetItems("")
	fmt.Println("总共api:", len(items))

	// 扫描路由
	ruleList := scanRuleList()
	var buffer bytes.Buffer
	var apis bytes.Buffer
	var apiFiles []Api
	var v0 bytes.Buffer
	v0.WriteString("package v0\n")
	var allapis = []string{}
	for _, item := range items {
		//fmt.Println(item.Folder, item.Name)
		request := *item.Request
		requestName := ""
		if len(request.Url.Path) > 0 {
			requestName = request.Url.Path[0]
		}
		acName := ""
		if len(request.Url.Path) > 1 {
			acName = request.Url.Path[1]
		}
		if requestName == "" {
			fmt.Printf("什么鸡吧: %+v \n", item)
			fmt.Printf("什么鸡吧: %+v \n", item.Request)
			continue
		}
		if ruleList.Exists(acName) {
			continue
		}
		var apiName = fmt.Sprintf("Api%v", clCommon.ConvertToCamelCase(acName))
		if clCommon.InArray(apiName, allapis) {
			fmt.Println("重复API:", apiName)
			continue
		}
		allapis = append(allapis, apiName)
		respBody := item.Response[0].Body
		var jsonBody string
		for _, bodyRow := range strings.Split(respBody, "\n") {
			if index := strings.Index(bodyRow, "//"); index > 0 {
				jsonBody += strings.TrimSpace(bodyRow[0:index])
			} else {
				jsonBody += strings.TrimSpace(bodyRow)
			}
		}
		jsonBody = parseResponse(jsonBody)
		var outStr, structStr string
		if jsonBody != "" {
			structStr, outStr = jsonToStructStr(acName, jsonBody)
			v0.WriteString(structStr + "\n")
		}
		_content := fmt.Sprintf(`func %v(_auth *clAuth.AuthInfo, _param *rule.HttpParam, _server *rule.ServerParam) string {
			return clResponse.Success(%v)
		}`, apiName, outStr)
		//v0Input := `"xbuild/src/controller/v0"`
		//if outStr == "" {
		//	v0Input = ""
		//}
		apis.WriteString(_content)
		apiContent := fmt.Sprintf(`package controller
import (
	"github.com/cxi7448/cxhttp/clResponse"
	"github.com/cxi7448/cxhttp/core/clAuth"
	"github.com/cxi7448/cxhttp/core/rule"
)
%v
%v
`, structStr, _content)
		apiFiles = append(apiFiles, Api{
			Name:    apiName,
			Content: apiContent,
			Folder:  item.Folder,
		})
		apis.WriteString("\n")
		var params string = ""
		if len(request.Url.Query) > 0 {
			for _, query := range request.Url.Query {
				if query.Key == "" {
					continue
				}
				_type := "rule.PTYPE_SAFE_STR"
				Static := "false"
				Default := query.Value
				if strings.Contains(query.Description, "整形") || strings.Contains(query.Description, "整数") || strings.Contains(query.Description, "int") || strings.Contains(query.Description, "long") {
					_type = "rule.PTYPE_INT"
				}
				if strings.Contains(query.Description, "必填") {
					Static = "true"
				}
				if strings.Contains(query.Description, "逗号隔开") {
					_type = "rule.PTYPE_NUMBER_LIST"
				}
				if strings.Contains(query.Description, "手机号码") {
					_type = "rule.PTYPE_PHONE"
				}
				if strings.Contains(query.Description, "日期时间") {
					_type = "rule.PTYPE_DATETIME"
				} else if strings.Contains(query.Description, "日期") {
					_type = "rule.PTYPE_DATE"
				}

				if strings.Contains(query.Description, "浮点数") {
					_type = "rule.PTYPE_FLOAT"
				}
				params += fmt.Sprintf(`{Name:"%v",ParamType:%v,Static:%v,Default: "%v"},
`, query.Key, _type, Static, Default)
			}
		}
		params += ""
		buffer.WriteString(fmt.Sprintf("// %v\n", item.Name))
		buffer.WriteString(fmt.Sprintf(`rule.AddRule(rule.Rule{
		Request:     "%v",
		Name:        "%v",
		Params:      []rule.ParamInfo{
	%v
},
		CallBack:    controller.%v,
	})`, requestName, acName, params, apiName))
		buffer.WriteString("\n")
	}
	//buffer.WriteString("}\n")

	appendToRuleList(gomodule, buffer.Bytes())
	//ioutil.WriteFile(ruleFile, buffer.Bytes(), 0700)

	// build控制器
	var contrlBuffer bytes.Buffer
	contrlBuffer.WriteString(fmt.Sprintf(`package controller
import (
	"github.com/cxi7448/cxhttp/clResponse"
	"github.com/cxi7448/cxhttp/core/clAuth"
	"github.com/cxi7448/cxhttp/core/rule"
	"xbuild/src/controller/v0"
)
%v
`, apis.String()))
	controllerPath := fmt.Sprintf("%v/controller", rootDir)
	//controllerFile := fmt.Sprintf("%v/controller.go", controllerPath)
	//v0Path := fmt.Sprintf("%v/controller/v0", rootDir)
	//v0File := fmt.Sprintf("%v/controller/v0/v0.go", rootDir)
	os.MkdirAll(controllerPath, 0700)
	//ioutil.WriteFile(v0File, v0.Bytes(), 0700)
	for _, file := range apiFiles {
		//os.MkdirAll(controllerPath+"/"+file.Folder, 0700)
		filename := fmt.Sprintf("%v/%v", controllerPath, file.Path())
		fmt.Println(filename)
		if clFile.IsFile(filename) {
			continue
		}
		fmt.Println(ioutil.WriteFile(filename, []byte(file.Content), 0700))
	}
	//ioutil.WriteFile(controllerFile, contrlBuffer.Bytes(), 0700)
	return nil
}

// 自动生成结构
func BuildModel(table string) (string, error) {
	db := clGlobal.GetMysql()
	res, err := db.Query("show columns from %v", table)
	if err != nil {
		clLog.Error("错误:%v", err)
		return "", err
	}
	var result = fmt.Sprintf("type %v struct { \n", clCommon.ConvertToCamelCase(table))
	for _, val := range res.ArrResult {
		column := val.GetStr("Field", "")
		column_type := parseColumnType(val.GetStr("Type", ""))
		pri_key := parsePRI(val.GetStr("Key", ""))
		result += fmt.Sprintf("%v %v `db:\"%v\" %v json:\"%v\"`\n", clCommon.ConvertToCamelCase(column), column_type, column, pri_key, column)
	}
	result += fmt.Sprintf("}\n")
	return result, nil
}

func parsePRI(key string) string {
	if strings.Contains(key, "PRI") {
		return `primary:"true"`
	}
	return ""
}

func parseColumnType(_type string) string {
	if strings.Contains(_type, "tinyint") {
		return "uint32"
	} else if strings.Contains(_type, "int") {
		return "uint32"
	} else if strings.Contains(_type, "decimal") || strings.Contains(_type, "float") {
		return "float32"
	} else {
		return "string"
	}
}

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func jsonToStruct(jsonStr string) string {
	js := clJson.New([]byte(jsonStr))
	return js.GetKey("data").ToStr()
	//json.Unmarshal([]byte(jsonStr), &result)
	//return result
}

func parseResponse(jsStr string) string {
	js := clJson.New([]byte(jsStr))
	data := js.GetStr("data")
	if data != "null" {
		return data
	} else {
		return ""
	}
}

func jsonToStructStr(name, str string) (string, string) {
	structName := clCommon.ConvertToCamelCase(name)
	var result string = fmt.Sprintf("type %v struct {\n", structName)
	js := clJson.New([]byte(str))
	if strings.HasPrefix(str, "{") {
		// map结构
		js.ToMap().Each(func(key string, val *clJson.JsonStream) {
			result += fmt.Sprintf("%v string `json:\"%v\"`\n", clCommon.ConvertToCamelCase(key), key)
		})
	} else {
		js.ToArray().Each(func(key int, value *clJson.JsonStream) bool {
			result, structName = jsonToStructStr(name, value.ToStr())
			return false
		})
		return result, fmt.Sprintf("[]%v", structName)
	}
	result += "}"
	return result, fmt.Sprintf("%v{}", structName)
}

func scanRuleList() RuleList {
	folder := "src/rule_list"
	rows := []Rule{}
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return rows
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		content, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", folder, file.Name()))
		if err != nil {
			clLog.Error("打开文件失败:%v", err)
			continue
		}
		reg := regexp.MustCompile(`Name:\s+"[a-zA-Z0-9_]+"`)
		allRules := reg.FindAll(content, -1)
		for _, rule := range allRules {
			acName := bytes.TrimSpace(rule[5:])
			rows = append(rows, Rule{
				Name: string(acName[1 : len(acName)-1]),
			})
		}
	}
	return rows
}

func appendToRuleList(gomodule string, rule []byte) {
	rulePath := fmt.Sprintf("%v/rule_list", rootDir)
	ruleFile := fmt.Sprintf("%v/rule_list.go", rulePath)
	if clFile.IsFile(ruleFile) {
		// 存在
		content, err := ioutil.ReadFile(ruleFile)
		if err != nil {
			clLog.Error("读取:[%v] 失败:%v", ruleFile, err)
		} else {
			lastIndex := bytes.LastIndex(content, []byte("}"))
			new_content := content[0:lastIndex]
			new_content = append(new_content, []byte("\n")...)
			new_content = append(new_content, rule...)
			new_content = append(new_content, []byte("\n}")...)
			ioutil.WriteFile(ruleFile, new_content, 0700)
		}
	} else {
		// 不存在
		os.MkdirAll(rulePath, 0700)
		var new_content bytes.Buffer
		new_content.WriteString(fmt.Sprintf(`package rule_list
	import (
		"%v/src/controller"
		"github.com/cxi7448/cxhttp/core/rule"
	)`, gomodule))
		new_content.WriteString("\n")
		new_content.WriteString("func InitRules(){\n")
		new_content.Write(rule)
		new_content.WriteString("\n")
		new_content.Write([]byte("}"))
		ioutil.WriteFile(ruleFile, new_content.Bytes(), 0700)
	}
}
