package xbuild

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clGlobal"
	"github.com/cxi7448/cxhttp/clUtil/clCommon"
	"github.com/cxi7448/cxhttp/clUtil/clFile"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"os"
)

/*
*
controll： 控制器目录
model: 模型目录
*/
func BuildCURD(table, controller, model, router string) error {
	if controller == "" {
		controller = "src/controller"
	}
	if model == "" {
		model = "src/table"
	}
	if router == "" {
		router = "src/router"
	}
	controller += "/" + table
	model += "/" + table
	modelFile := fmt.Sprintf("%v/%v_model.go", model, table)
	controllerFile := fmt.Sprintf("%v/%v_api.go", controller, table)
	routerFile := fmt.Sprintf("%v/router_%v.go", router, table)
	os.MkdirAll(controller, 0700)
	os.MkdirAll(model, 0700)
	db := clGlobal.GetMysql()
	res, err := db.Query("show columns from %v", table)
	if err != nil {
		clLog.Error("错误:%v", err)
		return err
	}
	modelName := clCommon.ConvertToCamelCase(table)
	routerResult := fmt.Sprintf("package router \n")
	routerResult += fmt.Sprintf(`
import (
	"github.com/cxi7448/cxhttp/core/rule"
)
func init%v(){
`, modelName)
	controllerResult := fmt.Sprintf("package %v\n", table)
	controllerResult += fmt.Sprintf(`
import (
	"github.com/cxi7448/cxhttp/clGlobal"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/cxi7448/cxhttp/clResponse"
	"github.com/cxi7448/cxhttp/core/clAuth"
	"github.com/cxi7448/cxhttp/core/rule"
)
`)
	modelResult := fmt.Sprintf("package %v \n", table)
	modelResult += fmt.Sprintf(`
import (
	"github.com/cxi7448/cxhttp/clGlobal"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
)
`)
	modelResult += fmt.Sprintf("const Table = \"%v\" \n", table)
	modelResult += fmt.Sprintf("type %v struct { \n", modelName)
	var columnMap string = "data := clJson.M{}\n"
	for _, val := range res.ArrResult {
		column := val.GetStr("Field", "")
		column_type := parseColumnType(val.GetStr("Type", ""))
		pri_key := parsePRI(val.GetStr("Key", ""))
		modelResult += fmt.Sprintf("%v %v `db:\"%v\" %v json:\"%v\"`\n", clCommon.ConvertToCamelCase(column), column_type, column, pri_key, column)
		if column_type == "string" {
			columnMap += fmt.Sprintf("data[\"%v\"] = _param.GetStr(\"%v\",\"\")\n", column, column)
		} else {
			columnMap += fmt.Sprintf("data[\"%v\"] = _param.Get%v(\"%v\",0)\n", column, clCommon.ConvertToCamelCase(column_type), column)
		}
	}
	modelResult += fmt.Sprintf("}\n")
	modelResult += fmt.Sprintf("func Get(id uint32)*%v{\n", modelName)
	modelResult += "db := clGlobal.GetMysql()\n"
	modelResult += fmt.Sprintf("row := &%v{}\n", modelName)
	modelResult += fmt.Sprintf("err := db.NewBuilder().Table(Table).Where(\"id = %%d \",id).FindOne(row)\n")
	modelResult += "if err != nil {\n"
	modelResult += "if err.Error() != \"not found\" {\n"
	modelResult += "clLog.Error(\"错误:%v\",err)\n"
	modelResult += "}\n"
	modelResult += "return nil\n"
	modelResult += "}\n"
	modelResult += "return row\n"
	modelResult += "}\n"

	// 列表
	controllerResult += fmt.Sprintf("func Api%vList(_auth *clAuth.AuthInfo, _param *rule.HttpParam, _server *rule.ServerParam)string{\n", modelName)
	controllerResult += "pageid := _param.GetInt32(\"pageid\",0)\n"
	controllerResult += "pcount := _param.GetInt32(\"pcount\",10)\n"
	controllerResult += fmt.Sprintf("rows :=[]%v.%v{}\n", table, modelName)
	controllerResult += "where := \"1 = 1\"\n"
	controllerResult += "db := clGlobal.GetMysql()\n"
	controllerResult += fmt.Sprintf(`
	total, err := db.NewBuilder().List(%v.Table, where, pageid, pcount, &rows, "id desc")
	if err != nil && err.Error() != "not found" {
		clLog.Error("错误:%%v", err)
	}
	return clResponse.Success(clJson.M{
		"list":  rows,
		"total": total,
	})
`, table)
	controllerResult += "\n"
	controllerResult += "}\n"

	routerResult += fmt.Sprintf(`

rule.AddRule(rule.Rule{
		Request: "request",
		Name:    "%v_list",
		Params: []rule.ParamInfo{
		},
		CallBack: %v.Api%vList,
		Method:   "POST",
	})

`, table, table, modelName)

	// 添加
	controllerResult += fmt.Sprintf(`
func Api%vAdd(_auth *clAuth.AuthInfo, _param *rule.HttpParam, _server *rule.ServerParam)string{
%v
db := clGlobal.GetMysql()
_, err := db.NewBuilder().Table(%v.Table).Add(data)
if err != nil {
	clLog.Error("错误:%%v",err)
	return clResponse.Error("添加失败:%%v",err)
}
return clResponse.Success()
}
`, modelName, columnMap, table)

	routerResult += fmt.Sprintf(`

rule.AddRule(rule.Rule{
		Request: "request",
		Name:    "%v_add",
		Params: []rule.ParamInfo{
		},
		CallBack: %v.Api%vAdd,
		Method:   "POST",
	})

`, table, table, modelName)
	// 编辑
	controllerResult += fmt.Sprintf(`
func Api%vEdit(_auth *clAuth.AuthInfo, _param *rule.HttpParam, _server *rule.ServerParam)string{
id := _param.GetUint32("id",0)
%v
db := clGlobal.GetMysql()
_, err := db.NewBuilder().Table(%v.Table).Where("id = %%d",id).Save(data)
if err != nil {
	clLog.Error("错误:%%v",err)
	return clResponse.Error("编辑失败:%%v",err)
}
return clResponse.Success()
}
`, modelName, columnMap, table)

	routerResult += fmt.Sprintf(`

rule.AddRule(rule.Rule{
		Request: "request",
		Name:    "%v_edit",
		Params: []rule.ParamInfo{
		},
		CallBack: %v.Api%vEdit,
		Method:   "POST",
	})

`, table, table, modelName)

	// 删除
	controllerResult += fmt.Sprintf(`
func Api%vDel(_auth *clAuth.AuthInfo, _param *rule.HttpParam, _server *rule.ServerParam)string{
ids := _param.GetStr("ids","")
db := clGlobal.GetMysql()
_, err := db.NewBuilder().Table(%v.Table).Where("id in(%%v)",ids).Del()
if err != nil {
	clLog.Error("错误:%%v",err)
	return clResponse.Error("编辑失败:%%v",err)
}
return clResponse.Success()
}
`, modelName, table)

	routerResult += fmt.Sprintf(`

rule.AddRule(rule.Rule{
		Request: "request",
		Name:    "%v_delete",
		Params: []rule.ParamInfo{
			{Name: "ids", ParamType: rule.PTYPE_NUMBER_LIST, Static: true},
		},
		CallBack: %v.Api%vDel,
		Method:   "POST",
	})

`, table, table, modelName)

	// 创建模型文件
	if !clFile.IsFile(modelFile) {
		// 自动生成，存在就不生成了
		os.WriteFile(modelFile, []byte(modelResult), 0700)
	}

	// 添加
	// 编辑
	// 删除
	// 创建控制器文件
	if !clFile.IsFile(controllerFile) {
		// 自动生成，存在就不生成了
		os.WriteFile(controllerFile, []byte(controllerResult), 0700)
	}

	routerResult += "\n}"
	// 创建模型文件
	if !clFile.IsFile(routerFile) {
		// 自动生成，存在就不生成了
		os.WriteFile(routerFile, []byte(routerResult), 0700)
	}

	return nil
}
