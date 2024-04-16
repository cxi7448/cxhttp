package xcdn

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	v20180606 "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cdn/v20180606"
)

type TxCDN struct {
	config Config
	client *v20180606.Client
	err    error
}

/*
*
腾讯云CDN管理
*/
func New(config Config) *TxCDN {
	result := &TxCDN{
		config: config,
	}
	client, err := v20180606.NewClientWithSecretId(config.AccessKey, config.SecretKey, config.Region)
	if err != nil {
		result.err = err
		clLog.Error("错误:%v", err)
	} else {
		result.client = client
	}
	return result
}

func (this *TxCDN) PushUrlCache(_url ...string) error {
	if this.client == nil {
		return fmt.Errorf("初始化失败!")
	}
	var url = []*string{}
	for _, val := range _url {
		url = append(url, &val)
	}
	request := v20180606.NewPushUrlsCacheRequest()
	request.Urls = url
	response, err := this.client.PushUrlsCache(request)
	if err != nil {
		clLog.Error("pushurl失败:%v", err)
		return err
	}
	if response.Response != nil {
		return nil
	}
	return fmt.Errorf("预热失败:%v", response.ToJsonString())
}
