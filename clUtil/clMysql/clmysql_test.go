package clMysql

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"testing"
)

func TestAddMaster(t *testing.T) {

	db := NewDBSimple("127.0.0.1", "root", "root", "miner_new")
	if db == nil {
		fmt.Printf("connect to mysql failed\n")
		return
	}

	tableList, err := db.GetTables("")
	if err != nil {
		clLog.Error("获取数据库表格错误: %v", err)
		return
	}
	for _, table := range tableList {
		clLog.Debug("表格: %v", table)
		columnList, err := db.GetTableColumns(table)
		if err != nil {
			clLog.Error("获取表格字段数据错误: %v", err)
			continue
		}
		clLog.Debug("表格: %v 字段数据: %+v", table, columnList)
	}

	clLog.Debug("生成完毕!")
}

type AddObjData struct {
	A uint32  `db:"a" primary:"TRUE"`
	B uint32  `db:"b"`
	C float64 `db:"c"`
}

func TestSqlBuider_AddMulti(t *testing.T) {

	db := NewDBSimple("127.0.0.1", "root", "root", "miner_new")
	if db == nil {
		fmt.Printf("connect to mysql failed\n")
		return
	}
	var data AddObjData
	data.A = 100
	data.B = 3
	data.C = 2292929.988778

	db.NewBuilder().Table("test").AddObj(data, true)

}

type AddObjMultiObj struct {
	A uint32 `db:"a" primary:"TRUE"`
	B uint32 `db:"b"`
	C uint32 `db:"c"`
}

func TestSqlBuider_AddObjMulti(t *testing.T) {
	db := NewDBSimple("127.0.0.1", "root", "root", "miner_new")
	if db == nil {
		fmt.Printf("connect to mysql failed\n")
		return
	}
	db.NewBuilder().Table("test").AddObjMulti([]interface{}{
		AddObjMultiObj{
			A: 1,
			B: 2,
			C: 3,
		},
		AddObjMultiObj{
			A: 1,
			B: 2,
			C: 3,
		},
		AddObjMultiObj{
			A: 1,
			B: 2,
			C: 3,
		},
	}, true)

}
