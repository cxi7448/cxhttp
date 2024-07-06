package xbuild

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clGlobal"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
)

type Table struct {
	Name    string
	Columns []TableColumn
}

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

func (this *Table) GetFormStr() string {
	var result = ""
	for _, val := range this.Columns {
		if val.Type == "string" {
			result += fmt.Sprintf("\"%v\":\"%v\",\n", val.Field, val.Default)
		} else {
			if val.Default == "" {
				val.Default = 0
			}
			result += fmt.Sprintf("\"%v\":%v,\n", val.Field, val.Default)
		}
	}
	return result
}

func (this *Table) ElTableColumn() string {
	var result = ""
	for _, val := range this.Columns {
		if val.Field == "status" {
			result += fmt.Sprintf("<el-table-column  prop=\"%v\" label=\"%v\" align=\"center\" :show-overflow-tooltip=\"true\" >\n<template #default=\"scope\"><el-switch :active-value=\"1\" :inactive-value=\"0\" v-model=\"scope.row.%v\"/></template>\n</el-table-column>\n", val.Field, val.Comment, val.Field)
		} else {
			result += fmt.Sprintf("<el-table-column  prop=\"%v\" label=\"%v\" align=\"center\" :show-overflow-tooltip=\"true\" />\n", val.Field, val.Comment)
		}
	}
	return result
}

func (this *Table) ElFormItem() string {
	var result = ""
	for _, val := range this.Columns {
		if val.Field == "status" || val.Field == "is_hot" {
			result += fmt.Sprintf("<el-form-item label=\"%v\" label-width=\"80px\"><el-switch :active-value=\"1\" :inactive-value=\"0\" v-model=\"form.%v\"/></el-form-item>\n", val.Comment, val.Field)
		} else if val.Type != "string" {
			result += fmt.Sprintf("<el-form-item label=\"%v\" label-width=\"80px\"><el-input class=\"width_400\" type=\"number\" placeholder=\"请输入%v\" v-model=\"form.%v\"></el-input></el-form-item>\n", val.Comment, val.Comment, val.Field)
		} else {
			result += fmt.Sprintf("<el-form-item label=\"%v\" label-width=\"80px\"><el-input class=\"width_400\" placeholder=\"请输入%v\" v-model=\"form.%v\"></el-input></el-form-item>\n", val.Comment, val.Comment, val.Field)
		}
	}
	return result
}
