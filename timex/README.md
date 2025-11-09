# timex - Go 时间处理工具包

`timex` 是一个用于处理时间和日期操作的工具包，提供了时间格式化、时区设置、时间计算、时间差计算等丰富的功能。

## 功能列表

### 时间格式化
* `FormatStrToTime(s string, format ...string) (time.Time, error)` - 将字符串转换为时间
* `FormatTimeToStr(t time.Time, format ...string) string` - 将时间转换为字符串

### 时区设置
* `SetTimeZoneV1(t time.Time, timeZone string) (time.Time, error)` - 通过时区名称设置时区
* `SetTimeZoneV2(t time.Time, name string, hour int) (time.Time, error)` - 通过时区偏移量设置时区

### 时间计算
* `AddMinute(t time.Time, minutes int) (time.Time, error)` - 添加或减去分钟数
* `AddHour(t time.Time, hours int) (time.Time, error)` - 添加或减去小时数
* `AddDay(t time.Time, days int) (time.Time, error)` - 添加或减去天数
* `AddWeek(t time.Time, weeks int) (time.Time, error)` - 添加或减去周数
* `AddMonth(t time.Time, months int) (time.Time, error)` - 添加或减去月数
* `AddYear(t time.Time, years int) (time.Time, error)` - 添加或减去年数
* `AddTimeDelta(t time.Time, delta *TimeDelta) (time.Time, error)` - 添加时间差

### 当前时间获取
* [GetNowDate() string](https://github.com/shijl0925/go-toolkits/blob/main/timex/core.go#L177-L180) - 获取当前日期 (yyyy-mm-dd)
* [GetNowTime() string](https://github.com/shijl0925/go-toolkits/blob/main/timex/core.go#L183-L186) - 获取当前时间 (hh-mm-ss)
* [GetNowDateTime() string](https://github.com/shijl0925/go-toolkits/blob/main/timex/core.go#L189-L192) - 获取当前日期时间 (yyyy-mm-dd hh-mm-ss)

### 时间差计算
* `GetDaysBetween(start, end time.Time) int` - 计算两个时间之间的天数差
* `GetMonthsBetween(start, end time.Time) int` - 计算两个时间之间的月数差
* `GetYearsBetween(start, end time.Time) int` - 计算两个时间之间的年数差
* `GetHoursBetween(start, end time.Time) float64` - 计算两个时间之间的小时数差
* `GetMinutesBetween(start, end time.Time) float64` - 计算两个时间之间的分钟数差
* `GetSecondsBetween(start, end time.Time) float64` - 计算两个时间之间的秒数差
* `GetDurationBetween(start, end time.Time) time.Duration` - 计算两个时间之间的持续时间
* `GetDurationPretty(start, end time.Time) (day, hour, minute, second int)` - 以可读格式返回时间差

### 时间比较
* `Min(t1 time.Time, times ...time.Time) time.Time` - 返回给定时间中的最早时间
* `Max(t1 time.Time, times ...time.Time) time.Time` - 返回给定时间中的最晚时间

### 时间差类型 (TimeDelta)
* [TimeDelta](https://github.com/shijl0925/go-toolkits/blob/main/timex/delta.go#L8-L16) - 表示时间差的结构体，包含周、天、小时、分钟、秒、毫秒、微秒等字段
* `Duration() (time.Duration, error)` - 返回可添加到时间的 `time.Duration`
* [Add(td2 *TimeDelta) TimeDelta](https://github.com/shijl0925/go-toolkits/blob/main/timex/delta.go#L83-L93) - 返回两个时间差的和
* `Subtract(td2 *TimeDelta) TimeDelta` - 返回两个时间差的差
* [Abs() TimeDelta](https://github.com/shijl0925/go-toolkits/blob/main/timex/delta.go#L109-L119) - 返回时间差的绝对值
* [String() string](https://github.com/shijl0925/go-toolkits/blob/main/utils_test.go#L20-L22) - 返回时间差的字符串表示

## 安装

```shell
go get github.com/shijl0925/go-toolkits/timex
```


## 使用示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/shijl0925/go-toolkits/timex"
)

func main() {
    // 时间格式化
    t, err := timex.FormatStrToTime("2023-12-25 15:30:45")
    if err != nil {
        panic(err)
    }
    fmt.Println("Parsed time:", t)
    
    // 时间转字符串
    timeStr := timex.FormatTimeToStr(t)
    fmt.Println("Formatted time:", timeStr)
    
    // 时间计算
    newTime, _ := timex.AddDay(t, 7)
    fmt.Println("7 days later:", newTime)
    
    // 使用 TimeDelta
    delta := &timex.TimeDelta{
        Days:  1,
        Hours: 2,
    }
    resultTime, _ := timex.AddTimeDelta(t, delta)
    fmt.Println("After adding delta:", resultTime)
    
    // 时间差计算
    now := time.Now()
    days := timex.GetDaysBetween(t, now)
    fmt.Printf("Days between: %d\n", days)
    
    // 获取当前时间
    fmt.Println("Current date:", timex.GetNowDate())
    fmt.Println("Current time:", timex.GetNowTime())
    fmt.Println("Current datetime:", timex.GetNowDateTime())
}
```
