package clMysql

import (
	"errors"
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clCommon"
	"github.com/cxi7448/cxhttp/clUtil/clSuperMap"
	"github.com/cxi7448/cxhttp/clUtil/clTime"
	"reflect"
	"strings"
)

type sqlJoin struct {
	TableName     string
	JoinCondition string
}

/*
*

	数据库语句生成器
*/
type SqlBuider struct {
	tablename  string
	dbname     string
	whereStr   string
	fieldStr   string
	updateData map[string]string
	orders     string
	random     *int // 随机因子
	limit      string
	group      string
	having     string
	expire     uint32
	timeout    uint32 // 超时时间，秒
	finalSql   string

	addColumns    []MySqlColumns // 保存的字段
	removeColumns []string       // 待删除的字段
	addIndexs     []string       // 待添加的索引
	removeIndexs  []string       // 待删除的索引

	lastColumns []MySqlColumns
	lastIndexs  []string
	primaryKeys string

	lastTable   string
	lastTableId string

	unionalls []string
	unions    []string

	duplicateKey []string

	leftJoin  []sqlJoin
	rightJoin []sqlJoin

	dbType    uint32
	dbPointer *DBPointer
	dbTx      *ClTranslate
}

type MySqlColumns struct {
	name     string // 字段名称
	typename string // 字段类型
	null     bool   // 是否为空
	defaults string // 默认值
	autoInc  bool   // 是否自动递增
	comment  string // 备注
}

func NewBuilder() *SqlBuider {

	sqlbuild := SqlBuider{}
	return &sqlbuild
}

// 使用DBPointer进行构建器创建
func (this *DBPointer) NewBuilder() *SqlBuider {

	sqlbuild := SqlBuider{
		dbPointer: this,
		dbType:    1,
		dbname:    this.Dbname,
		leftJoin:  make([]sqlJoin, 0),
		rightJoin: make([]sqlJoin, 0),
	}
	return &sqlbuild
}

// 使用DBPointer进行构建器创建
func (this *DBPointer) NewBuilderByTable(_tableName string) *SqlBuider {

	sqlbuild := SqlBuider{
		dbPointer: this,
		dbType:    1,
		dbname:    this.Dbname,
		leftJoin:  make([]sqlJoin, 0),
		rightJoin: make([]sqlJoin, 0),
	}
	return sqlbuild.Table(_tableName)
}

/**
 * 设置表格名称
 * @param tablename string  要设置的表格名称
 */
func (this *SqlBuider) Table(tablename string) *SqlBuider {
	this.tablename = tablename
	this.whereStr = ""
	this.fieldStr = ""
	this.updateData = make(map[string]string)
	this.orders = ""
	this.limit = ""
	this.group = ""
	this.having = ""
	this.expire = 0
	this.finalSql = ""

	this.addColumns = make([]MySqlColumns, 0) // 保存的字段
	this.removeColumns = make([]string, 0)    // 待删除的字段
	this.addIndexs = make([]string, 0)        // 待添加的索引
	this.removeIndexs = make([]string, 0)     // 待删除的索引

	this.lastColumns = make([]MySqlColumns, 0)
	this.lastIndexs = make([]string, 0)
	this.primaryKeys = ""

	this.lastTable = ""
	this.lastTableId = ""

	this.unions = make([]string, 0)
	this.unionalls = make([]string, 0)

	this.duplicateKey = make([]string, 0)

	return this
}

func (this *SqlBuider) List(table, where string, pageid, pcount int32, rows interface{}, _order ...string) (int32, error) {
	builder := this.Table(table)
	if len(_order) > 0 {
		builder.Order(_order[0])
	}
	err := builder.Where(where).Page(pageid, pcount).FindAll(rows)
	if err != nil {
		return 0, err
	}
	total := int32(reflect.ValueOf(rows).Elem().Len())
	if pageid == 0 && total == pcount {
		i, err := this.Table(table).Where(where).Count()
		if err != nil {
			fmt.Printf("Mysql.List错误:%v\n", err)
		}
		total = i
	}
	return total, nil
}

/*
*

	设置WHERE条件
	@param wherestr string WHERE条件文本
*/
func (this *SqlBuider) Where(wherestr string, args ...interface{}) *SqlBuider {

	if args == nil || len(args) == 0 {
		this.whereStr = wherestr
	} else {
		this.whereStr = fmt.Sprintf(wherestr, args...)
	}

	return this
}

/**
 * 设置重复要更新的key列表
 * @param keys 要更新的字段名字列表
 */
func (this *SqlBuider) OnDuplicateKey(keys []string) *SqlBuider {
	this.duplicateKey = keys
	return this
}

/**
 * 设置要查询的Field
 * @param fiedlStr string FIELD字段列表
 */
func (this *SqlBuider) Field(fieldStr string) *SqlBuider {
	this.fieldStr = fieldStr
	return this
}

/**
 * 设置要查询的Field
 * @param fiedlStr string FIELD字段列表
 */
func (this *SqlBuider) FieldList(fields []string) *SqlBuider {
	this.fieldStr = "`" + strings.Join(fields, "`,`") + "`"
	return this
}

/*
*

	设置分组
	@param group string 分组内容
*/
func (this *SqlBuider) Group(groupStr string) *SqlBuider {
	this.group = groupStr
	return this
}

/*
*

	设置排序方式
	@param orders string 排序内容
*/
func (this *SqlBuider) Order(orders string) *SqlBuider {
	this.orders = orders
	return this
}

/*
*

	设置排序方式
	@param orders string 排序内容
*/
func (this *SqlBuider) OrderRand(random int) *SqlBuider {
	this.random = &random
	return this
}

/*
*

	设置排序方式
	@param orders string 排序内容
*/
func (this *SqlBuider) OrderRandHour() *SqlBuider {
	ct, _ := clTime.NewDate("")
	var hour = int(ct.Hour)
	this.random = &hour
	return this
}

/*
*

	设置Cache 缓存时间
	@param expire int32 缓存有效期
*/
func (this *SqlBuider) Timeout(_timeout uint32) *SqlBuider {
	this.timeout = _timeout
	return this
}

/*
*

	设置LIMIT限制
	@param min int32 设置limit的最小值
	@param count int32 设置limit的数量
*/
func (this *SqlBuider) Limit(min int32, count int32) *SqlBuider {
	this.limit = fmt.Sprintf(" LIMIT %v, %v", min, count)
	return this
}

/*
*

	设置LIMIT限制
	@param min int32 设置limit的最小值
	@param count int32 设置limit的数量
*/
func (this *SqlBuider) Page(page int32, count int32) *SqlBuider {
	this.limit = fmt.Sprintf(" LIMIT %v, %v", page*count, count)
	return this
}

/*
*
设置DB名字
@param dbname string 设置DB名字
*/
func (this *SqlBuider) DB(dbname string) *SqlBuider {
	this.dbname = dbname
	return this
}

/**
 * 查询语句并返回结果集
 */
func (this *SqlBuider) Query() (*DbResult, error) {

	sqlStr, buildErr := this.buildQuerySql()
	if buildErr != nil {
		return nil, buildErr
	}
	this.finalSql = sqlStr

	var resp, err = this.QueryCustom(sqlStr)
	return resp, err
}

/**
 * 查询语句
 * 获取指定索引处的数据
 * @param idx int 索引id
 */
func (this *SqlBuider) Find(idx uint32) (*clSuperMap.SuperMap, error) {

	sqlStr, buildErr := this.buildQuerySql()
	if buildErr != nil {
		return nil, buildErr
	}
	this.finalSql = sqlStr

	resp, err := this.QueryCustom(sqlStr)

	if err != nil {
		return nil, err
	}

	if resp == nil || resp.Length <= idx {
		return nil, nil
	}

	return resp.ArrResult[idx], nil
}

/**
 * 查询语句
 * 获取指定索引处的数据
 * @param idx int 索引id
 */
func (this *SqlBuider) Count() (int32, error) {

	this.fieldStr = fmt.Sprintf("COUNT(*) as t_count")
	sqlStr, buildErr := this.buildQuerySql()
	if buildErr != nil {
		return 0, buildErr
	}

	this.finalSql = sqlStr
	resp, err := this.QueryCustom(sqlStr)

	if err != nil {
		return 0, err
	}
	if resp == nil || resp.Length == 0 {
		return 0, err
	}

	return resp.ArrResult[0].GetInt32("t_count", 0), nil
}

/**
 * 查询语句
 * 获取指定索引处的数据
 * @param idx int 索引id
 */
func (this *SqlBuider) Max(_field string) (uint64, error) {

	this.fieldStr = fmt.Sprintf("MAX(`%v`) as max_id", _field)
	sqlStr, buildErr := this.buildQuerySql()
	if buildErr != nil {
		return 0, buildErr
	}

	this.finalSql = sqlStr
	resp, err := this.QueryCustom(sqlStr)
	if err != nil {
		return 0, err
	}
	if resp == nil || resp.Length == 0 {
		return 0, err
	}

	return resp.ArrResult[0].GetUInt64("max_id", 0), nil
}

/*
*

	事务查询语句
*/
func (this *SqlBuider) SelectTx() (*DbResult, error) {

	sqlStr, buildErr := this.buildQuerySql()
	if buildErr != nil {
		return nil, buildErr
	}

	this.finalSql = sqlStr

	resp, err := this.QueryCustom(sqlStr)
	return resp, err
}

/*
*

	更新语句.
	@param data map[string] string 需要更新字段列表
	@return 修改成功个数, 错误
*/
func (this *SqlBuider) Save(data map[string]interface{}) (int64, error) {

	if this.tablename == "" {
		return 0, errors.New("EMPTY TABLE NAME")
	}

	fieldstr := ""
	valueStr := []interface{}{}
	for key, val := range data {
		if fieldstr != "" {
			fieldstr += ","
		}
		fieldstr += fmt.Sprintf("`%v` = ?", key)
		valueStr = append(valueStr, val)
	}

	if fieldstr == "" {
		return 0, errors.New("EMPTY UPDATE COLUMN LIST")
	}

	sqlStr := fmt.Sprintf("UPDATE %v SET %v WHERE %v", this.tablename, fieldstr, this.whereStr)
	this.finalSql = sqlStr

	//var resp, err = this.ExecCustom(this.finalSql)
	var resp, err = this.ExecPrepare(this.finalSql, valueStr...)
	return resp, err
}

// 基于对象的修改
func (this *SqlBuider) SaveObj(_resp interface{}) (int64, error) {

	fieldList := GetUpdateSql(_resp, false)

	this.finalSql = fmt.Sprintf("UPDATE `%v`.`%v` SET %v WHERE %v", this.dbname, this.tablename, strings.Join(fieldList, ","), this.whereStr)

	resp, err := this.ExecCustom(this.finalSql)
	return resp, err
}

/*
*

	删除语句.
	@return 删除个数, 错误
*/
func (this *SqlBuider) Del() (int64, error) {

	if this.tablename == "" {
		return 0, errors.New("EMPTY TABLE NAME")
	}

	sqlStr := fmt.Sprintf("DELETE FROM %v WHERE %v", this.tablename, this.whereStr)
	this.finalSql = sqlStr

	resp, err := this.ExecCustom(sqlStr)
	return resp, err
}

/*
*

	事务更新语句.
	@param data map[string] string 需要更新字段列表
	@return 修改成功个数, 错误
*/
func (this *SqlBuider) SaveTx(data map[string]interface{}) (int64, error) {

	if this.tablename == "" {
		return 0, errors.New("EMPTY TABLE NAME")
	}

	fieldstr := ""
	for key, val := range data {
		if fieldstr != "" {
			fieldstr += ","
		}
		fieldstr += fmt.Sprintf("%v = '%v'", key, val)
	}

	if fieldstr == "" {
		return 0, errors.New("EMPTY UPDATE COLUMN LIST")
	}

	sqlStr := fmt.Sprintf("UPDATE %v SET %v WHERE %v", this.tablename, fieldstr, this.whereStr)
	this.finalSql = sqlStr

	resp, err := this.ExecCustom(this.finalSql)
	return resp, err
}

/*
*

	初始化整个表.
	@return 返回是否发生错误
*/
func (this *SqlBuider) Truncate() error {

	if this.tablename == "" {
		return errors.New("EMPTY TABLE NAME")
	}

	this.finalSql = fmt.Sprintf("TRUNCATE TABLE %v", this.tablename)

	_, err := this.ExecCustom(this.finalSql)
	return err
}

/*
*

	添加语句
	@param data map[string] string 需要添加的字段列表
	@return 最后一条添加的id
*/
func (this *SqlBuider) Add(data map[string]interface{}) (int64, error) {

	if this.tablename == "" {
		return 0, errors.New("EMPTY TABLE NAME")
	}

	// 拼接字段区和值字段区
	fieldstr := strings.Builder{}
	valuestr := []interface{}{}
	prepareValue := strings.Builder{}
	for key, val := range data {
		if prepareValue.Len() > 0 {
			prepareValue.WriteString(",")
			fieldstr.WriteString(",")
		}

		fieldstr.WriteString(fmt.Sprintf("`%v`", key))
		prepareValue.WriteString("?")
		//if key == "guid" || key == "uid" || key == "id" {
		//	valuestr.WriteString(fmt.Sprintf("%v", val))
		//} else {
		//	valuestr.WriteString(fmt.Sprintf("'%v'", val))
		//}
		valuestr = append(valuestr, val)
	}

	if fieldstr.Len() == 0 {
		return 0, errors.New("EMPTY UPDATE COLUMN LIST")
	}

	// 拼接重复区
	onDuplicateStr := strings.Builder{}
	if this.duplicateKey != nil && len(this.duplicateKey) > 0 {
		onDuplicateStr.WriteString(" ON DUPLICATE KEY UPDATE ")

		for i, val := range this.duplicateKey {
			if i > 0 {
				onDuplicateStr.WriteString(",")
			}
			onDuplicateStr.WriteString(fmt.Sprintf("`%[1]v` = VALUES(`%[1]v`)", val))
		}
	}

	sqlStr := fmt.Sprintf("INSERT INTO %v (%v) VALUES(%v) %v", this.tablename, fieldstr.String(), prepareValue.String(), onDuplicateStr.String())
	this.finalSql = sqlStr
	resp, err := this.ExecPrepare(this.finalSql, valuestr...)
	return resp, err
}

/*
*

	批量添加语句
	@param data map[string] string 需要添加的字段列表
	@return 最后一条添加的id
*/
func (this *SqlBuider) AddMulti(_list []map[string]interface{}) (int64, error) {

	if this.tablename == "" {
		return 0, errors.New("EMPTY TABLE NAME")
	}
	if len(_list) == 0 {
		return 0, errors.New("EMPTY DATA")
	}

	// 拼接字段区和值字段区
	fieldstr := make([]string, 0)
	for field, _ := range _list[0] {
		fieldstr = append(fieldstr, fmt.Sprintf("%v", field))
	}

	valueList := make([]string, 0)

	for _, data := range _list {
		var valueArr = make([]string, len(fieldstr))
		for field, val := range data {

			idx := clCommon.IndexOfStringArray(fieldstr, field, true)
			if idx == -1 {
				continue
			}

			valueArr[idx] = fmt.Sprintf("'%v'", val)
		}
		valueList = append(valueList, "("+strings.Join(valueArr, ",")+")")
	}

	if len(valueList) == 0 || len(fieldstr) == 0 {
		return 0, errors.New("EMPTY UPDATE COLUMN LIST")
	}

	// 拼接重复区
	onDuplicateStr := strings.Builder{}
	if this.duplicateKey != nil && len(this.duplicateKey) > 0 {
		onDuplicateStr.WriteString(" ON DUPLICATE KEY UPDATE ")

		for i, val := range this.duplicateKey {
			if i > 0 {
				onDuplicateStr.WriteString(",")
			}
			onDuplicateStr.WriteString(fmt.Sprintf("`%[1]v` = VALUES(`%[1]v`)", val))
		}
	}

	this.finalSql = fmt.Sprintf("INSERT INTO %v (`%v`) VALUES %v %v", this.tablename, strings.Join(fieldstr, "`,`"), strings.Join(valueList, ","), onDuplicateStr.String())

	resp, err := this.ExecCustom(this.finalSql)
	return resp, err
}

/*
*

	添加语句
	@param data map[string] string 需要添加的字段列表
	@return 最后一条添加的id
*/
func (this *SqlBuider) Replace(data map[string]interface{}) (int64, error) {

	if this.tablename == "" {
		return 0, errors.New("EMPTY TABLE NAME")
	}

	fieldstr := strings.Builder{}
	valuestr := strings.Builder{}
	for key, val := range data {
		if valuestr.Len() > 0 {
			valuestr.WriteString(",")
			fieldstr.WriteString(",")
		}
		fieldstr.WriteString(fmt.Sprintf("`%v`", key))
		if key == "guid" || key == "uid" || key == "id" {
			valuestr.WriteString(fmt.Sprintf("%v", val))
		} else {
			valuestr.WriteString(fmt.Sprintf("'%v'", val))
		}
	}

	if fieldstr.Len() == 0 {
		return 0, errors.New("EMPTY UPDATE COLUMN LIST")
	}

	sqlStr := fmt.Sprintf("REPLACE INTO %v (%v) VALUES(%v)", this.tablename, fieldstr.String(), valuestr.String())
	this.finalSql = sqlStr
	resp, err := this.ExecCustom(sqlStr)
	return resp, err
}

/*
*

	添加语句
	@param data map[string] string 需要添加的字段列表
	@return 最后一条添加的id
*/
func (this *SqlBuider) ReplaceNew(data map[string]interface{}) (int64, error) {

	if this.tablename == "" {
		return 0, errors.New("EMPTY TABLE NAME")
	}

	fieldstr := strings.Builder{}
	valuestr := strings.Builder{}
	valueArgs := []interface{}{}
	for key, val := range data {
		if fieldstr.Len() > 0 {
			valuestr.WriteString(",")
			fieldstr.WriteString(",")
		}
		fieldstr.WriteString(fmt.Sprintf("`%v`", key))
		valuestr.WriteString("?")
		valueArgs = append(valueArgs, val)
	}

	if fieldstr.Len() == 0 {
		return 0, errors.New("EMPTY UPDATE COLUMN LIST")
	}
	sqlStr := fmt.Sprintf("REPLACE INTO %v (%v) VALUES(%v)", this.tablename, fieldstr.String(), valuestr.String())
	this.finalSql = sqlStr
	resp, err := this.ExecPrepare(sqlStr, valueArgs...)
	return resp, err
}

/*
*

	事务添加语句
	@param data map[string] string 需要添加的字段列表
	@return 最后一条添加的id
*/
func (this *SqlBuider) AddTx(data map[string]interface{}) (int64, error) {

	if this.tablename == "" {
		return 0, errors.New("EMPTY TABLE NAME")
	}

	fieldstr := strings.Builder{}
	valuestr := strings.Builder{}
	for key, val := range data {
		if valuestr.Len() > 0 {
			valuestr.WriteString(",")
			fieldstr.WriteString(",")
		}
		fieldstr.WriteString(fmt.Sprintf("`%v`", key))
		if key == "guid" || key == "uid" || key == "id" {
			valuestr.WriteString(fmt.Sprintf("%v", val))
		} else {
			valuestr.WriteString(fmt.Sprintf("'%v'", val))

		}
	}

	if fieldstr.Len() == 0 {
		return 0, errors.New("EMPTY UPDATE COLUMN LIST")
	}

	// 拼接重复区
	onDuplicateStr := strings.Builder{}
	if this.duplicateKey != nil && len(this.duplicateKey) > 0 {
		onDuplicateStr.WriteString(" ON DUPLICATE KEY UPDATE ")

		for i, val := range this.duplicateKey {
			if i > 0 {
				onDuplicateStr.WriteString(",")
			}
			onDuplicateStr.WriteString(fmt.Sprintf("`%[1]v` = VALUES(`%[1]v`)", val))
		}
	}

	sqlStr := fmt.Sprintf("INSERT INTO %v (%v) VALUES(%v) %v", this.tablename, fieldstr.String(), valuestr.String(), onDuplicateStr.String())
	this.finalSql = sqlStr
	resp, err := this.ExecCustom(this.finalSql)
	return resp, err
}

/*
*

	添加字段
	@param col string 字段名称
	@param tpname string 字段类型
	@param isnull bool 是否为空
	@param comment string 备注
*/
func (this *SqlBuider) AddColumn(col string, tpname string, isnull bool, defval string, comment string) *SqlBuider {

	this.addColumns = append(this.addColumns, MySqlColumns{
		name:     col,
		typename: tpname,
		null:     isnull,
		defaults: defval,
		comment:  comment,
		autoInc:  false,
	})
	return this
}

/*
*

	删除字段
	@param col 字段名称
*/
func (this *SqlBuider) RemoveColumn(col string) *SqlBuider {
	this.removeColumns = append(this.removeColumns, col)
	return this
}

/*
*

	设置为主键
	col : 需要设置为主键的字段名称
	auto : 是否自动递增
*/
func (this *SqlBuider) SetId(col string, autoInc bool) *SqlBuider {
	for key, val := range this.addColumns {
		if val.name == col {
			if autoInc {
				this.addColumns[key].autoInc = true
			}
			break
		} else {
			if autoInc {
				this.addColumns[key].autoInc = false
			}
		}
	}

	this.primaryKeys = col
	return this
}

/*
*

	添加索引
	cols: 要添加的索引，用逗号隔开
	unique: 是否不重复
*/
func (this *SqlBuider) AddIndex(cols string, unique bool) *SqlBuider {
	indexStr := ""
	if unique {
		indexStr = "UNIQUE KEY"
	} else {
		indexStr = "KEY"
	}
	indexStr += " (" + cols + ")"
	this.addIndexs = append(this.addIndexs, indexStr)
	return this
}

/*
*

	移除索引
*/
func (this *SqlBuider) RemoveIndex(cols string) *SqlBuider {
	this.removeIndexs = append(this.removeIndexs, cols)
	return this
}

/**
 * 创建表格
 * @param overWrite 是否覆盖, 如果为true则会先删除原先的表
 */
func (this *SqlBuider) CreateTable(overWrite bool) bool {
	//生成整个表格的SQL
	if this.tablename == "" {
		return false
	}

	if len(this.addColumns) == 0 {
		return false
	}

	sqlStr := strings.Builder{}
	sqlStr.WriteString("CREATE TABLE IF NOT EXISTS " + this.tablename + "(")
	for key, val := range this.addColumns {

		sqlStr.WriteString("`" + val.name + "` " + val.typename)
		if val.autoInc == true {
			sqlStr.WriteString(" AUTO_INCREMENT")
		} else {
			if val.null == false && val.typename != "text" {
				sqlStr.WriteString(" NOT NULL DEFAULT '" + val.defaults + "'")
			}
		}
		sqlStr.WriteString(" COMMENT '" + val.comment + "'")

		if key < len(this.addColumns)-1 {
			sqlStr.WriteString(",")
		}
	}
	if this.primaryKeys != "" {
		sqlStr.WriteString(", PRIMARY KEY (" + this.primaryKeys + ")")
	}

	if len(this.addIndexs) > 0 {
		sqlStr.WriteString("," + strings.Join(this.addIndexs, ","))
	}
	sqlStr.WriteString(") ENGINE=INNODB DEFAULT CHARSET=UTF8")

	var err error
	if overWrite {
		switch this.dbType {
		case 1: // Picker
			this.dbPointer.Exec("DROP TABLE IF EXISTS " + this.tablename)
		case 2: // 事务
			this.dbTx.Exec("DROP TABLE IF EXISTS " + this.tablename)
		}
	}

	this.finalSql = sqlStr.String()
	_, err = this.ExecCustom(this.finalSql)
	return err == nil
}

// 修改表格结构
func (this *SqlBuider) SaveTable() error {

	//生成整个表格的SQL
	if this.tablename == "" {
		return errors.New("TABLE NAME IS EMPTY")
	}

	SqlStr := strings.Builder{}
	SqlStr.WriteString("ALTER TABLE " + this.tablename)
	AnyThing := false

	if len(this.addColumns) > 0 {
		// 添加字段
		AnyThing = true
		for key, val := range this.addColumns {
			SqlStr.WriteString("ADD COLUMN " + val.name + " " + val.typename)
			if !val.null {
				SqlStr.WriteString(" NOT NULL")
			}

			SqlStr.WriteString(" DEFAULT '" + val.defaults + "' COMMENT '" + val.comment + "'")
			if key < len(this.addColumns)-1 {
				SqlStr.WriteString(",")
			}
		}
	}

	if len(this.removeColumns) > 0 {
		if AnyThing {
			SqlStr.WriteString(",")
		}
		AnyThing = true
		for key, val := range this.removeColumns {
			SqlStr.WriteString("DROP COLUMN " + val)
			if key < len(this.removeColumns)-1 {
				SqlStr.WriteString(",")
			}
		}
	}

	if len(this.addIndexs) > 0 {
		if AnyThing {
			SqlStr.WriteString(",")
		}
		AnyThing = true
		for key, val := range this.addIndexs {
			SqlStr.WriteString("ADD INDEX " + val)
			if key < len(this.addIndexs)-1 {
				SqlStr.WriteString(",")
			}
		}
	}

	if len(this.removeIndexs) > 0 {
		if AnyThing {
			SqlStr.WriteString(",")
		}
		for key, val := range this.removeIndexs {
			SqlStr.WriteString("DROP INDEX " + val)
			if key < len(this.removeIndexs)-1 {
				SqlStr.WriteString(",")
			}
		}
	}

	this.finalSql = SqlStr.String()
	_, err := this.ExecCustom(this.finalSql)
	return err
}

// 联表查询
// @param tablename string 联表名称
func (this *SqlBuider) UnionAll(tablename string) *SqlBuider {
	for _, val := range this.unionalls {
		if val == tablename {
			return this
		}
	}
	this.unionalls = append(this.unionalls, tablename)
	return this
}

// 强制联表查询
// @param tablename string 联表名称
func (this *SqlBuider) Union(tablename string) *SqlBuider {
	for _, val := range this.unions {
		if val == tablename {
			return this
		}
	}
	this.unions = append(this.unionalls, tablename)
	return this
}

/*
*
左内联
@param _tableName string 表名
@param _joinCondition string 条件
*/
func (this *SqlBuider) LeftJoin(_tableName string, _joinCondition string) *SqlBuider {
	this.leftJoin = append(this.leftJoin, sqlJoin{
		TableName:     _tableName,
		JoinCondition: _joinCondition,
	})
	return this
}

/*
*
右内联
@param _tableName string 表名
@param _joinCondition string 条件
*/
func (this *SqlBuider) RightJoin(_tableName string, _joinCondition string) *SqlBuider {
	this.rightJoin = append(this.rightJoin, sqlJoin{
		TableName:     _tableName,
		JoinCondition: _joinCondition,
	})
	return this
}

/*
获取sql语句
*/
func (this *SqlBuider) GetLastSql() string {
	return this.finalSql
}

// 获取查找
func (this *SqlBuider) FindAll(_resp interface{}) error {

	_value := reflect.ValueOf(_resp)
	_valueE := _value.Elem()
	_valueE = _valueE.Slice(0, _valueE.Cap())

	_element := _valueE.Type().Elem()
	if this.fieldStr == "" {
		fieldList := GetAllField(reflect.New(_element).Interface())
		this.fieldStr = "`" + strings.Join(fieldList, "`,`") + "`"
	}

	sqlStr, buildErr := this.buildQuerySql()
	if buildErr != nil {
		return buildErr
	}

	this.finalSql = sqlStr

	resp, err := this.QueryCustom(this.finalSql)

	if err != nil {
		return err
	}

	if resp != nil && resp.Length > 0 {
		i := 0
		for idx, row := range resp.ArrResult {
			// 需要添加
			if _valueE.Len() == idx {
				elemp := reflect.New(_element)
				Unmarsha(row, elemp.Interface())
				_valueE = reflect.Append(_valueE, elemp.Elem())
			}
			i++
		}

		_value.Elem().Set(_valueE.Slice(0, i))
	}

	return nil
}

// 获取查找
func (this *SqlBuider) FindOne(_resp interface{}) error {

	if this.fieldStr == "" {
		fieldList := GetAllField(_resp)
		this.fieldStr = "`" + strings.Join(fieldList, "`,`") + "`"
	}
	this.limit = "LIMIT 1"

	sqlStr, buildErr := this.buildQuerySql()
	if buildErr != nil {
		return buildErr
	}

	this.finalSql = sqlStr

	resp, err := this.QueryCustom(this.finalSql)
	if err != nil {
		return err
	}

	if resp == nil || resp.Length == 0 {
		return errors.New("not found")
	}
	Unmarsha(resp.ArrResult[0], _resp)
	return nil
}

// 获取查找
func (this *SqlBuider) AddObj(_resp interface{}, _include_primary bool) (int64, error) {

	fieldList, valuesList := GetInsertSql(_resp, _include_primary)

	// 拼接重复区
	onDuplicateStr := strings.Builder{}
	if this.duplicateKey != nil && len(this.duplicateKey) > 0 {
		onDuplicateStr.WriteString(" ON DUPLICATE KEY UPDATE ")

		for i, val := range this.duplicateKey {
			if i > 0 {
				onDuplicateStr.WriteString(",")
			}
			onDuplicateStr.WriteString(fmt.Sprintf("`%[1]v` = VALUES(`%[1]v`)", val))
		}
	}

	sqlStr := fmt.Sprintf("INSERT INTO `%v`.`%v` (`%v`) VALUES('%v') %v", this.dbname, this.tablename, strings.Join(fieldList, "`,`"), strings.Join(valuesList, "','"), onDuplicateStr.String())
	this.finalSql = sqlStr

	resp, err := this.ExecCustom(sqlStr)
	return resp, err
}

// 获取查找
func (this *SqlBuider) AddObjMulti(_resp []interface{}, _includePrimary bool) (int64, error) {

	fieldList, prepareList, valuesList := GetInsertSqlMultiNew(_resp, _includePrimary)
	// 拼接重复区
	onDuplicateStr := strings.Builder{}
	if this.duplicateKey != nil && len(this.duplicateKey) > 0 {
		onDuplicateStr.WriteString(" ON DUPLICATE KEY UPDATE ")

		for i, val := range this.duplicateKey {
			if i > 0 {
				onDuplicateStr.WriteString(",")
			}
			onDuplicateStr.WriteString(fmt.Sprintf("`%[1]v` = VALUES(`%[1]v`)", val))
		}
	}
	this.finalSql = fmt.Sprintf("INSERT INTO `%v`.`%v` (`%v`) VALUES %v %v",
		this.dbname,
		this.tablename,
		strings.Join(fieldList, "`,`"),
		strings.Join(prepareList, ","),
		onDuplicateStr.String())
	//resp, err := this.ExecCustom(this.finalSql)
	resp, err := this.ExecPrepare(this.finalSql, valuesList...)
	return resp, err
}

// 覆盖
func (this *SqlBuider) ReplaceObj(_resp interface{}, _include_primary bool) (int64, error) {

	fieldList, valuesList := GetInsertSql(_resp, _include_primary)

	sqlStr := fmt.Sprintf("REPLACE INTO `%v`.`%v` (`%v`) VALUES('%v')", this.dbname, this.tablename, strings.Join(fieldList, "`,`"), strings.Join(valuesList, "','"))
	this.finalSql = sqlStr

	resp, err := this.ExecCustom(sqlStr)
	return resp, err
}

func (this *SqlBuider) buildQuerySql() (string, error) {
	if this.tablename == "" {
		return "", errors.New("EMPTY TABLE NAME")
	}

	if this.whereStr == "" {
		this.whereStr = "1"
	}

	if this.fieldStr == "" {
		this.fieldStr = "*"
	}

	var extraSql = ""
	if this.group != "" {
		extraSql += "GROUP BY " + this.group
	}

	joinStr := ""
	var FinallySql = ""
	if len(this.unionalls) > 0 {
		// 使用一般联表查询
		for _, sub := range this.unionalls {
			if FinallySql == "" {
				FinallySql = fmt.Sprintf("SELECT %v FROM %v WHERE ( %v ) %v UNION ALL ", this.fieldStr, this.tablename, this.whereStr, extraSql)
			} else {
				FinallySql += " UNION ALL "
			}
			FinallySql += fmt.Sprintf("SELECT %v FROM %v WHERE ( %v ) %v ", this.fieldStr, sub, this.whereStr, extraSql)
		}
	} else if len(this.unions) > 0 {
		// 使用强制联表查询
		for _, sub := range this.unions {
			if FinallySql == "" {
				FinallySql = fmt.Sprintf("SELECT %v FROM %v WHERE ( %v ) %v UNION ", this.fieldStr, this.tablename, this.whereStr, extraSql)
			} else {
				FinallySql += " UNION "
			}
			FinallySql += fmt.Sprintf("SELECT %v FROM %v WHERE ( %v ) %v ", this.fieldStr, sub, this.whereStr, extraSql)
		}
	} else {

		for _, val := range this.leftJoin {
			joinStr += fmt.Sprintf(" LEFT JOIN %v ON (%v)", val.TableName, val.JoinCondition)
		}
		for _, val := range this.rightJoin {
			joinStr += fmt.Sprintf(" RIGHT JOIN %v ON (%v)", val.TableName, val.JoinCondition)
		}
	}

	if this.orders != "" || this.random != nil {
		extraSql += " ORDER BY "
		if this.orders != "" {
			extraSql += this.orders
		}
		if this.random != nil {
			if this.orders != "" {
				extraSql += ","
			}
			extraSql += fmt.Sprintf(" rand(%v) ", *this.random)
		}
	}

	if this.limit != "" {
		extraSql += " " + this.limit
	}

	if FinallySql == "" {
		FinallySql = fmt.Sprintf("SELECT %v FROM %v %v WHERE ( %v ) %v", this.fieldStr, this.tablename, joinStr, this.whereStr, extraSql)
	} else {
		FinallySql = fmt.Sprintf("SELECT %v FROM ( %v ) temp %v", this.fieldStr, FinallySql, extraSql)
	}
	return FinallySql, nil
}

func (this *SqlBuider) ExecPrepare(sqlStr string, args ...interface{}) (int64, error) {
	var resp int64
	var err error
	switch this.dbType {
	case 1: // Picker
		resp, err = this.dbPointer.ExecPrepare(this.timeout, sqlStr, args...)
	case 2: // 事务
		resp, err = this.dbTx.ExecPrepare(this.timeout, sqlStr, args...)
	default:
		err = errors.New("unknown mysql mode")
	}
	return resp, err
}

// 标准执行
func (this *SqlBuider) ExecCustom(sqlStr string) (int64, error) {
	var resp int64
	var err error
	switch this.dbType {
	case 1: // Picker
		if this.timeout > 0 {
			resp, err = this.dbPointer.ExecWithTimeout(this.timeout, sqlStr)
		} else {
			resp, err = this.dbPointer.Exec(sqlStr)
		}
	case 2: // 事务
		if this.timeout > 0 {
			resp, err = this.dbTx.ExecWithTimeout(this.timeout, sqlStr)
		} else {
			resp, err = this.dbTx.Exec(sqlStr)
		}
	default:
		err = errors.New("unknown mysql mode")
	}
	return resp, err
}

// 标准查询
func (this *SqlBuider) QueryCustom(sqlStr string) (*DbResult, error) {
	var resp *DbResult = nil
	var err error

	switch this.dbType {
	case 1: // Picker
		if this.timeout > 0 {
			resp, err = this.dbPointer.QueryWithTimeout(this.timeout, this.finalSql)
		} else {
			resp, err = this.dbPointer.Query(this.finalSql)
		}
	case 2: // 事务
		if this.timeout > 0 {
			resp, err = this.dbTx.QueryWithTimeout(this.timeout, this.finalSql)
		} else {
			resp, err = this.dbTx.Query(this.finalSql)
		}
	default:
		err = errors.New("unknown mysql mode")
	}
	return resp, err
}

func (this *SqlBuider) Inc(columns map[string]int) error {
	if this.tablename == "" {
		return fmt.Errorf("unknown table")
	}

	sql := fmt.Sprintf("update %v set ", this.tablename)
	update_column := []string{}
	for column, inc := range columns {
		update_column = append(update_column, fmt.Sprintf("%v = %v + %d", column, column, inc))
	}
	sql += fmt.Sprintf(" %v ", strings.Join(update_column, ","))
	if this.whereStr != "" {
		sql += fmt.Sprintf(" where %v ", this.whereStr)
	}
	if this.orders != "" {
		sql += fmt.Sprintf(" %v ", this.orders)
	}
	if this.limit != "" {
		sql += fmt.Sprintf(" %v ", this.limit)
	}
	this.finalSql = sql
	_, err := this.ExecCustom(this.finalSql)
	return err
}

func (this *SqlBuider) IncFloat(columns map[string]float64) error {
	if this.tablename == "" {
		return fmt.Errorf("unknown table")
	}

	sql := fmt.Sprintf("update %v set ", this.tablename)
	update_column := []string{}
	for column, inc := range columns {
		update_column = append(update_column, fmt.Sprintf("%v = %v + %v", column, column, inc))
	}
	sql += fmt.Sprintf(" %v ", strings.Join(update_column, ","))
	if this.whereStr != "" {
		sql += fmt.Sprintf(" where %v ", this.whereStr)
	}
	if this.orders != "" {
		sql += fmt.Sprintf(" %v ", this.orders)
	}
	if this.limit != "" {
		sql += fmt.Sprintf(" %v ", this.limit)
	}
	this.finalSql = sql
	_, err := this.ExecCustom(this.finalSql)
	return err
}
