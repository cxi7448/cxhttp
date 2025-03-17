package httpserver

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clResponse"
	"github.com/cxi7448/cxhttp/core/clAuth"
	"github.com/cxi7448/cxhttp/core/rule"
	"testing"
)

func TestCallHandler(t *testing.T) {

	rule.AddRule(rule.Rule{
		Request: "request",
		Name:    "paycallback",
		Params:  nil,
		CallBack: func(_auth *clAuth.AuthInfo, _param *rule.HttpParam, _server *rule.ServerParam) string {
			fmt.Println(_param.ToMap())
			fmt.Println(_server.RawData)
			return clResponse.Success()
		},
		CacheExpire:   0,
		CacheType:     0,
		CacheKeyParam: nil,
		Login:         false,
		Method:        "",
		RespContent:   "",
	})
	StartServer(9001)
}
