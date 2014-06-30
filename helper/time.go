package helper

import (
	"strings"
	"time"
)

var layout = `2006-01-02 15:04:05`
var loc, _ = time.LoadLocation("Asia/Shanghai")

func StrToTimestamp(stime string) int64 {
	stime = strings.TrimSpace(stime)
	t, err := time.ParseInLocation(layout, stime, loc)
	if err != nil {
		return 0
	}
	return t.Unix()
}

func TimestampToStr(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(layout)
}

func isLeapYear(year int) (result bool) {
	if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
		result = true
	}
	return
}

//计算该天属于该年的第几天
func DayOfYear(year, month, day int) int {
	day_of_year := 0
	if month > 12 || day > 31 {
		return day_of_year
	}
	if isLeapYear(year) {
		day_of_year = (275 * month / 9) - ((month + 9) / 12) + day - 30
	} else {
		day_of_year = (275 * month / 9) - (((month + 9) / 12) << 1) + day - 30
	}
	return day_of_year
}
