package timex_test

import (
	"github.com/shijl0925/go-toolkits/timex"
	"testing"
	"time"
)

func Test_StringFormatTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		format   []string
		wantErr  bool
		expected time.Time
	}{
		{
			name:    "Empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:     "Valid default format",
			input:    "2024-01-01 12:00:00",
			wantErr:  false,
			expected: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "Custom format valid",
			input:    "2024/01/01 12:00:00",
			format:   []string{"2006/01/02 15:04:05"},
			wantErr:  false,
			expected: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:    "Invalid time string",
			input:   "invalid-time",
			wantErr: true,
		},
		{
			name:    "Mismatched custom format",
			input:   "2024-01-01 12:00:00",
			format:  []string{"2006/01/02 15:04:05"},
			wantErr: true,
		},
		{
			name:    "Whitespace input",
			input:   "   ",
			wantErr: true,
		},
		{
			name:     "Empty format string",
			input:    "2024-01-01 12:00:00",
			format:   []string{},
			wantErr:  false,
			expected: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "Custom format valid 2",
			input:    "2024-01-01T12:00:00",
			format:   []string{"2006-01-02T15:04:05"},
			wantErr:  false,
			expected: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := timex.FormatStrToTime(tt.input, tt.format...)
			if tt.wantErr {
				if err == nil {
					t.Errorf("StringFormatTime() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("StringFormatTime() error = %v, wantErr %v", err, tt.wantErr)
				}
				if !got.Equal(tt.expected) {
					t.Errorf("StringFormatTime() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func Test_TimeFormatString(t *testing.T) {
	// 构造一个固定时间用于测试
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		input  time.Time
		format []string
		//wantErr  bool
		expected string
	}{
		{
			name:     "Use custom format",
			input:    testTime,
			format:   []string{"2006-01-02"},
			expected: "2024-01-01",
		},
		{
			name:     "Use default format",
			input:    testTime,
			format:   nil,
			expected: "2024-01-01 12:00:00",
		},
		{
			name:     "Empty format string",
			input:    testTime,
			format:   []string{""},
			expected: "",
		},
		{
			name:     "Format with spaces",
			input:    testTime,
			format:   []string{"  2006/01/02  "},
			expected: "2024/01/01",
		},
		{
			name:     "Use custom format 2",
			input:    testTime,
			format:   []string{"2006-01-02T15:04:05"},
			expected: "2024-01-01T12:00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timex.FormatTimeToStr(tt.input, tt.format...)
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGetDurationBetween(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected time.Duration
	}{
		{
			name:     "TC01 - End after Start",
			start:    now,
			end:      now.Add(2 * time.Hour),
			expected: 2 * time.Hour,
		},
		{
			name:     "TC02 - Start after End",
			start:    now.Add(3 * time.Hour),
			end:      now,
			expected: 3 * time.Hour,
		},
		{
			name:     "TC03 - Equal Times",
			start:    now,
			end:      now,
			expected: 0,
		},
		{
			name:     "TC04 - Exactly one day difference",
			start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			expected: 24 * time.Hour,
		},
		{
			name:     "TC05 - Cross year/month/day difference",
			start:    time.Date(2023, 12, 31, 23, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 1, 1, 30, 0, 0, time.UTC),
			expected: 2*time.Hour + 30*time.Minute,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := timex.GetDurationBetween(tc.start, tc.end)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

// 定义辅助函数用于创建固定时间点
func newTime(year int, month time.Month, day, hour, min, sec int) time.Time {
	return time.Date(year, month, day, hour, min, sec, 0, time.UTC)
}

func TestGetDurationPretty(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected [4]int // [day, hour, minute, second]
	}{
		{
			name:     "TC1 - 1秒差异",
			start:    newTime(2024, time.January, 1, 0, 0, 0),
			end:      newTime(2024, time.January, 1, 0, 0, 1),
			expected: [4]int{0, 0, 0, 1},
		},
		{
			name:     "TC2 - 61秒差异",
			start:    newTime(2024, time.January, 1, 0, 0, 0),
			end:      newTime(2024, time.January, 1, 0, 1, 1),
			expected: [4]int{0, 0, 1, 1},
		},
		{
			name:     "TC3 - 整小时差异",
			start:    newTime(2024, time.January, 1, 0, 0, 0),
			end:      newTime(2024, time.January, 1, 1, 0, 0),
			expected: [4]int{0, 1, 0, 0},
		},
		{
			name:     "TC4 - 整天差异",
			start:    newTime(2024, time.January, 1, 0, 0, 0),
			end:      newTime(2024, time.January, 2, 0, 0, 0),
			expected: [4]int{1, 0, 0, 0},
		},
		{
			name:     "TC5 - 多单位混合差异",
			start:    newTime(2024, time.January, 1, 0, 0, 0),
			end:      newTime(2024, time.July, 23, 1, 1, 1),
			expected: [4]int{204, 1, 1, 1},
		},
		{
			name:     "TC6 - end < start，自动交换",
			start:    newTime(2024, time.January, 1, 1, 0, 0),
			end:      newTime(2024, time.January, 1, 0, 0, 0),
			expected: [4]int{0, 1, 0, 0},
		},
		{
			name:     "TC7 - 相同时间",
			start:    newTime(2024, time.January, 1, 12, 0, 0),
			end:      newTime(2024, time.January, 1, 12, 0, 0),
			expected: [4]int{0, 0, 0, 0},
		},
		{
			name:     "TC8 - 多单位混合差异",
			start:    newTime(2024, time.January, 23, 12, 23, 1),
			end:      newTime(2024, time.July, 1, 1, 1, 1),
			expected: [4]int{159, 12, 38, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, h, m, s := timex.GetDurationPretty(tt.start, tt.end)
			if d != tt.expected[0] || h != tt.expected[1] || m != tt.expected[2] || s != tt.expected[3] {
				t.Errorf("expected (%d, %d, %d, %d), got (%d, %d, %d, %d)",
					tt.expected[0], tt.expected[1], tt.expected[2], tt.expected[3],
					d, h, m, s)
			}
		})
	}
}

func TestAddMinute(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		inputTime time.Time
		minutes   int
		wantTime  time.Time
		wantErr   bool
	}{
		{
			name:      "Zero Time Input",
			inputTime: time.Time{},
			minutes:   10,
			wantTime:  time.Time{},
			wantErr:   true,
		},
		{
			name:      "Add Positive Minutes",
			inputTime: baseTime,
			minutes:   30,
			wantTime:  baseTime.Add(30 * time.Minute),
			wantErr:   false,
		},
		{
			name:      "Subtract Minutes",
			inputTime: baseTime,
			minutes:   -60,
			wantTime:  baseTime.Add(-60 * time.Minute),
			wantErr:   false,
		},
		{
			name:      "Duration Overflow - Large Positive",
			inputTime: baseTime,
			minutes:   1e15, // This may cause overflow on 32-bit systems
			wantTime:  time.Time{},
			wantErr:   true,
		},
		{
			name:      "Duration Overflow - Large Negative",
			inputTime: baseTime,
			minutes:   -1e15,
			wantTime:  time.Time{},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, err := timex.AddMinute(tt.inputTime, tt.minutes)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !gotTime.Equal(tt.wantTime) {
					t.Errorf("Expected %v, got %v", tt.wantTime, gotTime)
				}
			}
		})
	}
}

func TestAddHour(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		inputTime time.Time
		hours     int
		wantTime  time.Time
		wantErr   bool
	}{
		{
			name:      "Zero Time Input",
			inputTime: time.Time{},
			hours:     10,
			wantTime:  time.Time{},
			wantErr:   true,
		},
		{
			name:      "Add Positive Minutes",
			inputTime: baseTime,
			hours:     12,
			wantTime:  baseTime.Add(12 * time.Hour),
			wantErr:   false,
		},
		{
			name:      "Subtract Minutes",
			inputTime: baseTime,
			hours:     -24,
			wantTime:  baseTime.Add(-24 * time.Hour),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, err := timex.AddHour(tt.inputTime, tt.hours)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !gotTime.Equal(tt.wantTime) {
					t.Errorf("Expected %v, got %v", tt.wantTime, gotTime)
				}
			}
		})
	}
}

func TestGetDaysBetween(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected int
	}{
		{
			name:     "Normal case: one day apart",
			start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "Less than a full day",
			start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "Reverse order with two days difference",
			start:    time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
			expected: 2,
		},
		{
			name:     "Same time",
			start:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "Cross leap year date (Feb to Mar)",
			start:    time.Date(2020, 2, 28, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
			expected: 2,
		},
		{
			name:     "Multiple full days",
			start:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := timex.GetDaysBetween(tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("Expected %d days, got %d", tt.expected, result)
			}
		})
	}
}

func TestGetMonthsBetween(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected int
	}{
		{
			name:     "Same date",
			start:    time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "Same month, end day > start day",
			start:    time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "Same month, end day < start day",
			start:    time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 5, 10, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "Cross month, end day < start day",
			start:    time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "Cross month, end day >= start day",
			start:    time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "Cross year, end day >= start day",
			start:    time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			expected: 13,
		},
		{
			name:     "Cross year, end day < start day",
			start:    time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC),
			expected: 12,
		},
		{
			name:     "End before start should swap",
			start:    time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			expected: 12,
		},
		{
			name:     "End of month to next month first day",
			start:    time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "Leap year Feb 29 to next year Feb 28",
			start:    time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2021, 2, 28, 0, 0, 0, 0, time.UTC),
			expected: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := timex.GetMonthsBetween(tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("Expected %d, but got %d", tt.expected, result)
			}
		})
	}
}

func TestGetYearsBetween(t *testing.T) {
	tests := []struct {
		name   string
		start  time.Time
		end    time.Time
		expect int
	}{
		{
			name:   "End after start, same month and day",
			start:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			end:    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			expect: 3,
		},
		{
			name:   "End after start, end day is later",
			start:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			end:    time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			expect: 3,
		},
		{
			name:   "End after start, end day is earlier",
			start:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			end:    time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC),
			expect: 2,
		},
		{
			name:   "End after start, same month but earlier day",
			start:  time.Date(2020, 3, 15, 0, 0, 0, 0, time.UTC),
			end:    time.Date(2023, 3, 14, 0, 0, 0, 0, time.UTC),
			expect: 2,
		},
		{
			name:   "Start after end (auto swap), same month and day",
			start:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			end:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expect: 3,
		},
		{
			name:   "Start after end (auto swap), same month but earlier day",
			start:  time.Date(2023, 3, 14, 0, 0, 0, 0, time.UTC),
			end:    time.Date(2020, 3, 15, 0, 0, 0, 0, time.UTC),
			expect: 2,
		},
		{
			name:   "Leap year edge case, end day before leap day",
			start:  time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			end:    time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC),
			expect: 3,
		},
		{
			name:   "Leap year edge case, end day on leap day",
			start:  time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			end:    time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			expect: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := timex.GetYearsBetween(tt.start, tt.end)
			if got != tt.expect {
				t.Errorf("GetYearsBetween() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestGetHoursBetween(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected float64
	}{
		{
			name:     "End after Start - 2 hours",
			start:    now,
			end:      now.Add(2 * time.Hour),
			expected: 2.0,
		},
		{
			name:     "Start after End - swap and return 2 hours",
			start:    now.Add(2 * time.Hour),
			end:      now,
			expected: 2.0,
		},
		{
			name:     "Same time - zero difference",
			start:    now,
			end:      now,
			expected: 0.0,
		},
		{
			name:     "Start in past, end in future",
			start:    now.Add(-1 * time.Hour),
			end:      now.Add(1 * time.Hour),
			expected: 2.0,
		},
		{
			name:     "End is 30 minutes before start",
			start:    now,
			end:      now.Add(-30 * time.Minute),
			expected: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := timex.GetHoursBetween(tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestGetMinutesBetween(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected float64
	}{
		{
			name:     "Start before End",
			start:    time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 1, 11, 30, 0, 0, time.UTC),
			expected: 90,
		},
		{
			name:     "Start after End",
			start:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC),
			expected: 90,
		},
		{
			name:     "Equal Times",
			start:    time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "Cross-day",
			start:    time.Date(2024, 1, 1, 23, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 2, 1, 30, 30, 0, time.UTC),
			expected: 150.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := timex.GetMinutesBetween(tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("Expected %f, but got %f", tt.expected, result)
			}
		})
	}
}

func TestGetSecondsBetween(t *testing.T) {
	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected float64
	}{
		{
			name:     "End after Start",
			start:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 1, 12, 1, 30, 0, time.UTC), // 90 seconds
			expected: 90.0,
		},
		{
			name:     "Start after End",
			start:    time.Date(2024, 1, 1, 12, 1, 30, 0, time.UTC),
			end:      time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: 90.0,
		},
		{
			name:     "Equal Times",
			start:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: 0.0,
		},
		{
			name:     "With Nanoseconds",
			start:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			end:      time.Date(2024, 1, 1, 12, 0, 1, 0, time.UTC),
			expected: 1.0, // 1 second + (987654321 - 123456789)/1e9
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := timex.GetSecondsBetween(tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("GetSecondsBetween() = %v; want %v", result, tt.expected)
			}
		})
	}
}
