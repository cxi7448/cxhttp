package clLog

import "testing"

func TestInfo(t *testing.T) {
	SetVersion("v1.0.0")
	Info("测试Info: %v", 1)
}
