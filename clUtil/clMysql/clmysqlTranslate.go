package clMysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clSuperMap"
	"strings"
	"time"
)

type ClTranslate struct {
	tx     *sql.Tx
	DBName string
}

// 查询事务
func (this *ClTranslate) Query(sqlstr string, args ...interface{}) (*DbResult, error) {

	if this.tx == nil {
		return nil, errors.New("错误: 事务指针为空")
	}

	if args != nil && len(args) > 0 {
		sqlstr = fmt.Sprintf(sqlstr, args...)
	}

	rows, err := queryTx(context.Background(), sqlstr, this.tx)
	if err != nil {
		return nil, err
	}
	var result DbResult
	result.ArrResult = make([]*clSuperMap.SuperMap, 0)
	result.Length = uint32(len(rows))

	for _, val := range rows {
		result.ArrResult = append(result.ArrResult, val)
	}

	return &result, nil
}

// 查询事务带超时时间
func (this *ClTranslate) QueryWithTimeout(timeout uint32, sqlstr string, args ...interface{}) (*DbResult, error) {

	if this.tx == nil {
		return nil, errors.New("错误: 事务指针为空")
	}

	if args != nil && len(args) > 0 {
		sqlstr = fmt.Sprintf(sqlstr, args...)
	}

	c := context.Background()
	if timeout > 0 {
		c, _ = context.WithTimeout(c, time.Duration(timeout)*time.Second)
	}
	rows, err := queryTx(c, sqlstr, this.tx)
	if err != nil {
		return nil, err
	}
	var result DbResult
	result.ArrResult = make([]*clSuperMap.SuperMap, 0)
	result.Length = uint32(len(rows))

	for _, val := range rows {
		result.ArrResult = append(result.ArrResult, val)
	}

	return &result, nil
}

// 执行事务

/*
*
lastSql := sqlstr

	if args != nil && len(args) != 0 {
		lastSql = fmt.Sprintf(sqlstr, args...)
	}
	if lastSql == "" {
		return 0, errors.New("SQL语句为空")
	}

	this.lastSql = lastSql
	if this.conn == nil {
		this.lastErr = "错误: SQL连线指针为空"
		return 0, errors.New("错误: SQL连线指针为nil pointer")
	}
	if timeout == 0 {
		stmt, err := this.conn.Prepare(sqlstr)
		if err != nil {
			return 0, err
		}
		res, err := stmt.Exec(args...)
		if err != nil {
			this.lastErr = fmt.Sprintf("执行失败! ERR:%v", err)
			this.lastSql = fmt.Sprintf(sqlstr, args...)
			return 0, err
		}

		if strings.HasPrefix(strings.ToLower(sqlstr), "insert") {
			return res.LastInsertId()
		}

		return res.RowsAffected()
	} else {
		stmt, err := this.conn.Prepare(sqlstr)
		if err != nil {
			return 0, err
		}
		c := context.Background()
		if timeout > 0 {
			c, _ = context.WithTimeout(c, time.Duration(timeout)*time.Second)
		}
		res, err := stmt.ExecContext(c, args...)
		if err != nil {
			this.lastErr = fmt.Sprintf("执行失败! ERR:%v", err)
			this.lastSql = fmt.Sprintf(sqlstr, args...)
			return 0, err
		}

		if strings.HasPrefix(strings.ToLower(sqlstr), "insert") {
			return res.LastInsertId()
		}

		return res.RowsAffected()
	}
*/
func (this *ClTranslate) ExecPrepare(timeout uint32, sqlstr string, args ...interface{}) (int64, error) {

	lastSql := sqlstr
	if args != nil && len(args) != 0 {
		lastSql = fmt.Sprintf(sqlstr, args...)
	}
	if lastSql == "" {
		return 0, errors.New("SQL语句为空")
	}

	if this.tx == nil {
		return 0, errors.New("错误: SQL连线指针为nil pointer")
	}
	if timeout == 0 {
		stmt, err := this.tx.Prepare(sqlstr)
		if err != nil {
			return 0, err
		}
		defer stmt.Close()
		res, err := stmt.Exec(args...)
		if err != nil {
			return 0, err
		}

		if strings.HasPrefix(strings.ToLower(sqlstr), "insert") {
			return res.LastInsertId()
		}

		return res.RowsAffected()
	} else {
		stmt, err := this.tx.Prepare(sqlstr)
		if err != nil {
			return 0, err
		}
		defer stmt.Close()
		c := context.Background()
		if timeout > 0 {
			c, _ = context.WithTimeout(c, time.Duration(timeout)*time.Second)
		}
		res, err := stmt.ExecContext(c, args...)
		if err != nil {
			return 0, err
		}

		if strings.HasPrefix(strings.ToLower(sqlstr), "insert") {
			return res.LastInsertId()
		}

		return res.RowsAffected()
	}
}

// 执行事务
func (this *ClTranslate) Exec(sqlstr string, args ...interface{}) (int64, error) {

	if this.tx == nil {
		return 0, errors.New("错误: 事务指针为 nil pointer")
	}

	if args != nil && len(args) != 0 {
		sqlstr = fmt.Sprintf(sqlstr, args...)
	}

	res, err := this.tx.Exec(sqlstr)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("%v, SQL:%v", err, sqlstr))
	}

	if strings.HasPrefix(strings.ToLower(sqlstr), "insert") {
		return res.LastInsertId()
	}

	return res.RowsAffected()
}

// 执行事务
func (this *ClTranslate) ExecWithTimeout(timeout uint32, sqlstr string, args ...interface{}) (int64, error) {

	if this.tx == nil {
		return 0, errors.New("错误: 事务指针为 nil pointer")
	}

	if args != nil && len(args) != 0 {
		sqlstr = fmt.Sprintf(sqlstr, args...)
	}

	c := context.Background()
	if timeout > 0 {
		c, _ = context.WithTimeout(c, time.Duration(timeout)*time.Second)
	}

	res, err := this.tx.ExecContext(c, sqlstr)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("%v, SQL:%v", err, sqlstr))
	}

	if strings.HasPrefix(strings.ToLower(sqlstr), "insert") {
		return res.LastInsertId()
	}

	return res.RowsAffected()
}

// 提交事务
func (this *ClTranslate) Commit() error {
	return this.tx.Commit()
}

// 回滚事务
func (this *ClTranslate) Rollback() error {
	return this.tx.Rollback()
}

// 使用DBPointer进行构建器创建
func (this *ClTranslate) NewBuilder() *SqlBuider {

	sqlbuild := SqlBuider{
		dbTx:   this,
		dbType: 2,
		dbname: this.DBName,
	}
	return &sqlbuild
}

/*
*

	是否存在这个表格
*/
func (this *ClTranslate) HasTable(tablename string) bool {

	var tables, _ = this.GetTables(tablename)
	for _, val := range tables {
		if strings.EqualFold(val, tablename) {
			return true
		}
	}
	return false
}

/**
 * 获取指定数据库下的所有表名字
 * @param dbname string 获取的数据库名
 * @param contain string 表名包含字符串，为空则取全部表
 * return
 * @1 数据表数组
 * @2 数据库
 */
func (this *ClTranslate) GetTables(contain string) ([]string, error) {

	querySql := ""
	if contain == "" {
		querySql = "SHOW TABLES"
	} else {
		querySql = "SHOW TABLES LIKE '%" + contain + "%'"
	}

	res, err := this.Query(querySql)
	if err != nil {
		return []string{}, err
	}

	if res.Length == 0 {
		return nil, nil
	}

	tables := make([]string, res.Length)
	for i := 0; i < int(res.Length); i++ {
		if contain != "" {
			tables[i] = res.ArrResult[i].GetStr("Tables_in_"+this.DBName+" (%"+contain+"%)", "")
		} else {
			tables[i] = res.ArrResult[i].GetStr("Tables_in_"+this.DBName, "")
		}
	}

	return tables, nil
}
