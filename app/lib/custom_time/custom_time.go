package custom_time

import (
	_ "fmt"
	"time"
)

func GetCustomTimeContextKey() string {
	return "custom_time"
}

func Now() time.Time {
	//[todo]: 検証用にdebugできるように?
	// debug //
	//return time.Date(2017, 7, 19, 23, 59, 0, 1, time.Local)
	//return time.Date(2017, 9, 21, 12, 0, 0, 1, time.Local)
	///////////

	return time.Now()
}

func GetDayStartTime(t time.Time) time.Time {
	// debug //
	//t = time.Date(2017, 5, 30, 12, 0, 0, 1, time.Local)
	///////////

	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func GetTimeWithoutMs(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()+1, 0, t.Location())
}

func GetTimeDayAfter(t time.Time, day_after int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day()+day_after, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

func GetDayOf(year int, month int, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Now().Location())
}

func GetTimeWithoutMsFromUnix(unixtime int64) time.Time {
	return GetTimeWithoutMs(time.Unix(unixtime, 0))
}

/*
func GetDayStartTime2(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
*/
