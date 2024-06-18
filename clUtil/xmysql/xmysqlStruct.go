package xmysql

import "database/sql"

type Conn struct {
	err    error // 错误信息
	option Option
	db     *sql.DB
}

type Option struct {
	Host            string
	Name            string
	User            string
	Pass            string
	Port            uint32
	Charset         string
	MaxOpenConns    int
	ConnMaxLifetime int
}
