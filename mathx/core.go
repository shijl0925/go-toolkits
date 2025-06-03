package mathx

import (
	"fmt"
	"math"
	"strings"

	toolkits "github.com/shijl0925/go-toolkits"
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
	trimmedS := strings.TrimSpace(s)
	if len(trimmedS) < 2 || trimmedS[len(trimmedS)-1] != '%' {
		return 0, fmt.Errorf("invalid percent string: must be in form '<number>%%', got: %s", s)
	}

	numStr := trimmedS[:len(trimmedS)-1]
	// Further trim the number string part in case there were spaces before % like "50 %"
	numStr = strings.TrimSpace(numStr)
	if numStr == "" {
		return 0, fmt.Errorf("invalid percent string: no number found before '%%', got: %s", s)
	}

	f, err := toolkits.AnyToFloat64(numStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert percent value to float64: %v", err)
	}

	return f / 100.0, nil
}
