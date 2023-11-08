package clAuth

import (
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"testing"
)

func TestGetUser(t *testing.T) {
	var auth *AuthInfo
	auth.SetLogin(1, "123123").SetItems(clJson.M{
		"aaa": 111,
	})
	clLog.Info("item: %v", auth.GetUint32("aaa"))
}
