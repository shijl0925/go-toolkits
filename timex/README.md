# timex

`timex` 提供一套更贴近日常业务的时间工具，涵盖时间格式化、时区转换、时间增减、差值计算和 `TimeDelta` 表达。适合接口层、报表层和调度类逻辑使用。

## 安装

```bash
go get github.com/shijl0925/go-toolkits/timex
```

## 核心能力

### 时间格式化

- `FormatStrToTime`
- `FormatTimeToStr`

### 时区处理

- `SetTimeZoneV1`
- `SetTimeZoneV2`

### 时间增减

- `AddMinute`
- `AddHour`
- `AddDay`
- `AddWeek`
- `AddMonth`
- `AddYear`
- `AddTimeDelta`

### 当前时间

- `GetNowDate`
- `GetNowTime`
- `GetNowDateTime`

### 差值与比较

- `GetDaysBetween`
- `GetMonthsBetween`
- `GetYearsBetween`
- `GetHoursBetween`
- `GetMinutesBetween`
- `GetSecondsBetween`
- `GetDurationBetween`
- `GetDurationPretty`
- `Min`
- `Max`

### TimeDelta

- `TimeDelta.Duration()`
- `TimeDelta.Add()`
- `TimeDelta.Subtract()`
- `TimeDelta.Abs()`
- `TimeDelta.String()`

## 快速示例

```go
t, _ := timex.FormatStrToTime("2024-01-01 08:00:00")
nextWeek, _ := timex.AddWeek(t, 1)
prettyDay, prettyHour, prettyMinute, prettySecond := timex.GetDurationPretty(t, nextWeek)

_ = prettyDay
_ = prettyHour
_ = prettyMinute
_ = prettySecond
```

## 默认行为

### 默认格式

`FormatStrToTime` 和 `FormatTimeToStr` 在未传入自定义 layout 时，默认使用：

```go
2006-01-02 15:04:05
```

### TimeDelta

`TimeDelta` 适合表示“周 / 天 / 小时 / 分钟 / 秒 / 毫秒 / 微秒”的组合偏移，例如：

```go
delta := &timex.TimeDelta{Days: 1, Hours: 2, Minutes: 30}
result, err := timex.AddTimeDelta(t, delta)
```

### 溢出保护

`TimeDelta.Duration()` 会检查 `time.Duration` 溢出；当偏移量过大时会返回错误，而不是产生不可预期结果。

## 使用建议

- **业务日期计算**：优先使用 `AddDay`、`AddMonth`、`AddYear` 等语义化函数，代码更易懂。
- **展示友好时长**：使用 `GetDurationPretty` 获取“天/时/分/秒”四元结果。
- **时区转换**：如果你已有标准时区名（如 `Asia/Shanghai`），优先使用 `SetTimeZoneV1`。
- **比较多个时间点**：使用 `Min` / `Max` 简化最早、最晚时间判断。

## 注意事项

- 月份和年份计算受日历规则影响，建议在关键业务场景补充测试。
- `GetHoursBetween`、`GetMinutesBetween`、`GetSecondsBetween` 返回浮点数，适合统计展示。
- `TimeDelta.String()` 内部依赖 `Duration()`；如果发生溢出，会返回 `TimeDelta(overflow)`。

