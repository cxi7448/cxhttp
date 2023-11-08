package clIP2Region

import (
	"fmt"
	"testing"
)

func TestCheckIP(t *testing.T) {

	var ip, e = CheckIP("192.168.0.1")
	fmt.Printf("check ip : %v %v\n", ip, e)
}

func TestLoadHeaderFromFile(t *testing.T) {
	buff, err := LoadContentFromFile("./ip2region.xdb")
	if err != nil {
		fmt.Printf("加载ip数据错误: %v\n", err)
		return
	}

	searcher, err := NewWithBuffer(buff)
	if err != nil {
		fmt.Printf("初始化搜寻器错误: %v\n", err)
		return
	}

	region, err := searcher.SearchByStr("171.214.182.201")
	if err != nil {
		fmt.Printf("搜寻ip地理位置错误: %v\n", err)
		return
	}

	fmt.Printf("地理位置: %v\n", region)
}

func TestInitWithFile(t *testing.T) {
	InitWithFile("./ip2region.xdb")
	r, e := FindRegionByStr("171.214.182.201")
	fmt.Printf("区域信息: %+v %v\n", *r, e)
}
