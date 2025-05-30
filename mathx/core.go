package mathx

import (
	"fmt"
	toolkits "github.com/shijl0925/go-toolkits"
	"math"
)

// RoundToFloat returns the float32/float64 of the specified precision from rounding
func RoundToFloat[T float64 | float32](f T, n int) float64 {
	pow10N := math.Pow(10.0, float64(n))
	val := float64(f)
	return math.Round(val*pow10N) / pow10N
}

// FloatToPercent returns the percent string of the specified precision
func FloatToPercent[T float64 | float32](f T, n uint) string {
	return fmt.Sprintf("%.*f%%", n, RoundToFloat(f*100.0, int(n)))
}

// PercentToFloat returns the float64 of the percent string
func PercentToFloat(s string) (float64, error) {
	if len(s) < 2 || s[len(s)-1] != '%' {
		return 0, fmt.Errorf("invalid percent string: must be in form '<number>%%', got: %s", s)
	}

	s = s[:len(s)-1]
	f, err := toolkits.AnyToFloat64(s)
	if err != nil {
		return 0, fmt.Errorf("failed to convert percent value to float64: %v", err)
	}

	return f / 100.0, nil
}
