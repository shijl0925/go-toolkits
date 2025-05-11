package structs_test

import (
	"github.com/shijl0925/go-toolkits/structs"
	"reflect"
	"testing"
	"time"
)

// TestStructToMap_BasicFields 测试普通结构体字段的转换
func TestStructToMap_BasicFields(t *testing.T) {
	type SampleStruct struct {
		Name string
		Age  int
	}

	input := SampleStruct{
		Name: "Alice",
		Age:  30,
	}

	expected := map[string]any{
		"Name": "Alice",
		"Age":  30,
	}

	result := structs.StructToMap(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestStructToMap_PointerStructFields(t *testing.T) {
	type SampleStruct struct {
		Name string
		Age  int
	}

	input := &SampleStruct{
		Name: "Alice",
		Age:  30,
	}

	expected := map[string]any{
		"Name": "Alice",
		"Age":  30,
	}

	result := structs.StructToMap(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestStructToMap_WithTag(t *testing.T) {
	type SampleStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	input := SampleStruct{
		Name: "Alice",
		Age:  30,
	}

	expected := map[string]any{
		"name": "Alice",
		"age":  30,
	}

	result := structs.StructToMap(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestStructToMap_EmptyStruct 测试空结构体
func TestStructToMap_EmptyStruct(t *testing.T) {
	type EmptyStruct struct{}

	input := EmptyStruct{}

	result := structs.StructToMap(input) // 输出: map[]

	if len(result) != 0 {
		t.Errorf("Expected empty map, got %v", result)
	}
}

// TestStructToMap_NestedStruct 测试嵌套结构体
func TestStructToMap_NestedStruct(t *testing.T) {
	type InnerStruct struct {
		X int
		Y int
	}

	type OuterStruct struct {
		ID   int
		Data InnerStruct `json:"data,flatten"`
	}

	input := OuterStruct{
		ID: 1,
		Data: InnerStruct{
			X: 10,
			Y: 20,
		},
	}

	expected := map[string]any{
		"ID": 1,
		"data": map[string]any{
			"X": 10,
			"Y": 20,
		},
	}

	result := structs.StructToMap(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestStructToMap_PointerFields 测试包含指针字段的结构体
func TestStructToMap_PointerFields(t *testing.T) {
	type StructWithPointer struct {
		Name *string
		Age  *int
	}

	name := "Bob"
	age := 25

	input := StructWithPointer{
		Name: &name,
		Age:  &age,
	}

	expected := map[string]any{
		"Name": "Bob",
		"Age":  25,
	}

	result := structs.StructToMap(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestStructToMap_ComplexTypes 测试包含复杂类型的结构体（slice, map）
func TestStructToMap_ComplexTypes(t *testing.T) {
	type ComplexStruct struct {
		Tags    []string
		Config  map[string]int
		Scores  [2]int
		Enabled bool
	}

	input := ComplexStruct{
		Tags: []string{"go", "test"},
		Config: map[string]int{
			"a": 1,
			"b": 2,
		},
		Scores:  [2]int{90, 85},
		Enabled: true,
	}

	expected := map[string]any{
		"Tags": []string{"go", "test"},
		"Config": map[string]int{
			"a": 1,
			"b": 2,
		},
		"Scores":  [2]int{90, 85},
		"Enabled": true,
	}

	result := structs.StructToMap(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestStructToMap_omitNested(t *testing.T) {
	type Server struct {
		Name string    `json:"name"`
		ID   int       `json:"id"`
		Time time.Time `json:"time,omitnested"` // do not convert to map[string]interface{}
	}

	input := Server{
		Name: "Zeynep",
		ID:   789012,
		Time: time.Now(),
	}

	result := structs.StructToMap(input)

	expected := map[string]any{
		"name": "Zeynep",
		"id":   789012,
		"time": result["time"],
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestStructToMap_OmitEmpty(t *testing.T) {
	type S struct {
		Name string `json:",omitempty"`
		Age  int    `json:",omitempty"`
	}
	input := S{}
	expected := map[string]any{}

	result := structs.StructToMap(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestStructToMap_omitEmptyWithValue(t *testing.T) {
	// By default field with struct types of zero values are processed too. We
	// can stop processing them via "omitempty" tag option.
	type Server struct {
		Name     string `json:"name,omitempty"`
		ID       int32  `json:"id,omitempty"`
		Location string
	}

	// Only add location
	input := Server{
		Location: "Tokyo",
	}

	expected := map[string]any{
		"Location": "Tokyo",
	}

	result := structs.StructToMap(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestStructToMap_MultiLevelNested(t *testing.T) {
	type Address struct {
		City string `json:",flatten"`
	}
	type Profile struct {
		Address Address `json:",flatten"`
	}
	type User struct {
		Name    string
		Profile Profile `json:",flatten"`
	}
	input := User{
		Name: "Frank",
		Profile: Profile{
			Address: Address{City: "Shanghai"},
		},
	}

	expected := map[string]any{
		"Name": "Frank",
		"Profile": map[string]any{
			"Address": map[string]any{
				"City": "Shanghai",
			},
		},
	}

	result := structs.StructToMap(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestStructToMap_AnonymousStruct(t *testing.T) {
	type Base struct {
		ID int
	}
	type Derived struct {
		Base `json:",flatten"`
		Name string
	}
	input := Derived{
		Base: Base{ID: 1},
		Name: "Grace",
	}
	expected := map[string]any{
		"ID":   1,
		"Name": "Grace",
	}

	result := structs.StructToMap(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
