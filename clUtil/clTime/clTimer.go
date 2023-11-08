package clTime

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clCommon"
	"math"
	"time"
)

// 偏移指定秒数
func (this *clTimer) AfterSecond(sec uint32) (*clTimer, error) {

	if sec == 0 {
		return this, nil
	}

	nowtime, _ := NewTime(this.TimeStamp + sec)
	return nowtime, nil
}

// 获取星期几的文本表示方式
func (this *clTimer) GetWeekStr() string {
	switch this.Week {
	case 0:
		return "星期一"
	case 1:
		return "星期二"
	case 2:
		return "星期三"
	case 3:
		return "星期四"
	case 4:
		return "星期五"
	case 5:
		return "星期六"
	case 6:
		return "星期日"
	default:
		return "NULL"
	}
}

// 获取本月开始时间
// @return uint32 当前时间的本月1号0时0分0秒的时间戳
// @return string 当前时间的本月1号0时0分0秒的日期时间格式
func (this *clTimer) GetCurMonth() ( /*timestamp*/ uint32 /*datestr*/, string) {
	dateformat := fmt.Sprintf("%04v-%02v-01 00:00:00", this.Year, this.Month)
	targetTime := getNowTime(dateformat)
	return uint32(targetTime.Unix()), targetTime.Format("2006-01-02 15:04:05")
}

// 获取本周的开始时间
// @return uint32 当前时间的本周星期一的时间戳
// @return string 当前时间的本周星期一的日期时间戳
func (this *clTimer) GetCurWeek() ( /*timestamp*/ uint32 /*datestr*/, string) {
	dateformat := fmt.Sprintf("%04v-%02v-%02v 00:00:00", this.Year, this.Month, this.Days)
	todayTime := getNowTime(dateformat)
	targetTime := todayTime.Add(-time.Duration(this.Week) * 24 * time.Hour)
	return uint32(targetTime.Unix()), targetTime.Format("2006-01-02 15:04:05")
}

// 获取指定跨度的月份时间戳区间
// @param offset int 跨度偏移， 0为当前月份时间周期
// @return uint32 指定时间区间起始时间
// @return uint32 指定时间区间结束时间
func (this *clTimer) GetMonthBetween(offset int) (uint32, uint32) {

	var beginMonth = int(this.Month) + offset
	var beginYear = int(this.Year)
	if beginMonth < 1 || beginMonth > 12 {
		if offset < 0 {
			beginYear = int(this.Year) + beginMonth/12 - 1
			beginMonth = 12 - (int(math.Abs(float64(beginMonth))) % 12)
		} else {
			beginMonth = (int(this.Month) + offset) % 12
			beginYear = int(this.Year) + (beginMonth-1)/12
		}
	}

	var endMonth = beginMonth + 1
	var endYear = beginYear
	if endMonth > 12 {
		endMonth = 1
		endYear++
	}

	begins := GetTimeStamp(fmt.Sprintf("%04v-%02v-01 00:00:00", beginYear, beginMonth))
	ends := GetTimeStamp(fmt.Sprintf("%04v-%02v-01 00:00:00", endYear, endMonth))
	return begins, ends - 1
}

// 跨年
func (this *clTimer) GetYearBetween(offset int) (uint32, uint32) {
	beginYear := int(this.Year) - offset
	begins := GetTimeStamp(fmt.Sprintf("%04v-01-01 00:00:00", beginYear))
	ends := GetTimeStamp(fmt.Sprintf("%04v-01-01 00:00:00", beginYear+1))
	return begins, ends - 1
}

// 获取指定跨度的月份时间戳区间
// @param offset int 跨度偏移， 0为当前月份时间周期
// @return uint64 指定时间区间起始时间
// @return uint64 指定时间区间结束时间
func (this *clTimer) GetMonthBetweenWithMSec(offset int) (uint64, uint64) {

	var beginMonth = int(this.Month) + offset
	var beginYear = int(this.Year)
	if beginMonth < 1 || beginMonth > 12 {
		if offset < 0 {
			beginYear = int(this.Year) + beginMonth/12 - 1
			beginMonth = 12 - (beginMonth % 12)
		} else {
			beginMonth = (int(this.Month) + offset) % 12
			beginYear = int(this.Year) + (beginMonth-1)/12
		}
	}

	var endMonth = beginMonth + 1
	var endYear = beginYear
	if endMonth > 12 {
		endMonth = 1
		endYear++
	}

	begins := GetTimeStampWithMSec(fmt.Sprintf("%04v-%02v-01 00:00:00", beginYear, beginMonth))
	ends := GetTimeStampWithMSec(fmt.Sprintf("%04v-%02v-01 00:00:00", endYear, endMonth))
	return begins, ends - 1000
}

// 获取指定跨度的小时时间戳区间
// @param offset int 跨度偏移， 0为当前小时时间周期
// @return uint32 指定时间区间起始时间
// @return uint32 指定时间区间结束时间
func (this *clTimer) GetHourBetween(offset int) (uint32, uint32) {

	nowTime := uint32(this.TimeStamp)
	beginTime := nowTime - (nowTime % 3600)
	endTime := beginTime + 3599

	return beginTime, endTime
}

// 获取指定跨度的周时间戳区间
// @param offset int 跨度偏移， 0为当前周时间周期
// @return uint32 指定时间区间起始时间
// @return uint32 指定时间区间结束时间
func (this *clTimer) GetWeekBetween(offset int) (uint32, uint32) {

	beginTime, _ := NewTime(this.TimeStamp + uint32(offset*7*86400))
	endTime, _ := NewTime(this.TimeStamp + uint32((offset+1)*7*86400))

	begins := GetTimeStamp(fmt.Sprintf("%04v-%02v-%02v 00:00:00", beginTime.Year, beginTime.Month, beginTime.Days))
	if beginTime.Week > 0 {
		begins = begins - uint32(endTime.Week)*86400
	}
	ends := GetTimeStamp(fmt.Sprintf("%04v-%02v-%02v 00:00:00", endTime.Year, endTime.Month, endTime.Days))
	if endTime.Week > 0 {
		ends = ends - uint32(endTime.Week)*86400
	}
	return begins, ends - 1
}

// 获取指定跨度的周时间戳区间
// @param offset int 跨度偏移， 0为当前周时间周期
// @return uint64 指定时间区间起始时间
// @return uint64 指定时间区间结束时间
func (this *clTimer) GetWeekBetweenWithMSec(offset int) (uint64, uint64) {

	beginTime, _ := NewTime(this.TimeStamp + uint32(offset*7*86400))
	endTime, _ := NewTime(this.TimeStamp + uint32((offset+1)*7*86400))

	begins := GetTimeStampWithMSec(fmt.Sprintf("%04v-%02v-%02v 00:00:00", beginTime.Year, beginTime.Month, beginTime.Days))

	if beginTime.Week > 0 {
		begins = begins - uint64(endTime.Week)*86400000
	}
	ends := GetTimeStampWithMSec(fmt.Sprintf("%04v-%02v-%02v 00:00:00", endTime.Year, endTime.Month, endTime.Days))
	if endTime.Week > 0 {
		ends = ends - uint64(endTime.Week)*86400000
	}
	return begins, ends - 1000
}

// 获取指定跨度的天时间戳区间
// @param offset int 跨度偏移， 0为当前周时间周期
// @return uint32 指定时间区间起始时间
// @return uint32 指定时间区间结束时间
func (this *clTimer) GetDayBetween(offset int) (uint32, uint32) {

	begins := GetTimeStamp(fmt.Sprintf("%04v-%02v-%02v 00:00:00", this.Year, this.Month, this.Days))
	btime := uint32(int32(begins) + int32(offset*86400))
	return btime, btime + 86400 - 1
}

// 获取指定跨度的天时间戳区间單位為 milisecond
// @param offset int 跨度偏移， 0为当前周时间周期
// @return uint64 指定时间区间起始时间
// @return uint64 指定时间区间结束时间
func (this *clTimer) GetDayBetweenWithMSec(offset int) (uint64, uint64) {

	begins := GetTimeStampWithMSec(fmt.Sprintf("%04v-%02v-%02v 00:00:00", this.Year, this.Month, this.Days))
	btime := uint64(int64(begins) + int64(offset*86400000))
	return btime, btime + (86400-1)*1000
}

// 格式化当前时间文本
func (this *clTimer) FormatBy(_layout string) string {
	return GetDateByFormat(this.TimeStamp, _layout)
}

// 格式化当前时间文本的日期
func (this *clTimer) GetDate() string {
	return GetDateByFormat(this.TimeStamp, "2006-01-02")
}

// 格式化当前时间文本的时间
func (this *clTimer) GetTime() string {
	return GetDateByFormat(this.TimeStamp, "15:04:05")
}

// 格式化当前时间的日期时间
func (this *clTimer) GetDateTime() string {
	return GetDateByFormat(this.TimeStamp, "2006-01-02 15:04:05")
}

// 格式化当前时间的日期的数字格式
func (this *clTimer) GetSDate() uint32 {
	return clCommon.Uint32(GetDateByFormat(this.TimeStamp, "20060102"))
}
