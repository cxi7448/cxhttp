package clRedis

import (
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"testing"
)

func TestNew(t *testing.T) {

	redis, err := New("127.0.0.1:6379", "", "")
	if err != nil {
		clLog.Error("连接错误: %v", err)
		return
	}

	if redis.SetEx("hello", "1", 100) {
		clLog.Debug("设置成功!!")
	} else {
		clLog.Debug("设置失败!!")
	}

	if redis.SetExpire("hello", 200) {
		clLog.Debug("设置时效成功!!")
	} else {
		clLog.Debug("设置时效失败!!")
	}

	clLog.Debug("设置redis...")
}
