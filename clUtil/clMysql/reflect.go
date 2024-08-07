package clMysql

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clSuperMap"
	"reflect"
	"strings"
)

// 根据数据结构取得字段列表
func GetAllField(_val interface{}) []string {

	_fields := make([]string, 0)
	_type := reflect.TypeOf(_val).Elem()

	for i := 0; i < _type.NumField(); i++ {
		if _type.Field(i).Anonymous {
			// 匿名字段,说明是嵌套的结构
			_value := reflect.ValueOf(_val)
			t := _value.Elem().Field(i)
			for j := 0; j < _value.Elem().Field(i).NumField(); j++ {
				fieldName := t.Type().Field(j).Tag.Get("db")
				if fieldName == "" || fieldName == "-" {
					continue
				}
				alias := t.Type().Field(j).Tag.Get("alias")
				if alias != "" {
					fieldName = fmt.Sprintf("%v`.`%v", alias, fieldName)
				}
				_fields = append(_fields, fieldName)
			}
			continue
		}
		fieldName := _type.Field(i).Tag.Get("db")
		if fieldName == "" || fieldName == "-" {
			continue
		}
		alias := _type.Field(i).Tag.Get("alias")
		if alias != "" {
			fieldName = fmt.Sprintf("%v`.`%v", alias, fieldName)
		}
		_fields = append(_fields, fieldName)
	}
	return _fields
}

// 将数据反序列化为一个对象
func Unmarsha(_row *clSuperMap.SuperMap, _inter interface{}) {

	_type := reflect.TypeOf(_inter)
	_value := reflect.ValueOf(_inter)
	_valueE := _value.Elem()
	for i := 0; i < _value.Elem().NumField(); i++ {
		if _type.Elem().Field(i).Anonymous {
			t := _value.Elem().Field(i)
			for j := 0; j < _value.Elem().Field(i).NumField(); j++ {
				setValue(t.Type().Field(j), _value.Elem().Field(i).Field(j), _row)
			}
			continue
		}
		//field_name := _type.Elem().Field(i).Tag.Get("db")
		setValue(_type.Elem().Field(i), _valueE.Field(i), _row)
		//switch _type.Elem().Field(i).Type.String() {
		//case "uint32":
		//	_valueE.Field(i).SetUint(uint64(_row.GetUInt32(field_name, 0)))
		//case "uint64":
		//	_valueE.Field(i).SetUint(_row.GetUInt64(field_name, 0))
		//case "int32":
		//	_valueE.Field(i).SetInt(int64(_row.GetInt32(field_name, 0)))
		//case "int64":
		//	_valueE.Field(i).SetInt(_row.GetInt64(field_name, 0))
		//case "string":
		//	_valueE.Field(i).SetString(_row.GetStr(field_name, ""))
		//case "bool":
		//	_valueE.Field(i).SetBool(_row.GetBool(field_name, false))
		//case "float32":
		//	_valueE.Field(i).SetFloat(float64(_row.GetFloat32(field_name, 0)))
		//case "float64":
		//	_valueE.Field(i).SetFloat(_row.GetFloat64(field_name, 0))
		//}
	}
}

func setValue(_type reflect.StructField, _valueE reflect.Value, _row *clSuperMap.SuperMap) {
	field_name := _type.Tag.Get("db")
	switch _valueE.Type().String() {
	case "uint32":
		_valueE.SetUint(uint64(_row.GetUInt32(field_name, 0)))
	case "uint64":
		_valueE.SetUint(_row.GetUInt64(field_name, 0))
	case "int32":
		_valueE.SetInt(int64(_row.GetInt32(field_name, 0)))
	case "int64":
		_valueE.SetInt(_row.GetInt64(field_name, 0))
	case "string":
		_valueE.SetString(_row.GetStr(field_name, ""))
	case "bool":
		_valueE.SetBool(_row.GetBool(field_name, false))
	case "float32":
		_valueE.SetFloat(float64(_row.GetFloat32(field_name, 0)))
	case "float64":
		_valueE.SetFloat(_row.GetFloat64(field_name, 0))
	}
}

// 根据数据结构取得字段列表
func GetInsertSql(_val interface{}, _primary bool) ([]string, []string) {

	_fields := make([]string, 0)
	_values := make([]string, 0)

	_type := reflect.TypeOf(_val)
	_value := reflect.ValueOf(_val)

	for i := 0; i < _type.NumField(); i++ {

		if _type.Field(i).Name[0] >= 97 {
			continue // 小写的过滤掉
		}

		if !_primary {
			if strings.ToUpper(_type.Field(i).Tag.Get("primary")) == "TRUE" {
				continue
			}
		}

		if _type.Field(i).Tag.Get("db") == "" {
			continue
		}

		_fields = append(_fields, _type.Field(i).Tag.Get("db"))
		if strings.HasPrefix(_value.Field(i).Type().String(), "float") {
			_values = append(_values, fmt.Sprintf("%f", _value.Field(i).Float()))
		} else {
			_values = append(_values, fmt.Sprintf("%v", _value.Field(i).Interface()))
		}

	}
	return _fields, _values
}

// 根据数据结构取得字段列表
func GetInsertSqlMulti(_dataList []interface{}, _primary bool) ([]string, []string) {

	_fields := make([]string, 0)
	_valueList := make([]string, 0)

	for k, _val := range _dataList {

		_type := reflect.TypeOf(_val)
		_value := reflect.ValueOf(_val)

		_values := make([]string, 0)

		for i := 0; i < _type.NumField(); i++ {

			if _type.Field(i).Name[0] >= 97 {
				continue // 小写的过滤掉
			}

			if !_primary {
				if strings.ToUpper(_type.Field(i).Tag.Get("primary")) == "TRUE" {
					continue
				}
			}

			if _type.Field(i).Tag.Get("db") == "" {
				continue
			}

			if k == 0 {
				_fields = append(_fields, _type.Field(i).Tag.Get("db"))
			}
			if strings.HasPrefix(_value.Field(i).Type().String(), "float") {
				_values = append(_values, fmt.Sprintf("%f", _value.Field(i).Float()))
			} else {
				_values = append(_values, fmt.Sprintf("%v", _value.Field(i).Interface()))
			}
		}
		_valueList = append(_valueList, fmt.Sprintf("(%v)", strings.Join(_values, ",")))
	}

	return _fields, _valueList
}

// 根据数据结构取得字段列表
func GetInsertSqlMultiNew(_dataList []interface{}, _primary bool) ([]string, []string, []interface{}) {

	_fields := make([]string, 0)
	_prepares := make([]string, 0)
	_valueList := make([]interface{}, 0)

	for k, _val := range _dataList {

		_type := reflect.TypeOf(_val)
		_value := reflect.ValueOf(_val)
		tmp_prepare := []string{}
		//_values := make([]string, 0)
		for i := 0; i < _type.NumField(); i++ {

			if _type.Field(i).Name[0] >= 97 {
				continue // 小写的过滤掉
			}

			if !_primary {
				if strings.ToUpper(_type.Field(i).Tag.Get("primary")) == "TRUE" {
					continue
				}
			}

			if _type.Field(i).Tag.Get("db") == "" {
				continue
			}

			if k == 0 {
				_fields = append(_fields, _type.Field(i).Tag.Get("db"))
			}
			if strings.HasPrefix(_value.Field(i).Type().String(), "float") {
				_valueList = append(_valueList, fmt.Sprintf("%f", _value.Field(i).Float()))
			} else {
				_valueList = append(_valueList, fmt.Sprintf("%v", _value.Field(i).Interface()))
			}
			tmp_prepare = append(tmp_prepare, "?")

		}
		_prepares = append(_prepares, fmt.Sprintf("(%v)", strings.Join(tmp_prepare, ",")))
	}

	return _fields, _prepares, _valueList
}

// 根据数据结构取得字段列表
func GetUpdateSql(_val interface{}, _primary bool) []string {

	_fields := make([]string, 0)

	_type := reflect.TypeOf(_val)
	_value := reflect.ValueOf(_val)

	for i := 0; i < _type.NumField(); i++ {

		if _type.Field(i).Name[0] >= 97 {
			continue // 小写的过滤掉
		}

		if !_primary {
			if strings.ToUpper(_type.Field(i).Tag.Get("primary")) == "TRUE" {
				continue
			}
		}

		var fileName = _type.Field(i).Tag.Get("db")
		var fieldValue = _value.Field(i).Interface()

		if fileName == "" {
			continue
		}

		if _type.Field(i).Tag.Get("db") == "" {
			continue
		}

		_fields = append(_fields, fmt.Sprintf("`%v` = '%v'", fileName, fieldValue))
	}
	return _fields
}
