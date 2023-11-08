package clAliSDK

import (
	"encoding/json"
	"github.com/cxi7448/cxhttp/clUtil/clCommon"
	"github.com/cxi7448/cxhttp/clUtil/clHttpClient"
	"github.com/cxi7448/cxhttp/clUtil/clTime"
	"time"
)

type RespGetBalance struct {
	Message   string         `json:"Message"`
	RequestId string         `json:"RequestId"`
	Data      GetBalanceData `json:"Data"`
	Code      string         `json:"Code"`
	Success   bool           `json:"success"`
}

type GetBalanceData struct {
	AvailableCashAmount string `json:"AvailableCashAmount"`
	MybankCreditAmount  string `json:"MybankCreditAmount"`
	Currency            string `json:"Currency"`
	AvailableAmount     string `json:"AvailableAmount"`
	CreditAmount        string `json:"CreditAmount"`
}

// 获取余额
func ApiGetBalance(_accessKey, _secretKey string) (error, *RespGetBalance) {

	err, resp := DoReq(Map{
		"Action": "QueryAccountBalance",
	}, _accessKey, _secretKey)
	if err != nil {
		return err, nil
	}
	var data RespGetBalance
	err = json.Unmarshal([]byte(resp), &data)
	if err != nil {
		return err, nil
	}
	return nil, &data
}

// 请求参数
func DoReq(_param Map, _accessKey, _secretKey string) (error, string) {
	var nowTime = uint32(time.Now().Unix())
	var hc = clHttpClient.NewClient("https://business.aliyuncs.com")

	nonceStr := clCommon.GenNonceStr(32)
	globalMap := Map{
		"Format":           "JSON",
		"Version":          "2017-12-14",
		"AccessKeyId":      _accessKey,
		"SignatureMethod":  "HMAC-SHA1",
		"Timestamp":        clTime.GetDateByFormat(nowTime-8*3700, "2006-01-02T15:04:05Z"),
		"SignatureVersion": "1.0",
		"SignatureNonce":   nonceStr,
	}
	for k, v := range _param {
		globalMap[k] = v
	}

	params := GlobalSigned(globalMap, _secretKey)

	for paramKey, paramVal := range params {
		hc.AddParam(paramKey, paramVal)
	}

	hc.SetMethod("GET")
	resp, err := hc.Do()
	if err != nil {
		return err, ""
	}
	return err, resp.Body
}
