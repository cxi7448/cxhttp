package tinypng

import (
	"bytes"
	"fmt"
	"github.com/cxi7448/cxhttp/clCommon"
	"github.com/cxi7448/cxhttp/clUtil/clFile"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/cxi7448/cxhttp/clUtil/clTime"
	"io/ioutil"
	"sync"
	"time"
)

var filename = "tinypng.dat"

type Account struct {
	Email  string `json:"email"`
	Apikey string `json:"apikey"`
	Month  uint32 `json:"month"` // 当前月份
	Times  int    `json:"times"` // 次数
}

var accountPools = []Account{
	{Email: "yfeiyu174@gmail.com", Apikey: "98RtDTNzvmJBl1m2406NP6RqwQPLhJdQ"},
	{Email: "chax70308@gmail.com", Apikey: "VPTjfGcB8sV3gzrGfwhG2CCdk0vYVLMx"},
	{Email: "da3208628@gmail.com", Apikey: "CRmZ6BScMVwJfkWKNqxLYQgh312pWCQK"},
	{Email: "lanniao615@gmail.com", Apikey: "CWKhrC44t6NSy8tgF22bqG8RqBYBsccz"},
	{Email: "hai44geg@gmail.com", Apikey: "rTtKNmvlLJTM0sy45dlkxZYMdkdvHxMr"},
	{Email: "sx195788@gmail.com", Apikey: "8zV1TC6DxczsxKC6pTcgrTLNQCqSzj3D"},
	{Email: "xinx27176@gmail.com", Apikey: "4N2WNZn2sJD6QKQby5l1DzvSzcbwt5Lq"},
	{Email: "xiao93990@gmail.com", Apikey: "GYG6KV7ZcW6R6VfXM284Xxw7TDHDVxXy"},
	{Email: "fyun9754@gmail.com", Apikey: "SXCwH5s9cMbJ9Z5QSSZT5d9ZWhqxsgLc"},
	{Email: "lxingqi5@gmail.com", Apikey: "t95bwj4bL69hld9j9RqPP0lMZFMQ6S5L"},
	{Email: "zhou42614@gmail.com", Apikey: "DLvNgMJBJpzzvxg39XKpJBr0LxZrTphD"},
	{Email: "xkong9455@gmail.com", Apikey: "FTbRzyG2Px59hsvZp1Js58LLLF3NrrFK"},
	{Email: "lyue92705@gmail.com", Apikey: "N4ppt41nqDLLRttQbHWYhf7mzZ5BWVMk"},
	{Email: "linb57874@gmail.com", Apikey: "83qwss2sGpvtpzbdfYbrCBbK4VhyMPlX"},
	{Email: "tun025375@gmail.com", Apikey: "DCSW386F2r9zx5fqnFmt3WGNG0fN787M"},
	{Email: "haixing656@gmail.com", Apikey: "CQMqyJrhVqvLlFq1QcJ24FnZfPgjCHbg"},
	{Email: "lanse773@gmail.com", Apikey: "4Hl0wdjsGgbVS2KsYxJ9tjLMKr7S4y04"},
	{Email: "yse88786@gmail.com", Apikey: "r21FlN8s9hGFMV8gHYtwy8RpbS55P1dX"},
}
var locker sync.RWMutex

func LoadFile(filePath string) {
	if !clFile.IsFile(filePath) {
		return
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		clLog.Error("加载tinypng账号失败:%v", err)
		return
	}
	filename = filePath
	split := []byte("\n")
	if bytes.Contains(content, []byte("\r")) {
		split = []byte("\r\n")
	}
	var accounts = []Account{}
	for _, row := range bytes.Split(content, split) {
		rows := bytes.Split(row, []byte("|"))
		if len(rows) != 2 {
			continue
		}
		var month uint32
		var times int
		splits := bytes.Split(rows[1], []byte(","))
		if len(splits) > 1 {
			month = clCommon.Uint32(string(splits[1]))
		}
		if len(splits) > 2 {
			times = clCommon.Int(string(splits[2]))
		}
		accounts = append(accounts, Account{
			Email:  string(splits[0]),
			Apikey: string(rows[0]),
			Month:  month,
			Times:  times,
		})
	}
	locker.Lock()
	defer locker.Unlock()
	accountPools = accounts
}

func GetOne() *Account {
	locker.RLock()
	defer locker.RUnlock()
	month := clCommon.Uint32(clTime.GetDateByFormat(uint32(time.Now().Unix()), "01"))
	for _, row := range accountPools {
		if row.Month == month && row.Times == -1 {
			continue
		}
		return &row
	}
	return nil
}

func (this *Account) SetTimes(times int) {
	locker.Lock()
	defer locker.Unlock()
	var content = ""
	for key, row := range accountPools {
		if row.Email == this.Email {
			accountPools[key].Month = clCommon.Uint32(clTime.GetDateByFormat(uint32(time.Now().Unix()), "01"))
			if times == -1 {
				accountPools[key].Times = -1
			} else {
				accountPools[key].Times += 1
			}
		}
		content += fmt.Sprintf("%v|%v,%v,%v\n", row.Apikey, row.Email, accountPools[key].Month, accountPools[key].Times)
	}
	ioutil.WriteFile(filename, []byte(content), 0700)
}
