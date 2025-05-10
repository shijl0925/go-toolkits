package mapx_test

import (
	"github.com/shijl0925/go-toolkits/mapx"
	"github.com/shijl0925/go-toolkits/slice"
	"reflect"
	"testing"
)

func TestKeys(t *testing.T) {
	testCases := []struct {
		name    string
		input   map[int]int
		wantRes []int
	}{
		{
			name:    "nil",
			input:   nil,
			wantRes: []int{},
		},
		{
			name:    "empty",
			input:   map[int]int{},
			wantRes: []int{},
		},
		{
			name: "single",
			input: map[int]int{
				1: 11,
			},
			wantRes: []int{1},
		},
		{
			name: "multiple",
			input: map[int]int{
				1: 11,
				2: 12,
			},
			wantRes: []int{1, 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := mapx.Keys[int, int](tc.input)
			res := slice.DiffSet(got, tc.wantRes)
			if len(res) > 0 {
				t.Errorf("Keys() expected %v, got %v", tc.wantRes, got)
			}
		})
	}
}

func TestValues(t *testing.T) {
	testCases := []struct {
		name    string
		input   map[int]int
		wantRes []int
	}{
		{
			name:    "nil",
			input:   nil,
			wantRes: []int{},
		},
		{
			name:    "empty",
			input:   map[int]int{},
			wantRes: []int{},
		},
		{
			name: "single",
			input: map[int]int{
				1: 11,
			},
			wantRes: []int{11},
		},
		{
			name: "multiple",
			input: map[int]int{
				1: 11,
				2: 12,
			},
			wantRes: []int{11, 12},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := mapx.Values[int, int](tc.input)
			res := slice.DiffSet(got, tc.wantRes)
			if len(res) > 0 {
				t.Errorf("Values() expected %v, got %v", tc.wantRes, got)
			}
		})
	}
}

func Test_SortMap(t *testing.T) {
	balVal := map[string]int{
		"alpha":   34,
		"beta":    12,
		"gamma":   56,
		"delta":   78,
		"epsilon": 90,
		"zeta":    13,
		"eta":     35,
		"theta":   57,
		"iota":    79,
	}
	t.Run("test1", func(t *testing.T) {
		want1 := []mapx.KV[string, int]{
			{Key: "alpha", Value: 34},
			{Key: "beta", Value: 12},
			{Key: "delta", Value: 78},
			{Key: "epsilon", Value: 90},
			{Key: "eta", Value: 35},
			{Key: "gamma", Value: 56},
			{Key: "iota", Value: 79},
			{Key: "theta", Value: 57},
			{Key: "zeta", Value: 13},
		}

		if got := mapx.SortMap(balVal); !reflect.DeepEqual(got, want1) {
			t.Errorf("SortMap() = %v, want %v", got, want1)
		}
	})

	t.Run("SortMap with custom comparator", func(t *testing.T) {
		want1 := []mapx.KV[string, int]{
			{Key: "beta", Value: 12},
			{Key: "zeta", Value: 13},
			{Key: "alpha", Value: 34},
			{Key: "eta", Value: 35},
			{Key: "gamma", Value: 56},
			{Key: "theta", Value: 57},
			{Key: "delta", Value: 78},
			{Key: "iota", Value: 79},
			{Key: "epsilon", Value: 90},
		}
		if got := mapx.SortMap(balVal, func(a, b mapx.KV[string, int]) bool {
			return a.Value < b.Value // return mapx.SortByValue(a, b)
		}); !reflect.DeepEqual(got, want1) {
			t.Errorf("SortMap() = %v, want %v", got, want1)
		}
	})
}
