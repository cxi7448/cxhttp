package clTime

import (
	"fmt"
	"strings"
	"time"
)

/**
  时间类
*/

// 透过时间日期新建一个时间对象
func NewDate(date string) (*clTimer, error) {

	nowtime := getNowTime(date)

	weeks := 0
	switch strings.ToLower(nowtime.Weekday().String()) {
	case "monday":
		weeks = 0
	case "tuesday":
		weeks = 1
	case "wednesday":
		weeks = 2
	case "thursday":
		weeks = 3
	case "friday":
		weeks = 4
	case "saturday":
		weeks = 5
	case "sunday":
		weeks = 6
	}

	var cltimer = clTimer{
		TimeStamp: uint32(nowtime.Unix()),
		Year:      uint32(nowtime.Year()),
		Month:     uint8(nowtime.Month()),
		Days:      uint8(nowtime.Day()),
		Hour:      uint8(nowtime.Hour()),
		Minuter:   uint8(nowtime.Minute()),
		Second:    uint8(nowtime.Second()),
		Week:      uint8(weeks),
	}
	return &cltimer, nil
}

// 透过时间戳新建一个时间对象
func NewTime(timestamp uint32) (*clTimer, error) {

	nowtime := getTargetTime(timestamp)

	weeks := 0
	switch strings.ToLower(nowtime.Weekday().String()) {
	case "monday":
		weeks = 0
	case "tuesday":
		weeks = 1
	case "wednesday":
		weeks = 2
	case "thursday":
		weeks = 3
	case "friday":
		weeks = 4
	case "saturday":
		weeks = 5
	case "sunday":
		weeks = 6
	}

	var cltimer = clTimer{
		TimeStamp: uint32(nowtime.Unix()),
		Year:      uint32(nowtime.Year()),
		Month:     uint8(nowtime.Month()),
		Days:      uint8(nowtime.Day()),
		Hour:      uint8(nowtime.Hour()),
		Minuter:   uint8(nowtime.Minute()),
		Second:    uint8(nowtime.Second()),
		Week:      uint8(weeks),
	}
	return &cltimer, nil
}




// 获取东八区时间
// @param date string 时间日期格式
// @return *time.Time 返回这个日期格式生成的时间指针
func getNowTime(_date string) *time.Time {
	nowTime := time.Now()
	if _date != "" {
		var err error
		loc, errLoc := time.LoadLocation("Asia/Shanghai")
		if errLoc != nil {
			fmt.Printf("LoadLocation error: %v\n", err)
			return nil
		}
		nowTime, err = time.ParseInLocation("2006-01-02 15:04:05", _date, loc)
		if err != nil {
			return nil
		}

		return &nowTime
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Printf("LoadLocation error: %v\n", err)
		return nil
	}

	utcTime := nowTime.In(loc)
	return &utcTime
}

func getTargetTime(timestamp uint32) *time.Time {
	utc := time.Unix(int64(timestamp), 0)
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil
	}

	utcTime := utc.In(loc)
	return &utcTime
}

// 获取指定时间日期的时间戳
func GetTimeStamp(date string) uint32 {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0
	}

	utc, _ := time.ParseInLocation("2006-01-02 15:04:05", date, loc)
	return uint32(utc.In(loc).Unix())
}

//指定格式获取日期的时间戳
func GetTimeStamp2(date string,format string) uint32 {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0
	}

	utc, _ := time.ParseInLocation(format, date, loc)

	return uint32(utc.In(loc).Unix())
}


// 获取指定时间日期的时间戳單位為 milisecond
func GetTimeStampWithMSec(date string) uint64 {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0
	}
	timedate, _ := time.ParseInLocation("2006-01-02 15:04:05", date, loc)
	return uint64(timedate.UnixNano() / int64(time.Millisecond))
}

// 获取指定时间戳的日期格式
func GetDate(timestamp uint32) string {
	utc := time.Unix(int64(timestamp), 0)
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return "1970-01-01 00:00:00"
	}

	return utc.In(loc).Format("2006-01-02 15:04:05")
}

// 获取指定时间戳的日期格式
func GetDateByFormat(timestamp uint32, format string) string {
	utc := time.Now()
	if timestamp > 0 {
		utc = time.Unix(int64(timestamp), 0)
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return ""
	}

	return utc.In(loc).Format(format)
}



// 检测某个时间是否是在除夕中
// @param timestamp uint32 需要检测的时间戳
func CheckIsChuxi(timestamp uint32) bool {
	targetTime, _ := NewTime(timestamp)

	btime, etime := GetYearTimeBetween(int32(targetTime.Year))

	if btime <= timestamp && etime > timestamp {
		return true
	}
	return false
}

// 获取指定年份的除夕时间区间
// @param year int32 年份
func GetYearTimeBetween(year int32) (/*btime*/uint32, /*etime*/uint32) {
	beginTime := uint32(0)
	switch year {
	case 2018:
		chuxiBegin, _ := NewDate("2018-02-15 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2019:
		chuxiBegin, _ := NewDate("2019-02-04 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2020:
		chuxiBegin, _ := NewDate("2020-01-24 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2021:
		chuxiBegin, _ := NewDate("2021-02-11 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2022:
		chuxiBegin, _ := NewDate("2022-01-31 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2023:
		chuxiBegin, _ := NewDate("2023-01-21 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2024:
		chuxiBegin, _ := NewDate("2024-02-09 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2025:
		chuxiBegin, _ := NewDate("2025-01-28 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2026:
		chuxiBegin, _ := NewDate("2026-02-16 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2027:
		chuxiBegin, _ := NewDate("2027-02-05 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2028:
		chuxiBegin, _ := NewDate("2028-01-25 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2029:
		chuxiBegin, _ := NewDate("2029-02-12 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2030:
		chuxiBegin, _ := NewDate("2030-02-02 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2031:
		chuxiBegin, _ := NewDate("2031-01-22 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2032:
		chuxiBegin, _ := NewDate("2032-02-10 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2033:
		chuxiBegin, _ := NewDate("2033-01-30 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2034:
		chuxiBegin, _ := NewDate("2034-02-18 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2035:
		chuxiBegin, _ := NewDate("2035-02-07 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2036:
		chuxiBegin, _ := NewDate("2036-01-27 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2037:
		chuxiBegin, _ := NewDate("2037-02-14 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	case 2038:
		chuxiBegin, _ := NewDate("2038-02-03 00:00:00")
		beginTime = chuxiBegin.TimeStamp
	}
	return beginTime, beginTime+7*86400
}


// 间隔多少天的除夕
func BetweenChuxiDays(timestamp uint32, target uint32) uint32 {

	// 过滤垃圾参数
	if timestamp > target {
		return 0
	}
	oldtime, _ := NewTime(timestamp)
	newtime, _ := NewTime(target)

	if oldtime.Year == newtime.Year {
		btime, etime := GetYearTimeBetween(int32(oldtime.Year))

		if oldtime.TimeStamp < btime && newtime.TimeStamp < btime {
			// 都在除夕之前
			return 0
		} else if oldtime.TimeStamp > etime && newtime.TimeStamp > etime {
			// 都在除夕之后
			return 0
		}
		// 一前一后
		return 7
	}


	yearBetween := newtime.Year - oldtime.Year
	betweenDays := (yearBetween-1) * 7

	obtime, _ := GetYearTimeBetween(int32(oldtime.Year))

	if oldtime.TimeStamp < obtime {
		// 旧时间在除夕之前 + 7天
		betweenDays += 7
	}

	_, netime := GetYearTimeBetween(int32(newtime.Year))
	if newtime.TimeStamp > netime {
		betweenDays += 7
	}

	return betweenDays
}



// 获取当前时间戳
func GetNowTime() uint32 {
	return uint32(time.Now().Unix())
}

// 获取当前时间格式
func GetNowDateTime() string {
	return GetDate(uint32(time.Now().Unix()))
}

// 获取当前时间格式
func GetNowFormat(_format string) string {
	return GetDateByFormat(uint32(time.Now().Unix()), _format)
}
