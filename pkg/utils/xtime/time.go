package xtime

import "time"

func GetNowMillisecond() int64 {
	return time.Now().UnixNano() / 1000000
}

func GetNowSecond() int64 {
	return time.Now().Unix()
}

func GetTodaySecond(loc *time.Location) int64 {
	y, m, d := time.Now().In(loc).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, loc).Unix()
}
