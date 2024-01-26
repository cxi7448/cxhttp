package xbuild

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cxi7448/cxhttp/clGlobal"
	"github.com/cxi7448/cxhttp/clUtil/clCommon"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"io/ioutil"
	"os"
	"strings"
)

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
	items := rows.GetItems()
	fmt.Println("总共api:", len(items))
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf(`package rule_list
import (
	"%v/src/controller"
	"github.com/cxi7448/cxhttp/core/rule"
)`, gomodule))
	buffer.WriteString("\n")
	buffer.WriteString("func InitRules(){\n")
	var apis bytes.Buffer
	var allapis = []string{}
	for _, item := range items {
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
		var apiName = fmt.Sprintf("Api%v", clCommon.ConvertToCamelCase(acName))
		if clCommon.InArray(apiName, allapis) {
			fmt.Println("重复API:", apiName)
			continue
		}
		allapis = append(allapis, apiName)
		apis.WriteString(fmt.Sprintf(`func %v(_auth *clAuth.AuthInfo, _param *rule.HttpParam, _server *rule.ServerParam) string {
			return clResponse.Success()
		}`, apiName))
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
	buffer.WriteString("}\n")
	rootDir := "src"
	rulePath := fmt.Sprintf("%v/rule", rootDir)
	ruleFile := fmt.Sprintf("%v/rule.go", rulePath)
	os.MkdirAll(rulePath, 0700)
	ioutil.WriteFile(ruleFile, buffer.Bytes(), 0700)

	// build控制器
	var contrlBuffer bytes.Buffer
	contrlBuffer.WriteString(fmt.Sprintf(`package controller
import (
	"github.com/cxi7448/cxhttp/clResponse"
	"github.com/cxi7448/cxhttp/core/clAuth"
	"github.com/cxi7448/cxhttp/core/rule"
	v0 "xbuild/src/controller/v0"
)
%v
`, apis.String()))
	controllerPath := fmt.Sprintf("%v/controller", rootDir)
	controllerFile := fmt.Sprintf("%v/controller.go", controllerPath)
	os.MkdirAll(controllerPath, 0700)
	ioutil.WriteFile(controllerFile, contrlBuffer.Bytes(), 0700)
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
