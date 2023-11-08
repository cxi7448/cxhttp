package clIP2Region

import (
	"errors"
	"strings"
)

// 进行二次封装，提高易用性
var mRegionBuffer []byte
var mSearcher *Searcher

func InitWithFile(_filePath string) error {
	buff, err := LoadContentFromFile("./ip2region.xdb")
	if err != nil {
		return err
	}

	searcher, err := NewWithBuffer(buff)
	if err != nil {
		return err
	}
	mSearcher = searcher
	return nil
}

// 区域
type Region struct {
	Country  string // 国家
	Province string // 省份
	City     string // 城市
	Provider string // 服务商
}

// 通过IP字符串找到地区
func FindRegionByStr(_ipStr string) (*Region, error) {
	var regionStr, err = mSearcher.SearchByStr(_ipStr)
	if err != nil {
		return nil, err
	}
	var regionItem = strings.Split(regionStr, "|")
	if len(regionItem) != 5 {
		return nil, errors.New("数据格式错误或IP不合法")
	}
	return &Region{
		Country:  regionItem[0],
		Province: regionItem[2],
		City:     regionItem[3],
		Provider: regionItem[4],
	}, nil
}
