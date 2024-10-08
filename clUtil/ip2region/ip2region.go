package ip2region

import (
	"github.com/oschwald/geoip2-golang"
	"net"
)

var xdb *geoip2.Reader
var dbUrl = "https://github.com/cxi7448/cxhttp/raw/refs/heads/main/clUtil/ip2region/GeoIP2-City.mmdb.gz"

// "github.com/lionsoul2014/ip2region/binding/golang/xdb"
func LoadFromFile(dbPath string) error {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return err
	}
	xdb = db
	return nil
}

type XIP struct {
	IP       string
	Country  string
	Err      error
	Province string
	City     string
	Record   *geoip2.City
}

func Get(ip string) XIP {
	result := XIP{}
	if xdb == nil {
		return result
	}
	record, err := xdb.City(net.ParseIP(ip))
	if err != nil {
		result.Err = err
		return result
	}
	result.Country = record.Country.Names["zh-CN"]
	result.Record = record
	if result.Country == "" {
		result.Country = record.Country.Names["en"]
	}
	result.City = record.City.Names["zh-CN"]
	if result.City == "" {
		result.City = record.City.Names["en"]
	}
	if len(record.Subdivisions) > 0 {
		result.Province = record.Subdivisions[0].Names["zh-CN"]
		if result.Province == "" {
			result.Province = record.Subdivisions[0].Names["en"]
		}
	}
	return result
}
