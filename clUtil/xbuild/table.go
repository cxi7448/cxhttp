package xbuild

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clGlobal"
	"github.com/cxi7448/cxhttp/clUtil/clCommon"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
)

type Table struct {
	Name    string
	Columns []TableColumn
}

var UploadColumns = []string{"icon", "image"}

type TableColumn struct {
	Field   string
	Type    string
	Primary string
	Default interface{}
	Comment string
}

func GenTable(table string) *Table {
	row := &Table{
		Name:    table,
		Columns: []TableColumn{},
	}
	db := clGlobal.GetMysql()
	res, err := db.Query("show full columns  from %v", table)
	if err != nil {
		clLog.Error("错误:%v", err)
		return nil
	}

	for _, val := range res.ArrResult {
		column := TableColumn{
			Field:   val.GetStr("Field", ""),
			Type:    parseColumnType(val.GetStr("Type", "")),
			Primary: parsePRI(val.GetStr("Key", "")),
			Default: val.GetStr("Default", ""),
			Comment: val.GetStr("Comment", ""),
		}
		if column.Comment == "" {
			column.Comment = column.Field
		}
		row.Columns = append(row.Columns, column)
	}
	return row
}

// 生成文件上传js代码
func (this *Table) GenScript(islang bool) string {
	// 自动检测是否存在图片相关字段: icon,image
	//for _,row := range this.Columns {
	//
	//}
	var result = ""
	if islang {
		result += "import language from '@/lang'\n"
	}
	// import {fetchImage} from "@/util/func.ts";
	uploadColumns := this.GetUploadColumns()
	if len(uploadColumns) > 0 {
		result += "import {fetchImage} from '@/util/func.ts'\n"
		for _, val := range uploadColumns {
			field := clCommon.ConvertToCamelCase(val)
			result += fmt.Sprintf(`
const beforeUpload%v = (file:any)=>{
  console.log("上传文件：%%v",file)
  return true
}
const startUpload%v = (param:any)=>{
  var headers = {
    'Name':"tcp"
  }
  proxy?.$http.uploadFile(param.file,headers).then((res:any)=>{
    if (res.code == 0){
      form.value.%v = res.data.filename
    }else{
      ElMessage({message:res.msg,type:"error"})
    }
  })
}
`, field, field, val)
			result += "\n"
			result += "\n"
		}
	}
	return result
}

func (this *Table) GetUploadColumns() []string {
	result := []string{}
	for _, val := range this.Columns {
		if clCommon.InArray(val.Field, UploadColumns) {
			result = append(result, val.Field)
		}
	}
	return result
}

func (this *Table) GenLanguageImport(islang bool) string {
	if islang {
		return `import language from '@/lang'`
	}
	return ""
}

func (this *Table) GetFormStr(islang bool) string {
	var result = ""
	for _, val := range this.Columns {
		if val.Type == "string" {
			result += fmt.Sprintf("%v:\"%v\",\n", val.Field, val.Default)
		} else {
			if val.Default == "" {
				val.Default = 0
			}
			result += fmt.Sprintf("%v:%v,\n", val.Field, val.Default)
		}
	}
	if islang {
		result += `
lang:{
	CN:"",
  EN:"",
  THA:""
},
`
	}
	return result
}

func (this *Table) ElTableColumn() string {
	var result = ""
	for _, val := range this.Columns {
		if val.Field == "status" {
			result += fmt.Sprintf("<el-table-column  prop=\"%v\" label=\"状态\" align=\"center\" :show-overflow-tooltip=\"true\" >\n<template #default=\"scope\"><el-switch disabled :active-value=\"1\" :inactive-value=\"0\" v-model=\"scope.row.%v\"/></template>\n</el-table-column>\n", val.Field, val.Field)
		} else if clCommon.InArray(val.Field, UploadColumns) {
			result += fmt.Sprintf("<el-table-column  prop=\"%v\" label=\"%v\" align=\"center\" :show-overflow-tooltip=\"true\" >\n      <template #default=\"scope\">\n        <el-image :src=\"fetchImage(scope.row.%v)\"></el-image>\n      </template>\n    </el-table-column>\n", val.Field, val.Comment, val.Field)
		} else {
			result += fmt.Sprintf("<el-table-column  prop=\"%v\" label=\"%v\" align=\"center\" :show-overflow-tooltip=\"true\" />\n", val.Field, val.Comment)
		}
	}
	return result
}

func (this *Table) ElFormItem(islang bool) string {
	var result = ""
	for _, val := range this.Columns {
		is_add := ""
		if val.Field == "id" {
			is_add = ` v-if="!isAdd" `
		}
		if val.Field == "status" || val.Field == "is_hot" {
			result += fmt.Sprintf("<el-form-item %v label=\"状态\" ><el-switch :active-value=\"1\" :inactive-value=\"0\" v-model=\"form.%v\"/></el-form-item>\n", is_add, val.Field)
		} else if val.Type != "string" {
			result += fmt.Sprintf("<el-form-item %v label=\"%v\" ><el-input class=\"width_400\" type=\"number\" placeholder=\"请输入%v\" v-model=\"form.%v\"></el-input></el-form-item>\n", is_add, val.Comment, val.Comment, val.Field)
		} else if clCommon.InArray(val.Field, UploadColumns) {
			field := clCommon.ConvertToCamelCase(val.Field)
			result += fmt.Sprintf("<el-form-item  label=\"%v\"><el-input class=\"width_300\" placeholder=\"请输入%v\" v-model=\"form.%v\">\n</el-input><el-upload ref=\"upload%v\"\n                      class=\"upload-demo\"\n                      style=\"height: 32px;\"\n                      accept=\"image/*\"\n                      action=\"#\"\n                      :show-file-list=\"false\"\n                      :before-upload=\"beforeUpload%v\"\n                      :http-request=\"startUpload%v\"\n                      :auto-upload=\"true\">\n  <template #trigger>\n    <el-button type=\"primary\">选择文件</el-button>\n  </template>\n</el-upload>\n  <div>\n    <el-image :src=\"fetchImage(form.%v)\" v-if=\"form.%v\" style=\"max-width: 200px; max-height: 200px\"></el-image>\n  </div>\n</el-form-item>\n", val.Comment, val.Comment, val.Field, field, field, field, val.Field, val.Field)
			// :src="fetchImage(form.icon)" v-if="form.icon"
		} else {
			result += fmt.Sprintf("<el-form-item %v label=\"%v\" ><el-input class=\"width_400\" placeholder=\"请输入%v\" v-model=\"form.%v\"></el-input></el-form-item>\n", is_add, val.Comment, val.Comment, val.Field)
		}
	}
	if islang {
		result += fmt.Sprint(`
<el-form-item label="语言包">
        <el-input v-for="lang of language" class="width_400" v-model="form['lang'][lang.key]" :placeholder="'请输入'+lang.title">
          <template #prepend>{{ lang.key }}</template>
        </el-input>
      </el-form-item>
`)
		result += "\n"
	}
	return result
}
