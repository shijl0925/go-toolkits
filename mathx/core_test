package mathx_test

import (
	"github.com/shijl0925/go-toolkits/mathx"
	"testing"
)

func Test_RoundToFloat(t *testing.T) {
	tests := []struct {
		name     string
		inputF   any // 支持 float32 和 float64
		inputN   int
		expected float64
	}{
		{"TC1: 3.14159 with n=2", 3.14159, 2, 3.14},
		{"TC2: 2.71828 with n=0", 2.71828, 0, 3.0},
		{"TC3: 1.2345 with n=-1", 1.2345, -1, 0.0},
		{"TC4: 0.499999 with n=0", 0.499999, 0, 0.0},
		{"TC5: 0.5 with n=0", 0.5, 0, 1.0},
		{"TC6: 1.005 with n=2", 1.005, 2, 1.00},
		{"TC7: float32(3.14159) with n=2", float32(3.14159), 2, 3.14},
		{"TC8: 12345.6789 with n=3", 12345.6789, 3, 12345.679},
		{"TC9: 0.0 with n=2", 0.0, 2, 0.0},
		{"TC10: -2.71828 with n=1", -2.71828, 1, -2.7},
		{"TC11: 314.159 with n=-2", 314.159, -2, 300},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch f := tt.inputF.(type) {
			case float64:
				got := mathx.RoundToFloat(f, tt.inputN)
				if got != tt.expected {
					t.Errorf("RoundToFloat(%f, %d) = %f; want %f", f, tt.inputN, got, tt.expected)
				}
			case float32:
				got := mathx.RoundToFloat(f, tt.inputN)
				if got != tt.expected {
					t.Errorf("RoundToFloat(%f, %d) = %f; want %f", f, tt.inputN, got, tt.expected)
				}
			default:
				t.Errorf("unsupported type %T", f)
			}
		})
	}
}

func Test_FloatToPercent(t *testing.T) {
	tests := []struct {
		name     string
		inputF   any // 使用 any 来支持泛型参数
		inputN   uint
		expected string
	}{
		{"TC1: Normal float64", 1.23456, 2, "123.46%"},
		{"TC2: Rounding up", 0.99999, 1, "100.0%"},
		{"TC3: Zero input", 0.0, 0, "0%"},
		{"TC4: Negative number", -0.50123, 1, "-50.1%"},
		{"TC5: Large percentage", 0.10000, 0, "10%"},
		{"TC7: High precision", 1.23456, 5, "123.45600%"},
		{"TC8: float32 input", float32(1.23456), 3, "123.456%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.inputF.(type) {
			case float64:
				got := mathx.FloatToPercent(v, tt.inputN)
				if got != tt.expected {
					t.Errorf("FloatToPercent(%f, %d) = %s; want %s", v, tt.inputN, got, tt.expected)
				}
			case float32:
				got := mathx.FloatToPercent(v, tt.inputN)
				if got != tt.expected {
					t.Errorf("FloatToPercent(%f, %d) = %s; want %s", v, tt.inputN, got, tt.expected)
				}
			default:
				t.Fatalf("Unsupported type")
			}
		})
	}
}

func TestPercentToFloat(t *testing.T) {
	type test struct {
		name    string
		input   string
		wantVal float64
		wantErr bool
	}
	tests := []test{
		{
			name:    "Empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "Missing % suffix",
			input:   "50",
			wantErr: true,
		},
		{
			name:    "Invalid number format",
			input:   "abc%",
			wantErr: true,
		},
		{
			name:    "Negative percentage",
			input:   "-1%",
			wantVal: -0.01,
			wantErr: false,
		},
		{
			name:    "Over 100 percentage",
			input:   "101%",
			wantVal: 1.01,
			wantErr: false,
		},
		{
			name:    "Zero percentage",
			input:   "0%",
			wantVal: 0,
		},
		{
			name:    "Normal case 50%",
			input:   "50%",
			wantVal: 0.5,
		},
		{
			name:    "Max valid percentage",
			input:   "100%",
			wantVal: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := mathx.PercentToFloat(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if gotVal != tt.wantVal {
					t.Errorf("expected value: %.2f, got: %.2f", tt.wantVal, gotVal)
				}
			}
		})
	}
}
