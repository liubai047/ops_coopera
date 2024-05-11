package ctime

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// AddTimeDuration 从指定时间计算，增加/减少多少时间。支持尾缀格式:"ns"、"us" "µs"、"ms"、"s"、"m"、"h"、"day",
func AddTimeDuration(baseTime time.Time, str string) (time.Time, error) {
	var allowSuffix = []string{"ns", "us", "µs", "ms", "s", "m", "h", "day"}
	var allowStr = false
	var unSuffixVal = ""
	// 判断后缀是否合法
	for _, v := range allowSuffix {
		if strings.HasSuffix(str, v) {
			allowStr = true
			unSuffixVal = strings.TrimSuffix(str, v)
			break
		}
	}
	if !allowStr {
		return time.Time{}, errors.New("非法后缀")
	}
	// 去掉后缀后，判断字符串是否是int类型
	intVal, err := strconv.Atoi(unSuffixVal)
	if err != nil {
		return time.Time{}, errors.New("时间格式有误,非数字类型")
	}
	// 自定义原生方法(目前只有day)
	if strings.HasSuffix(str, "day") {
		return baseTime.Add(time.Duration(intVal*24*3600) * time.Second), nil
	}
	// 调用官方原生方法
	dura, err := time.ParseDuration(str)
	if err != nil {
		return time.Time{}, err
	}
	return baseTime.Add(dura), nil
}

// GetDayUnix 获取某day对应的时间戳值（ps:1day = 86400）
func GetDayUnix(str string) (int64, error) {
	if !strings.HasSuffix(str, "day") {
		return 0, errors.New("字符格式有误")
	}
	unSuffixVal := strings.TrimSuffix(str, "day")
	intVal, err := strconv.Atoi(unSuffixVal)
	if err != nil {
		return 0, errors.New("时间格式有误,非数字类型")
	}
	return int64(intVal * 24 * 3600), nil
}

// GetDayStartTime 获取某日起始时间点
func GetDayStartTime(tm time.Time) time.Time {
	return time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location())
}

// GetDayEndTime 获取某日结束时间点
func GetDayEndTime(tm time.Time) time.Time {
	return time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 59, 59, 999999999, tm.Location())
}

// GetTime 获取时间
func GetTime(tm time.Time, format ...string) string {
	var f = time.DateTime
	if len(format) > 0 {
		f = format[0]
	}
	return tm.Format(f)
}

// TmStrToTmStamp 时间格式字符串转为时间戳
func TmStrToTmStamp(tmStr string, format ...string) (int64, error) {
	fm := time.DateTime
	if len(format) > 0 {
		fm = format[0]
	}
	tm, err := time.Parse(fm, tmStr)
	if err != nil {
		return 0, err
	}
	return tm.Unix(), nil
}

// TmStampToTmStr 时间戳改为时间格式
func TmStampToTmStr(timestamp int64, format ...string) string {
	tm := time.Unix(timestamp, 0)
	return GetTime(tm, format...)
}

// StampToTime 时间戳转为时间格式
func StampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// TimeStrToTime 时间格式字符串转为时间格式
func TimeStrToTime(tmStr string, format ...string) (time.Time, error) {
	fm := time.DateTime
	if len(format) > 0 {
		fm = format[0]
	}
	tm, err := time.Parse(fm, tmStr)
	if err != nil {
		return time.Time{}, err
	}
	return tm, nil
}
