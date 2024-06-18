package xmysql

// 错误信息
func (this *Conn) Err() error {
	return this.err
}
