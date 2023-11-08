package clAliSDK

import (
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"testing"
)

func TestApiGetBalance(t *testing.T) {
	err, data := ApiGetBalance("LTAI5tBHtn33tn1sHN2aC8D9", "naBSo2fFwtHr9Oqc4hY1AmV6WpFADU")
	if err != nil {
		clLog.Error("获取余额错误: %v", err)
		return
	}

	clLog.Debug("账号: %v 余额: %v", "video-best", data.Data.AvailableCashAmount)
}

func TestApiGetResourceBalance(t *testing.T) {

	ApiGetResourceBalance("LTAI5tBHtn33tn1sHN2aC8D9", "naBSo2fFwtHr9Oqc4hY1AmV6WpFADU")
}
