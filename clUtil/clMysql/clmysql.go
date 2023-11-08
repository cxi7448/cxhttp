package clMysql

/**
 *	高性能数据库封装类
 *
 *
 */
import (
	"context"
	"database/sql"
	"github.com/cxi7448/cxhttp/clUtil/clSuperMap"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//var clconnections ClmysqlConnections

func init() {

	// 初始化连接池
	//DBPointerPool = make(map[string]*DBPointer)
	//go clconnections.StartToCheck()
}

// 内部查询用
func query(sqlstr string, curdb *sql.DB, _timeout uint32) ([]*clSuperMap.SuperMap, error) {
	c := context.Background()
	if _timeout > 0 {
		c, _ = context.WithTimeout(c, time.Duration(_timeout)*time.Second)
	}
	row, err := curdb.QueryContext(c, sqlstr)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	columns, _ := row.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	// 将values的内存地址保存在scanArgs中
	for i := range values {
		scanArgs[i] = &values[i]
	}
	result := make([]*clSuperMap.SuperMap, 0)
	for row.Next() {
		records := clSuperMap.NewSuperMap()
		row.Scan(scanArgs...) // 获取扫描后的数组
		for i, col := range values {
			if col == nil {
				continue
			}
			records.Add(columns[i], string(col.([]byte)))
		}
		result = append(result, records)
	}
	return result, nil
}

// 内部查询用
func queryTx(_context context.Context, sqlstr string, tx *sql.Tx) ([]*clSuperMap.SuperMap, error) {

	row, err := tx.QueryContext(_context, sqlstr)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	columns, _ := row.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	// 将values的内存地址保存在scanArgs中
	for i := range values {
		scanArgs[i] = &values[i]
	}
	result := make([]*clSuperMap.SuperMap, 0)
	for row.Next() {
		records := clSuperMap.NewSuperMap()
		row.Scan(scanArgs...) // 获取扫描后的数组
		for i, col := range values {
			if col == nil {
				continue
			}
			records.Add(columns[i], string(col.([]byte)))
		}
		result = append(result, records)
	}
	return result, nil
}
