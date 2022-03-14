package utils

import (
	"time"
)

func StrToTime(date string) time.Time {
	local, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", date, local)
	return t
}

func StrToDate(date string)time.Time{
	local, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02", date, local)
	return t
}
