package mapx_test

import (
	"github.com/shijl0925/go-toolkits/mapx"
	"github.com/shijl0925/go-toolkits/slice"
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
