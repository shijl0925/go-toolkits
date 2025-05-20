package toolkits_test

import (
	"github.com/shijl0925/go-toolkits"
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
			got, err := toolkits.StringFormatTime(tt.input, tt.format...)
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
			got := toolkits.TimeFormatString(tt.input, tt.format...)
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}
