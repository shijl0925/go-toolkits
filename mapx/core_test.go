package mapx_test

import (
	"github.com/shijl0925/go-toolkits/mapx"
	"github.com/shijl0925/go-toolkits/setx"
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
			if !setx.NewFromSlice(got).Equal(setx.NewFromSlice(tc.wantRes)) {
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
			if !setx.NewFromSlice(got).Equal(setx.NewFromSlice(tc.wantRes)) {
				t.Errorf("Values() expected %v, got %v", tc.wantRes, got)
			}
		})
	}
}

// Test_HasKey_KeyExists tests the case where the key exists in the map.
func Test_HasKey_KeyExists(t *testing.T) {
	m := map[string]int{
		"one": 1,
		"two": 2,
	}
	key := "one"

	if !mapx.HasKey(m, key) {
		t.Errorf("Expected key %q to exist, but it was not found", key)
	}
}

// Test_HasKey_KeyNotExists tests the case where the key does not exist in the map.
func Test_HasKey_KeyNotExists(t *testing.T) {
	m := map[string]int{
		"one": 1,
		"two": 2,
	}
	key := "three"

	if mapx.HasKey(m, key) {
		t.Errorf("Expected key %q to not exist, but it was found", key)
	}
}

// Test_HasKey_EmptyMap tests the behavior when the map is empty.
func Test_HasKey_EmptyMap(t *testing.T) {
	m := map[int]string{}
	key := 123

	if mapx.HasKey(m, key) {
		t.Errorf("Expected key %v to not exist in an empty map", key)
	}
}

// Test_HasKey_NilMap tests the behavior when the map is nil.
func Test_HasKey_NilMap(t *testing.T) {
	var m map[string]bool
	key := "test"

	if mapx.HasKey(m, key) {
		t.Errorf("Expected key %q to not exist in a nil map", key)
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name     string
		inputs   []map[string]int // Using concrete type for testing
		expected map[string]int
	}{
		{
			name:     "Empty input list",
			inputs:   []map[string]int{},
			expected: map[string]int{},
		},
		{
			name: "Multiple non-empty maps without duplicates",
			inputs: []map[string]int{
				{"a": 1, "b": 2},
				{"c": 3, "d": 4},
			},
			expected: map[string]int{"a": 1, "b": 2, "c": 3, "d": 4},
		},
		{
			name: "Multiple maps with duplicate keys",
			inputs: []map[string]int{
				{"a": 1, "b": 2},
				{"b": 99, "c": 3},
			},
			expected: map[string]int{"a": 1, "b": 99, "c": 3},
		},
		{
			name: "Nil maps in input",
			inputs: []map[string]int{
				nil,
				{"x": 10},
				nil,
				{"y": 20},
			},
			expected: map[string]int{"x": 10, "y": 20},
		},
		{
			name: "All empty maps",
			inputs: []map[string]int{
				{},
				{},
			},
			expected: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapx.Merge(tt.inputs...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Merge() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestGetOrDefault tests the GetOrDefault function with various scenarios.
func TestGetOrDefault(t *testing.T) {
	tests := []struct {
		name       string
		inputMap   map[string]int
		key        string
		defaultVal int
		expected   int
	}{
		{
			name:       "Key exists in map",
			inputMap:   map[string]int{"a": 1, "b": 2},
			key:        "a",
			defaultVal: 0,
			expected:   1,
		},
		{
			name:       "Key does not exist in map",
			inputMap:   map[string]int{"a": 1, "b": 2},
			key:        "c",
			defaultVal: 0,
			expected:   0,
		},
		{
			name:       "Map is nil",
			inputMap:   nil,
			key:        "x",
			defaultVal: -1,
			expected:   -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapx.GetOrDefault(tt.inputMap, tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("GetOrDefault() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestGetOrDefault_GenericType verifies generic behavior with different types.
func TestGetOrDefault_GenericType(t *testing.T) {
	// Test with int keys and string values
	stringMap := map[int]string{1: "one", 2: "two"}
	if val := mapx.GetOrDefault(stringMap, 1, "default"); val != "one" {
		t.Errorf("Expected 'one', got %s", val)
	}
	if val := mapx.GetOrDefault(stringMap, 3, "default"); val != "default" {
		t.Errorf("Expected 'default', got %s", val)
	}

	// Test with struct type
	type User struct {
		Name string
	}
	userMap := map[string]User{"alice": {Name: "Alice"}}
	expectedUser := User{Name: "Default"}
	if val := mapx.GetOrDefault(userMap, "bob", expectedUser); !reflect.DeepEqual(val, expectedUser) {
		t.Errorf("Expected %+v, got %+v", expectedUser, val)
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

		if got := mapx.SortByKey(balVal, func(a, b string) bool { return a < b }); !reflect.DeepEqual(got, want1) {
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
		if got := mapx.SortByValue(balVal, func(a, b int) bool { return a < b }); !reflect.DeepEqual(got, want1) {
			t.Errorf("SortMap() = %v, want %v", got, want1)
		}
	})
}

// TestInvertNilMap tests the case when input is nil.
func TestInvertNilMap(t *testing.T) {
	var m map[string]int = nil
	got, err := mapx.InvertWithErr(m)
	if err != nil {
		t.Errorf("InvertWithErr(nil) returned error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("InvertWithErr(nil) = %v, want empty map", got)
	}
}

// TestInvertEmptyMap tests the case when input is an empty map.
func TestInvertEmptyMap(t *testing.T) {
	m := make(map[string]int)
	got, err := mapx.InvertWithErr(m)
	if err != nil {
		t.Errorf("InvertWithErr(empty) returned error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("InvertWithErr(empty) = %v, want empty map", got)
	}
}

// TestInvertUniqueValues tests inversion of a map with unique values.
func TestInvertUniqueValues(t *testing.T) {
	m := map[string]int{
		"a": 1,
		"b": 2,
	}
	want := map[int]string{
		1: "a",
		2: "b",
	}
	got, err := mapx.InvertWithErr(m)
	if err != nil {
		t.Errorf("InvertWithErr(unique values) returned error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("InvertWithErr(unique values) = %v, want %v", got, want)
	}
}

// TestInvertMultipleDuplicates tests inversion where multiple duplicates exist.
func TestInvertDuplicateValues(t *testing.T) {
	m := map[int]string{
		1: "x",
		2: "y",
		3: "x",
	}
	_, err := mapx.InvertWithErr(m)
	if err == nil {
		t.Errorf("InvertWithErr(duplicate values) returned error: %v", err)
	}
}
