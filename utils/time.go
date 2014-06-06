package utils

import "time"

//返回距离午夜的秒数
func GetMidnightSeconds() int64 {
	now := time.Now()
	midnight := (23-now.Hour())*60*60 + (59-now.Minute())*60 + 59 - now.Second() + 1
	return int64(midnight)
}
