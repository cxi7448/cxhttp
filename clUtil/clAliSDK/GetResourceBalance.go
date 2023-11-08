package clAliSDK

import (
	"github.com/cxi7448/cxhttp/clUtil/clLog"
)

func ApiGetResourceBalance(_accessKey, _secretKey string) {
	err, resp := DoReq(Map{
		"Action":       "QueryDPUtilizationDetail",
		"StartTime":    "2022-11-01 12:00:00",
		"EndTime":      "2022-12-01 12:00:00",
		"IncludeShare": false,
	}, _accessKey, _secretKey)
	if err != nil {
		clLog.Error("err: %v", resp)
		return
	}

	clLog.Debug("resp: %v", resp)
}
