package timex

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const defaultTimeFormat = "2006-01-02 15:04:05"

var ErrTimeIsZero = fmt.Errorf("time is zero")
var OverflowError = fmt.Errorf("duration overflow detected")

// FormatStrToTime converts a string to a time.Time
func FormatStrToTime(s string, format ...string) (time.Time, error) {
	if strings.TrimSpace(s) == "" {
		return time.Time{}, fmt.Errorf("input time string is empty or whitespace")
	}

	var f string
	if len(format) > 0 {
		f = strings.TrimSpace(format[0])
	} else {
		f = defaultTimeFormat
	}

	return time.Parse(f, s)
}

// FormatTimeToStr converts a time.Time to a string
func FormatTimeToStr(t time.Time, format ...string) string {
	var f string
	if len(format) > 0 {
		f = strings.TrimSpace(format[0])
	} else {
		f = defaultTimeFormat
	}
	// 可选：添加格式合法性验证（根据业务需求决定是否启用）
	// 示例：尝试使用Parse方法验证格式是否合法
	// 注意：这种方式会带来一定性能开销，适用于调试阶段或低频调用场景
	/*
		if _, err := time.Parse(f, t.Format(f)); err != nil {
			panic("invalid time format: " + f)
		}
	*/

	return t.Format(f)
}

//// Timestamp 时间戳
//func Timestamp(t time.Time) int64 {
//	return t.Unix()
//}
//
//// LocalTime 本地时间
//func LocalTime(timestamp int64) time.Time {
//	return time.Unix(timestamp, 0).Local()
//}

// SetTimeZoneV1 设置时区, timeZone 为时区名称: "Asia/Shanghai"
func SetTimeZoneV1(t time.Time, timeZone string) (time.Time, error) {
	if t.IsZero() {
		return time.Time{}, ErrTimeIsZero
	}

	if timeZone == "" {
		return time.Time{}, fmt.Errorf("time zone name cannot be empty")
	}

	local, err := time.LoadLocation(timeZone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time zone: %s, %s", timeZone, err)
	}

	return t.In(local), nil
}

// SetTimeZoneV2 设置时区, hour 为时区偏移量(-12 - 14), name 为时区名称: "UTC", "CST", "EST"
func SetTimeZoneV2(t time.Time, name string, hour int) (time.Time, error) {
	if t.IsZero() {
		return time.Time{}, ErrTimeIsZero
	}

	// 检查 hour 是否在合理范围内（-12 到 +14 是常见时区偏移范围）
	if hour < -12 || hour > 14 {
		return time.Time{}, fmt.Errorf("invalid hour offset: %d, must be between -12 and +14", hour)
	}

	offset := hour * 3600
	local := time.FixedZone(name, offset)

	return t.In(local), nil
}

// AddMinute add or sub minutes to the time.
func AddMinute(t time.Time, minutes int) (time.Time, error) {
	if t.IsZero() {
		return time.Time{}, ErrTimeIsZero
	}

	// Check for duration overflow
	duration := time.Duration(minutes) * time.Minute
	if (minutes > 0 && duration < 0) || (minutes < 0 && duration > 0) {
		return time.Time{}, OverflowError
	}

	return t.Add(duration), nil
}

// AddHour add or sub hours to the time.
func AddHour(t time.Time, hours int) (time.Time, error) {
	if t.IsZero() {
		return time.Time{}, ErrTimeIsZero
	}

	const maxHours = int64(math.MaxInt64 / int64(time.Hour))
	const minHours = int64(math.MinInt64 / int64(time.Hour))
	if int64(hours) > maxHours || int64(hours) < minHours {
		return time.Time{}, fmt.Errorf("duration overflow in AddHour")
	}

	duration := time.Duration(hours) * time.Hour

	return t.Add(duration), nil
}

// AddDay add or sub days to the time.
func AddDay(t time.Time, days int) (time.Time, error) {
	if t.IsZero() {
		return time.Time{}, ErrTimeIsZero
	}
	// Using AddDate is safer and handles DST changes correctly.
	return t.AddDate(0, 0, days), nil
}

// AddWeek add or sub weeks to the time.
func AddWeek(t time.Time, weeks int) (time.Time, error) {
	if t.IsZero() {
		return time.Time{}, ErrTimeIsZero
	}
	// Using AddDate is safer and handles DST changes correctly.
	return t.AddDate(0, 0, weeks*7), nil
}

// AddMonth add or sub months to the time.
func AddMonth(t time.Time, months int) (time.Time, error) {
	if t.IsZero() {
		return time.Time{}, ErrTimeIsZero
	}

	target := time.Date(t.Year(), t.Month()+time.Month(months), 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	lastDay := time.Date(target.Year(), target.Month()+1, 0, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()).Day()
	day := t.Day()
	if day > lastDay {
		day = lastDay
	}

	return time.Date(target.Year(), target.Month(), day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()), nil
}

// AddYear add or sub year to the time.
func AddYear(t time.Time, years int) (time.Time, error) {
	if t.IsZero() {
		return time.Time{}, ErrTimeIsZero
	}
	return t.AddDate(years, 0, 0), nil
}

func AddTimeDelta(t time.Time, delta *TimeDelta) (time.Time, error) {
	if t.IsZero() {
		return time.Time{}, ErrTimeIsZero
	}
	if delta == nil {
		return time.Time{}, fmt.Errorf("delta is nil")
	}
	d, err := delta.Duration()
	if err != nil {
		// Propagate the error from Duration(), which could be an OverflowError
		return time.Time{}, fmt.Errorf("failed to get duration from TimeDelta: %w", err)
	}
	return t.Add(d), nil
}

// GetNowDate return format yyyy-mm-dd of current date.
func GetNowDate() string {
	now := time.Now()
	return now.Format("2006-01-02")
}

// GetNowTime return format hh-mm-ss of current time.
func GetNowTime() string {
	now := time.Now()
	return now.Format("15:04:05")
}

// GetNowDateTime return format yyyy-mm-dd hh-mm-ss of current datetime.
func GetNowDateTime() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
}

// GetDaysBetween returns the number of days between two times.
func GetDaysBetween(start, end time.Time) int {
	// 确保 end 不早于 start，以得到非负结果或一致的顺序。
	// 如果调用者期望根据原始顺序得到可能为负的结果，则此交换逻辑应移除或调整。
	if end.Before(start) {
		start, end = end, start
	}

	// 将时间标准化到各自日期的午夜，以计算日历天数的差异。
	// 使用各自的原始 Location 来处理可能的时区差异。
	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endDay := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	// 计算标准化日期之间的持续时间，并转换为天数。
	// Sub().Hours()/24 对于已标准化到午夜的时间点是计算日历天差异的常用方法。
	return int(endDay.Sub(startDay).Hours() / 24)
}

// GetMonthsBetween returns the number of months between two times.
func GetMonthsBetween(start, end time.Time) int {
	// 确保 end 不早于 start，防止负数计算
	if end.Before(start) {
		start, end = end, start
	}

	year1, month1, day1 := start.Date()
	year2, month2, day2 := end.Date()

	// 计算总月份数
	months := (year2-year1)*12 + int(month2-month1)

	// 如果结束日的日期小于开始日的日期，则减去一个月，只统计完整的月份差
	if day2 < day1 {
		months--
	}

	return months
}

// GetYearsBetween returns the number of years between two times.
func GetYearsBetween(start, end time.Time) int {
	// 确保 end 不早于 start，防止负数计算
	if end.Before(start) {
		start, end = end, start
	}
	years := end.Year() - start.Year()
	if end.Month() < start.Month() || (end.Month() == start.Month() && end.Day() < start.Day()) {
		years--
	}
	return years
}

// GetHoursBetween returns the number of hours between two times.
func GetHoursBetween(start, end time.Time) float64 {
	// 确保 end 不早于 start，防止负数计算
	if end.Before(start) {
		start, end = end, start
	}
	return end.Sub(start).Hours()
}

// GetMinutesBetween returns the number of minutes between two times.
func GetMinutesBetween(start, end time.Time) float64 {
	// 确保 end 不早于 start，防止负数计算
	if end.Before(start) {
		start, end = end, start
	}
	return end.Sub(start).Minutes()
}

// GetSecondsBetween returns the number of seconds
func GetSecondsBetween(start, end time.Time) float64 {
	// 确保 end 不早于 start，防止负数计算
	if end.Before(start) {
		start, end = end, start
	}
	return end.Sub(start).Seconds()
}

// GetDurationBetween returns the duration between two times, ensuring a non-negative result.
func GetDurationBetween(start, end time.Time) time.Duration {
	// 确保 end 不早于 start，防止负数计算
	if end.Before(start) {
		start, end = end, start
	}
	return end.Sub(start)
}

func GetDurationPretty(start, end time.Time) (day, hour, minute, second int) {
	// 确保 end 不早于 start，防止负数计算
	if end.Before(start) {
		start, end = end, start
	}

	// 计算总秒数（整数）
	seconds := int(end.Sub(start).Seconds())

	// 依次拆分 day, hour, minute, second
	minute, second = seconds/60, seconds%60
	hour, minute = minute/60, minute%60
	day, hour = hour/24, hour%24

	return day, hour, minute, second
}

// Min returns the earliest time among the given times.
func Min(t1 time.Time, times ...time.Time) time.Time {
	minTime := t1
	for _, t := range times {
		if t.Before(minTime) {
			minTime = t
		}
	}
	return minTime
}

// Max returns the latest time among the given times.
func Max(t1 time.Time, times ...time.Time) time.Time {
	maxTime := t1
	for _, t := range times {
		if t.After(maxTime) {
			maxTime = t
		}
	}
	return maxTime
}
