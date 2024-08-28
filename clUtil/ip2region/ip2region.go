package ip2region

import (
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"strings"
)

var xdbSearcher *xdb.Searcher

// "github.com/lionsoul2014/ip2region/binding/golang/xdb"
func LoadFromFile(dbPath string) error {
	cBuff, err := xdb.LoadContentFromFile(dbPath)
	if err != nil {
		return err
	}
	searcher, err := xdb.NewWithBuffer(cBuff)
	if err != nil {
		return err
	}
	xdbSearcher = searcher
	return nil
}

type Ip struct {
	IP       string
	Country  string
	Err      error
	Province string
	City     string
	Origin   string
}

func Get(ip string) Ip {
	result := Ip{IP: ip, Country: "中国"}
	// 中国|0|香港|0|联通
	// 国家|区域|省份|城市|ISP
	if xdbSearcher == nil {
		result.Err = fmt.Errorf("xdb error")
		return result
	}
	res, err := xdbSearcher.SearchByStr(ip)
	if err == nil {
		rows := strings.Split(res, "|")
		if len(rows) > 0 {
			result.Country = rows[0]
		}
		if len(rows) > 2 {
			result.Province = rows[2]
		}
		if len(rows) > 3 {
			result.City = rows[3]
		}
		result.Origin = res
	} else {
		result.Err = err
	}
	return result
}

func (this Ip) GetLang() string {
	switch this.Country {
	case "中国":
		return "CN"
	case "泰国":
		return "THA"
	default:
		return "EN"
	}
}
