package xmysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"regexp"
	"time"
)

/**
数据库链接
*/

func New() {
	// 自动从环境变量里面读取配置
	os.Getenv()
}

func NewWithOption(option Option) *Conn {
	if option.Charset == "" {
		option.Charset = "utf8"
	}
	host := option.Host
	if !regexp.MustCompile(`:[0-9]+$`).MatchString(option.Host) {
		// 自动补齐端口
		if option.Port == 0 {
			option.Port = 3306
		}
		host += fmt.Sprintf(":%v", option.Port)
	}
	conn := &Conn{}
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=%v", option.User, option.Pass, host, option.Name, option.Charset)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		conn.err = err
	} else {
		conn.db = db
		if option.MaxOpenConns > 0 {
			conn.db.SetMaxOpenConns(option.MaxOpenConns)
		}
		if option.ConnMaxLifetime > 0 {
			conn.db.SetConnMaxLifetime(time.Duration(option.ConnMaxLifetime) * time.Second)
		}
	}
	return conn
}
