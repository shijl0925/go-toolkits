package toolkits

import (
	"fmt"
	"strings"
	"time"
)

const defaultTimeFormat = "2006-01-02 15:04:05"

// StringFormatTime converts a string to a time.Time
func StringFormatTime(s string, format ...string) (time.Time, error) {
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

// TimeFormatString converts a time.Time to a string
func TimeFormatString(t time.Time, format ...string) string {
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
