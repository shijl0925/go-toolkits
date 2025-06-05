package mapx_test

import (
	"reflect"
	"testing"

	"github.com/shijl0925/go-toolkits/mapx"
	"github.com/shijl0925/go-toolkits/setx"
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
			if !setx.New[int](got).Equal(setx.New[int](tc.wantRes)) {
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
			if !setx.New[int](got).Equal(setx.New[int](tc.wantRes)) {
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

// TestIntersect_EmptyMaps tests when both maps are empty.
func TestIntersect_EmptyMaps(t *testing.T) {
	src := map[int]string{}
	dst := map[int]string{}
	expected := map[int]string{}

	result := mapx.Intersect(src, dst)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Intersect(empty maps) = %v, want %v", result, expected)
	}
}

// TestIntersect_IdenticalMaps tests when both maps are identical.
func TestIntersect_IdenticalMaps(t *testing.T) {
	src := map[string]int{"a": 1, "b": 2}
	dst := map[string]int{"a": 1, "b": 2}
	expected := map[string]int{"a": 1, "b": 2}

	result := mapx.Intersect(src, dst)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Intersect(identical maps) = %v, want %v", result, expected)
	}
}

// TestIntersect_PartialOverlap tests when only some keys match with same values.
func TestIntersect_PartialOverlap(t *testing.T) {
	src := map[string]int{"a": 1, "b": 2, "c": 3}
	dst := map[string]int{"a": 1, "b": 99, "d": 4}
	expected := map[string]int{"a": 1}

	result := mapx.Intersect(src, dst)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Intersect(partial overlap) = %v, want %v", result, expected)
	}
}

// TestIntersect_KeyExistsButDifferentValue tests when key exists but values differ.
func TestIntersect_KeyExistsButDifferentValue(t *testing.T) {
	src := map[int][]int{1: {1, 2}, 2: {3}}
	dst := map[int][]int{1: {1, 2}, 2: {4}}
	expected := map[int][]int{1: {1, 2}}

	result := mapx.Intersect(src, dst)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Intersect(key exists but different value) = %v, want %v", result, expected)
	}
}

// TestIntersect_SrcEmpty tests when source is empty.
func TestIntersect_SrcEmpty(t *testing.T) {
	src := map[int]int{}
	dst := map[int]int{1: 2}
	expected := map[int]int{}

	result := mapx.Intersect(src, dst)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Intersect(src empty) = %v, want %v", result, expected)
	}
}

// TestIntersect_DstEmpty tests when destination is empty.
func TestIntersect_DstEmpty(t *testing.T) {
	src := map[int]int{1: 2}
	dst := map[int]int{}
	expected := map[int]int{}

	result := mapx.Intersect(src, dst)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Intersect(dst empty) = %v, want %v", result, expected)
	}
}

// TestIntersect_StructValues tests with struct values.
func TestIntersect_StructValues(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	src := map[int]User{
		1: {"Alice", 30},
		2: {"Bob", 25},
	}
	dst := map[int]User{
		1: {"Alice", 30},
		2: {"Bob", 26},
	}
	expected := map[int]User{
		1: {"Alice", 30},
	}

	result := mapx.Intersect(src, dst)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Intersect(struct values) = %v, want %v", result, expected)
	}
}

// TestIntersect_PointerValues tests with pointer values pointing to same data.
func TestIntersect_PointerValues_Same(t *testing.T) {
	a := 42
	src := map[string]*int{"x": &a}
	dst := map[string]*int{"x": &a}
	expected := map[string]*int{"x": &a}

	result := mapx.Intersect(src, dst)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Intersect(pointer values) = %v, want %v", result, expected)
	}
}

// TestIntersect_PointerValues_Different tests with pointer values pointing to different data.
func TestIntersect_PointerValues_Different(t *testing.T) {
	a, b := 42, 43
	src := map[string]*int{"x": &a}
	dst := map[string]*int{"x": &b}
	expected := map[string]*int{} // Not equal because pointers point to different addresses

	result := mapx.Intersect(src, dst)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Intersect(pointer values) = %v, want %v", result, expected)
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

// TestChain_MultipleMapsNoOverlap 测试多个 map 合并且无重复 key 的情况
func TestChain_MultipleMapsNoOverlap(t *testing.T) {
	m1 := map[string]int{"a": 1}
	m2 := map[string]int{"b": 2}
	expected := map[string]int{"a": 1, "b": 2}

	result := mapx.Chain(m1, m2)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestChain_WithDuplicateKeys 测试多个 map 包含重复 key，确保只保留第一个出现的值
func TestChain_WithDuplicateKeys(t *testing.T) {
	m1 := map[int]string{1: "one"}
	m2 := map[int]string{1: "uno", 2: "two"}
	expected := map[int]string{1: "one", 2: "two"}

	result := mapx.Chain(m1, m2)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestChain_WithNilMaps 测试包含 nil map 的情况，应跳过 nil map
func TestChain_WithNilMaps(t *testing.T) {
	var m1 map[string]bool = nil
	m2 := map[string]bool{"enabled": true}
	expected := map[string]bool{"enabled": true}

	result := mapx.Chain(m1, m2)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestChain_EmptyInput 测试无输入的情况，应返回空 map
func TestChain_EmptyInput(t *testing.T) {
	expected := map[int]string{}
	result := mapx.Chain[int, string]()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestChain_SingleMap 测试单个 map 输入，应返回其副本
func TestChain_SingleMap(t *testing.T) {
	m := map[float64]int{3.14: 1, 2.71: 2}
	expected := map[float64]int{3.14: 1, 2.71: 2}

	result := mapx.Chain(m)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestFilterByKey(t *testing.T) {
	tests := []struct {
		name     string
		input    map[int]string
		filterFn func(int) bool
		expected map[int]string
	}{
		{
			name:     "EmptyMap",
			input:    map[int]string{},
			filterFn: func(_ int) bool { return true },
			expected: map[int]string{},
		},
		{
			name:     "SomeKeysMatch",
			input:    map[int]string{1: "a", 2: "b"},
			filterFn: func(k int) bool { return k > 1 },
			expected: map[int]string{2: "b"},
		},
		{
			name:     "EvenKeysOnly",
			input:    map[int]string{1: "a", 2: "b", 3: "c"},
			filterFn: func(k int) bool { return k%2 == 0 },
			expected: map[int]string{2: "b"},
		},
		{
			name:     "AllKeysMatch",
			input:    map[int]string{1: "a", 2: "b"},
			filterFn: func(k int) bool { return k != 0 },
			expected: map[int]string{1: "a", 2: "b"},
		},
		{
			name:     "NoKeysMatch",
			input:    map[int]string{1: "a", 2: "b"},
			filterFn: func(k int) bool { return k == 0 },
			expected: map[int]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapx.FilterByKey(tt.input, tt.filterFn)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("FilterByKey() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// 可选：测试非 int 类型的 key，例如 string
func TestFilterByKey_GenericKeyType(t *testing.T) {
	input := map[string]int{"a": 1, "bb": 2, "c": 3}
	fn := func(k string) bool {
		return len(k) == 1
	}
	expected := map[string]int{"a": 1, "c": 3}

	got := mapx.FilterByKey(input, fn)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("FilterByKey() with string keys = %v, want %v", got, expected)
	}
}

// TestFilterByValue 是 FilterByValue 的单元测试
func TestFilterByValue(t *testing.T) {
	tests := []struct {
		name     string
		inputMap map[any]any
		filterFn func(any) bool
		expected map[any]any
	}{
		{
			name:     "Empty Map",
			inputMap: map[any]any{},
			filterFn: func(_ any) bool { return true },
			expected: map[any]any{},
		},
		{
			name:     "Nil Map",
			inputMap: nil,
			filterFn: func(_ any) bool { return true },
			expected: map[any]any{},
		},
		{
			name:     "All Elements Match",
			inputMap: map[any]any{1: "apple", 2: "banana"},
			filterFn: func(_ any) bool { return true },
			expected: map[any]any{1: "apple", 2: "banana"},
		},
		{
			name:     "No Elements Match",
			inputMap: map[any]any{1: "apple", 2: "banana"},
			filterFn: func(_ any) bool { return false },
			expected: map[any]any{},
		},
		{
			name:     "Some Elements Match - String Value",
			inputMap: map[any]any{1: "apple", 2: "banana", 3: "cherry"},
			filterFn: func(v any) bool { return v.(string) == "banana" },
			expected: map[any]any{2: "banana"},
		},
		{
			name:     "Some Elements Match - Int Value",
			inputMap: map[any]any{"a": 10, "b": 20, "c": 30},
			filterFn: func(v any) bool { return v.(int) > 15 },
			expected: map[any]any{"b": 20, "c": 30},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapx.FilterByValue(tt.inputMap, tt.filterFn)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("FilterByValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// 测试用例
func TestFilter(t *testing.T) {
	tests := []struct {
		name     string
		input    map[int]string
		filterFn func(int, string) bool
		expected map[int]string
	}{
		{
			name:     "NilMap_ReturnsEmpty",
			input:    nil,
			filterFn: func(_ int, _ string) bool { return true },
			expected: map[int]string{},
		},
		{
			name:     "EmptyMap_ReturnsEmpty",
			input:    map[int]string{},
			filterFn: func(_ int, _ string) bool { return true },
			expected: map[int]string{},
		},
		{
			name:     "AllElementsMatch_ReturnsSame",
			input:    map[int]string{1: "a", 2: "b"},
			filterFn: func(_ int, _ string) bool { return true },
			expected: map[int]string{1: "a", 2: "b"},
		},
		{
			name:     "NoElementsMatch_ReturnsEmpty",
			input:    map[int]string{1: "a", 2: "b"},
			filterFn: func(_ int, _ string) bool { return false },
			expected: map[int]string{},
		},
		{
			name:     "FilterByKeyEven",
			input:    map[int]string{1: "a", 2: "b", 3: "c", 4: "dd", 5: "e"},
			filterFn: func(k int, v string) bool { return k%2 == 0 && len(v) == 1 },
			expected: map[int]string{2: "b"},
		},
		{
			name:     "FilterByValueLengthGreaterThanOne",
			input:    map[int]string{1: "a", 2: "ab", 3: "abc"},
			filterFn: func(_ int, v string) bool { return len(v) > 1 },
			expected: map[int]string{2: "ab", 3: "abc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapx.Filter(tt.input, tt.filterFn)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Filter() = %v, want %v", got, tt.expected)
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

func TestSetIfAbsent(t *testing.T) {
	tests := []struct {
		name       string
		inputMap   map[string]int
		key        string
		defaultVal int
		expected   map[string]int
	}{
		{
			name:       "EmptyMap_KeyNotPresent",
			inputMap:   make(map[string]int),
			key:        "a",
			defaultVal: 1,
			expected: map[string]int{
				"a": 1,
			},
		},
		{
			name: "NonEmptyMap_KeyNotPresent",
			inputMap: map[string]int{
				"b": 2,
			},
			key:        "a",
			defaultVal: 1,
			expected: map[string]int{
				"a": 1,
				"b": 2,
			},
		},
		{
			name: "NonEmptyMap_KeyPresent",
			inputMap: map[string]int{
				"a": 99,
			},
			key:        "a",
			defaultVal: 1,
			expected: map[string]int{
				"a": 99,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapx.SetIfAbsent(tt.inputMap, tt.key, tt.defaultVal)
			if !reflect.DeepEqual(tt.inputMap, tt.expected) {
				t.Errorf("SetIfAbsent() got %v, want %v", tt.inputMap, tt.expected)
			}
		})
	}
}

func TestSortByKey_EmptyMap(t *testing.T) {
	m := map[int]string{}
	result := mapx.SortByKey(m, func(a, b int) bool { return a < b })
	if len(result) != 0 {
		t.Errorf("Expected an empty slice, got %v", result)
	}
}

func TestSortByKey_IntKeysAscending(t *testing.T) {
	m := map[int]string{
		3: "three",
		1: "one",
		2: "two",
	}
	expected := []mapx.KV[int, string]{
		{Key: 1, Value: "one"},
		{Key: 2, Value: "two"},
		{Key: 3, Value: "three"},
	}
	result := mapx.SortByKey(m, func(a, b int) bool { return a < b })
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("SortByKey() expected %v, got %v", expected, result)
	}
}

func TestSortByKey_StringKeysDescending(t *testing.T) {
	m := map[string]int{
		"apple":  1,
		"banana": 2,
		"cherry": 3,
	}
	expected := []mapx.KV[string, int]{
		{"cherry", 3},
		{"banana", 2},
		{"apple", 1},
	}
	result := mapx.SortByKey(m, func(a, b string) bool { return a > b })
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("SortByKey() expected %v, got %v", expected, result)
	}
}

func TestSortByKey_CustomType(t *testing.T) {
	type CustomKey struct {
		ID   int
		Name string
	}
	m := map[CustomKey]string{
		{ID: 2, Name: "B"}: "ValueB",
		{ID: 1, Name: "A"}: "ValueA",
		{ID: 3, Name: "C"}: "ValueC",
	}
	expected := []mapx.KV[CustomKey, string]{
		{Key: CustomKey{ID: 1, Name: "A"}, Value: "ValueA"},
		{Key: CustomKey{ID: 2, Name: "B"}, Value: "ValueB"},
		{Key: CustomKey{ID: 3, Name: "C"}, Value: "ValueC"},
	}
	result := mapx.SortByKey(m, func(a, b CustomKey) bool { return a.ID < b.ID })
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("SortByKey() expected %v, got %v", expected, result)
	}
}

// 测试用例 1：空 map
func Test_SortByValue_EmptyMap(t *testing.T) {
	input := map[string]int{}
	result := mapx.SortByValue(input, func(a, b int) bool { return a < b })

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %v", result)
	}
}

// 测试用例 2：单个元素
func Test_SortByValue_SingleElement(t *testing.T) {
	input := map[string]int{"a": 1}
	expected := []mapx.KV[string, int]{{Key: "a", Value: 1}}

	result := mapx.SortByValue(input, func(a, b int) bool { return a < b })

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("SortByValue() expected %v, got %v", expected, result)
	}
}

// 测试用例 3：多个整数，升序排序
func Test_SortByValue_MultipleInts_Asc(t *testing.T) {
	input := map[string]int{
		"a": 3,
		"b": 1,
		"c": 2,
	}
	expected := []mapx.KV[string, int]{
		{"b", 1},
		{"c", 2},
		{"a", 3},
	}

	result := mapx.SortByValue(input, func(a, b int) bool { return a < b })
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("SortByValue() expected %v, got %v", expected, result)
	}
}

// 测试用例 4：多个整数，降序排序
func Test_SortByValue_MultipleInts_Desc(t *testing.T) {
	input := map[string]int{
		"a": 3,
		"b": 1,
		"c": 2,
	}
	expected := []mapx.KV[string, int]{
		{"a", 3},
		{"c", 2},
		{"b", 1},
	}

	result := mapx.SortByValue(input, func(a, b int) bool { return a > b })
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("SortByValue() expected %v, got %v", expected, result)
	}
}

// 测试用例 5：字符串值，字典序排序
func Test_SortByValue_Strings(t *testing.T) {
	input := map[int]string{
		1: "banana",
		2: "apple",
		3: "cherry",
	}
	expected := []mapx.KV[int, string]{
		{2, "apple"},
		{1, "banana"},
		{3, "cherry"},
	}

	expectedKeys := []int{2, 1, 3} // apple < banana < cherry

	result := mapx.SortByValue(input, func(a, b string) bool { return a < b })

	for i, kv := range result {
		if kv.Key != expectedKeys[i] {
			t.Errorf("Expected key %d at index %d, got %d", expectedKeys[i], i, kv.Key)
		}
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("SortByValue() expected %v, got %v", expected, result)
	}
}

// 测试用例 6：自定义结构体
type Person struct {
	Name string
	Age  int
}

func Test_SortByValue_CustomStruct(t *testing.T) {
	input := map[string]Person{
		"a": {"Alice", 30},
		"b": {"Bob", 25},
		"c": {"Charlie", 35},
	}
	expected := []mapx.KV[string, Person]{
		{"b", Person{"Bob", 25}},
		{"a", Person{"Alice", 30}},
		{"c", Person{"Charlie", 35}},
	}

	expectedAges := []int{25, 30, 35}

	result := mapx.SortByValue(input, func(a, b Person) bool { return a.Age < b.Age })

	for i, age := range expectedAges {
		if result[i].Value.Age != age {
			t.Errorf("Expected age %d at index %d, got %d", age, i, result[i].Value.Age)
		}
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("SortByValue() expected %v, got %v", expected, result)
	}
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

// 定义测试结构体
type TestStruct struct {
	Name    string   `json:"name"`
	Age     int      `json:"-"`
	Tags    []string `json:"tags"`
	Score   float64  `json:"score"`
	Active  bool     `json:"active"`
	Weight  *int
	Address Address `json:"address"`
}

type Address struct {
	City string `json:"city"`
}

func TestMapToStruct(t *testing.T) {
	// 测试用例 1: 成功转换简单结构体 (TC001)
	t.Run("SuccessfulConversion", func(t *testing.T) {
		m := map[string]any{
			"name":   "Alice",
			"age":    30,
			"score":  85.0,
			"active": true,
			"tags":   []string{"golang", "python"},
			"address": map[string]interface{}{
				"city": "New York",
			},
		}

		var s TestStruct
		if err := mapx.MapToStruct(m, &s); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		//fmt.Printf("s: %+v\n", s)

		if s.Name != "Alice" || s.Age != 0 || s.Address.City != "New York" || s.Score != 85.0 || s.Active != true || !reflect.DeepEqual(s.Tags, []string{"golang", "python"}) {
			t.Errorf("Unexpected values: %+v", s)
		}
	})

	// 测试用例 2: 忽略带 "-" tag 的字段 (TC002)
	t.Run("IgnoreFieldWithTagMinus", func(t *testing.T) {
		m := map[string]any{
			"age": 30, // age 字段被标记为 "-"
		}

		var s TestStruct
		if err := mapx.MapToStruct(m, &s); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if s.Age != 0 {
			t.Errorf("Age should be ignored but got %d", s.Age)
		}
	})

	// 测试用例 3: 处理未导出字段 (TC003)
	t.Run("UnexportedFields", func(t *testing.T) {
		// TestStruct 中没有未导出字段，创建一个新结构体
		type S struct {
			unexported string
		}

		m := map[string]any{
			"unexported": "value",
		}

		var s S
		if err := mapx.MapToStruct(m, &s); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if s.unexported != "" {
			t.Errorf("Unexported field should not be set")
		}
	})

	// 测试用例 4: 非结构体指针输入 (TC009)
	t.Run("NonStructPointerInput", func(t *testing.T) {
		m := map[string]any{}
		var i int

		if err := mapx.MapToStruct(m, &i); err == nil {
			t.Error("Expected error for non-struct pointer input, got nil")
		} else if err.Error() != "must be a struct pointer" {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	// 测试用例 5: nil 指针输入 (TC010)
	t.Run("NilPointerInput", func(t *testing.T) {
		m := map[string]any{}
		var s *TestStruct

		if err := mapx.MapToStruct(m, s); err == nil {
			t.Error("Expected error for nil pointer input, got nil")
		} else if err.Error() != "must be a non empty struct pointer" {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	// 测试用例 6: 处理指针字段 - 零值 (TC004)
	t.Run("ZeroValuePointerField", func(t *testing.T) {
		m := map[string]any{}

		var s TestStruct
		s.Weight = new(int) // 初始化为非 nil(指向0的指针)
		if err := mapx.MapToStruct(m, &s); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		//fmt.Printf("s: %+v\n", s)

		if *s.Weight != 0 {
			t.Errorf("Weight should be nil")
		}
	})

	// 测试用例 7: 处理指针字段 - 非零值 (TC005)
	t.Run("NonZeroValuePointerField", func(t *testing.T) {
		m := map[string]any{
			"weight": 75,
		}

		var s TestStruct
		if err := mapx.MapToStruct(m, &s); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		//fmt.Printf("s: %+v\n", s)
		if s.Weight == nil || *s.Weight != 75 {
			t.Errorf("Weight should be 75 and not nil")
		}
	})

	// 测试用例 8: 嵌套结构体转换 (TC006)
	t.Run("NestedStructConversion", func(t *testing.T) {
		m := map[string]any{
			"address": map[string]any{
				"city": "Beijing",
			},
		}

		var s TestStruct
		if err := mapx.MapToStruct(m, &s); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		//fmt.Printf("s: %+v\n", s)
		if s.Address.City != "Beijing" {
			t.Errorf("Address.City should be Beijing")
		}
	})

	// 测试用例 9: 不支持的类型转换 (TC007)
	t.Run("UnsupportedTypeConversion", func(t *testing.T) {
		type S struct {
			Field int
		}

		m := map[string]any{
			"field": "string_value",
		}

		var s S
		if err := mapx.MapToStruct(m, &s); err == nil {
			t.Error("Expected error for unsupported type conversion, got nil")
		}
	})
}
