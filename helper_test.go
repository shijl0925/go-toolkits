package toolkits_test

import (
	toolkits "github.com/shijl0925/go-toolkits"
	"math"
	"testing"
)

func TestEqualFloat64(t *testing.T) {
	tests := []struct {
		name     string
		a        float64
		b        float64
		decimal  int
		expected bool
		panicMsg string
	}{
		{
			name:     "equal within 2 decimals",
			a:        0.1 + 0.2,
			b:        0.3,
			decimal:  2,
			expected: true,
		},
		{
			name:     "equal within 5 decimals",
			a:        1.00001,
			b:        1.000012,
			decimal:  5,
			expected: true,
		},
		{
			name:     "not equal within 5 decimals",
			a:        1.00001,
			b:        1.00002,
			decimal:  5,
			expected: true,
		},
		{
			name:     "equal at decimal 0",
			a:        1.9,
			b:        2.1,
			decimal:  0,
			expected: true,
		},
		{
			name:     "not equal at decimal 0",
			a:        1.4,
			b:        2.6,
			decimal:  0,
			expected: true,
		},
		{
			name: "decimal is 15 (max allowed)",

			a:        1.000000000000001,
			b:        1.000000000000002,
			decimal:  15,
			expected: true,
		},
		{
			name:     "decimal less than 0",
			a:        1.0,
			b:        1.0,
			decimal:  -1,
			panicMsg: "decimal must be between 0 and 15 inclusive",
		},
		{
			name:     "decimal greater than 15",
			a:        1.0,
			b:        1.0,
			decimal:  16,
			panicMsg: "decimal must be between 0 and 15 inclusive",
		},
		{
			name:     "both are NaN",
			a:        math.NaN(),
			b:        math.NaN(),
			decimal:  2,
			expected: false,
		},
		{
			name:     "one is NaN",
			a:        math.NaN(),
			b:        1.0,
			decimal:  2,
			expected: false,
		},
		{
			name:     "both are Inf same sign",
			a:        math.Inf(1),
			b:        math.Inf(1),
			decimal:  2,
			expected: true,
		},
		{
			name:     "both are Inf different sign",
			a:        math.Inf(1),
			b:        math.Inf(-1),
			decimal:  2,
			expected: false,
		},
		{
			name:     "one is Inf",
			a:        math.Inf(1),
			b:        1.0,
			decimal:  2,
			expected: false,
		},
		{
			name:     "zero comparison",
			a:        0.0,
			b:        -0.0,
			decimal:  10,
			expected: true,
		},
		{
			name:     "very small numbers",
			a:        1e-16,
			b:        2e-16,
			decimal:  15,
			expected: true,
		},
		{
			name:     "large numbers with small relative difference",
			a:        1e15,
			b:        1e15 + 1e10,
			decimal:  5,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tt.panicMsg == "" {
						t.Errorf("unexpected panic: %v", r)
					} else if r != tt.panicMsg {
						t.Errorf("expected panic message: %q, got: %q", tt.panicMsg, r)
					}
				}
			}()

			result := toolkits.EqualFloat64(tt.a, tt.b, tt.decimal)
			if result != tt.expected {
				t.Errorf("EqualFloat64(%v, %v, %d) = %v; want %v", tt.a, tt.b, tt.decimal, result, tt.expected)
			}
		})
	}
}
