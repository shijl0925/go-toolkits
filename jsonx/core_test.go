package jsonx_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shijl0925/go-toolkits/jsonx"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

var testUser = Person{"Alice", "Smith", 18}

func Test_Encode(t *testing.T) {
	want := `{"first_name":"Alice","last_name":"Smith","age":18}`
	t.Run("test1", func(t *testing.T) {
		got, err := jsonx.Encode(testUser)
		if err != nil {
			t.Errorf("Encode() error = %v", err)
		}
		if string(got) != want {
			t.Errorf("Encode() expected %v, got %v", want, got)
		}
	})
}

func Test_Dumps(t *testing.T) {
	want := `{"first_name":"Alice","last_name":"Smith","age":18}`
	t.Run("test1", func(t *testing.T) {
		got, err := jsonx.Dumps(testUser)
		if err != nil {
			t.Errorf("Dumps() error = %v", err)
		}
		if got != want {
			t.Errorf("Dumps() expected %v, got %v", want, got)
		}
	})
}

/*func TestDumps_InvalidInput(t *testing.T) {
	type hasChan struct {
		C chan int
	}

	t.Run("Invalid type (chan)", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to invalid type, but did not get one")
			}
		}()

		jsonx.Dumps(hasChan{C: make(chan int)})
	})
}*/

func Test_DumpsPretty(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		indent   int
		expected string
	}{
		{
			name:     "Positive Indent - Simple Map",
			input:    map[string]int{"a": 1, "b": 2},
			indent:   2,
			expected: "{\n  \"a\": 1,\n  \"b\": 2\n}",
		},
		{
			name:     "Zero Indent - Uses Dumps",
			input:    []int{1, 2, 3},
			indent:   0,
			expected: "[1,2,3]",
		},
		{
			name:     "Negative Indent - Uses Dumps",
			input:    struct{}{},
			indent:   -1,
			expected: "{}",
		},
		{
			name:     "Nil Input",
			input:    nil,
			indent:   4,
			expected: "null",
		},
		{
			name:     "Nested Structure",
			input:    map[string]any{"a": 1, "b": map[string]int{"c": 3}},
			indent:   3,
			expected: "{\n   \"a\": 1,\n   \"b\": {\n      \"c\": 3\n   }\n}",
		},
		{
			name:     "Nested Structure with Indent",
			input:    testUser,
			indent:   4,
			expected: "{\n    \"first_name\": \"Alice\",\n    \"last_name\": \"Smith\",\n    \"age\": 18\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonx.DumpsPretty(tt.input, tt.indent)
			if err != nil {
				t.Errorf("DumpsPretty() error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("DumpsPretty() = %s, want %s", got, tt.expected)
			}
		})
	}
}

/*func TestDumpsPretty_InvalidInput(t *testing.T) {
	type hasChan struct {
		C chan int
	}

	t.Run("Invalid type (chan)", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic due to invalid type, but did not get one")
			}
		}()

		jsonx.DumpsPretty(hasChan{C: make(chan int)}, 2)
	})
}*/

func Test_EncodeToWriter(t *testing.T) {
	t.Run("TC01 - Valid JSON file", func(t *testing.T) {
		buf := &bytes.Buffer{}
		want := `{"first_name":"Alice","last_name":"Smith","age":18}
`
		err := jsonx.EncodeToWriter(testUser, buf)
		if err != nil {
			t.Errorf("EncodeToWriter() error: %v", err)
		}
		if buf.String() != want {
			t.Errorf("EncodeToWriter() expected %v, got %v", want, buf.String())
		}
	})

	t.Run("TC02 - Valid JSON file", func(t *testing.T) {
		tmpfile, _ := ioutil.TempFile("", "example.json")
		err := jsonx.EncodeToWriter(testUser, tmpfile)
		want := `{"first_name":"Alice","last_name":"Smith","age":18}
`

		if err != nil {
			t.Errorf("EncodeToWriter() error: %v", err)
		}

		buf, err := ioutil.ReadFile(tmpfile.Name())
		if err != nil {
			t.Errorf("ReadFile() error: %v", err)
		}
		if string(buf) != want {
			t.Errorf("EncodeToWriter() expected %v, got %v", want, string(buf))
		}
	})
}

func Test_Decode(t *testing.T) {
	var user = &Person{}
	str := `{"first_name":"Alice","last_name":"Smith","age":18}`
	t.Run("test1", func(t *testing.T) {
		err := jsonx.Decode([]byte(str), user)
		if err != nil {
			t.Errorf("Decode() error = %v", err)
		}
		if user.FirstName != "Alice" {
			t.Errorf("Decode() expected %v, got %v", "Alice", user.FirstName)
		}
		if user.LastName != "Smith" {
			t.Errorf("Decode() expected %v, got %v", "Smith", user.LastName)
		}
		if user.Age != 18 {
			t.Errorf("Decode() expected %v, got %v", 18, user.Age)
		}
	})
}

func Test_Loads(t *testing.T) {
	var user = &Person{}
	str := `{"first_name":"Alice","last_name":"Smith","age":18}`
	t.Run("test1", func(t *testing.T) {
		err := jsonx.Loads(str, user)
		if err != nil {
			t.Errorf("Loads() error = %v", err)
		}
		if user.FirstName != "Alice" {
			t.Errorf("Loads() expected %v, got %v", "Alice", user.FirstName)
		}
		if user.LastName != "Smith" {
			t.Errorf("Loads() expected %v, got %v", "Smith", user.LastName)
		}
		if user.Age != 18 {
			t.Errorf("Loads() expected %v, got %v", 18, user.Age)
		}
	})
}

func Test_DecodeFromReader(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		var user1 = &Person{}

		want := `{"first_name":"Alice","last_name":"Smith","age":18}`

		buf := &bytes.Buffer{}
		buf.WriteString(want)
		fmt.Printf("buf: %s\n", buf.String())

		err := jsonx.DecodeFromReader(buf, user1)
		if err != nil {
			t.Errorf("DecodeFromReader() error: %v", err)
		}
		if user1.FirstName != "Alice" {
			t.Errorf("DecodeFromReader() expected %v, got %v", "Alice", user1.FirstName)
		}
		if user1.LastName != "Smith" {
			t.Errorf("DecodeFromReader() expected %v, got %v", "Smith", user1.LastName)
		}
		if user1.Age != 18 {
			t.Errorf("DecodeFromReader() expected %v, got %v", 18, user1.Age)
		}
	})

	t.Run("test2", func(t *testing.T) {
		var user2 = &Person{}
		file, _ := os.Open("testdata/test.json")
		err := jsonx.DecodeFromReader(file, user2)
		if err != nil {
			t.Errorf("DecodeFromReader() error: %v", err)
		}
		if user2.FirstName != "Alice" {
			t.Errorf("DecodeFromReader() expected %v, got %v", "Alice", user2.FirstName)
		}
		if user2.LastName != "Smith" {
			t.Errorf("DecodeFromReader() expected %v, got %v", "Smith", user2.LastName)
		}
		if user2.Age != 18 {
			t.Errorf("DecodeFromReader() expected %v, got %v", 18, user2.Age)
		}
	})
}

func Test_ReadFile(t *testing.T) {
	t.Run("Valid JSON file", func(t *testing.T) {
		var user2 = &Person{}
		err := jsonx.ReadFile("testdata/test.json", user2)
		if err != nil {
			t.Errorf("ReadFile() error: %v", err)
		}
		if user2.FirstName != "Alice" {
			t.Errorf("ReadFile() expected %v, got %v", "Alice", user2.FirstName)
		}
		if user2.LastName != "Smith" {
			t.Errorf("ReadFile() expected %v, got %v", "Smith", user2.LastName)
		}
		if user2.Age != 18 {
			t.Errorf("ReadFile() expected %v, got %v", 18, user2.Age)
		}
	})
}
func Test_ReadFile_InvalidFile(t *testing.T) {
	t.Run("Invalid File", func(t *testing.T) {
		var data Person
		err := jsonx.ReadFile("nonexistent.json", &data)
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
	})

	t.Run("Invalid JSON content", func(t *testing.T) {
		// 创建临时文件并写入非法 JSON 内容
		tmpfile, err := os.CreateTemp("", "test_invalid_*.json")
		if err != nil {
			t.Errorf("Error creating temporary file: %v", err)
		}
		defer os.Remove(tmpfile.Name())

		_, err = tmpfile.WriteString("invalid json content")
		if err != nil {
			t.Errorf("Error writing to temporary file: %v", err)
		}
		tmpfile.Close()

		var data Person
		err = jsonx.ReadFile(tmpfile.Name(), &data)
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
	})
}

func Test_WriteFile(t *testing.T) {
	t.Run("TC01 - Valid JSON file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test.json")
		data := map[string]string{"key": "value"}

		err := jsonx.WriteFile(filePath, data)
		if err != nil {
			t.Errorf("WriteFile() error: %v", err)
		}

		// 验证文件内容
		content, _ := os.ReadFile(filePath)
		var result map[string]string
		err = json.Unmarshal(content, &result)
		if err != nil {
			t.Errorf("Error unmarshalling JSON: %v", err)
		}
		if !reflect.DeepEqual(result, data) {
			t.Errorf("WriteFile() expected %v, got %v", data, result)
		}
	})
	t.Run("TC02 - Valid JSON file", func(t *testing.T) {
		filePath := "testdata/test.json"
		err := jsonx.WriteFile(filePath, testUser)
		want := `{"first_name":"Alice","last_name":"Smith","age":18}
`
		if err != nil {
			t.Errorf("WriteFile() error: %v", err)
		}

		data, _ := ioutil.ReadFile(filePath)
		if string(data) != want {
			t.Errorf("WriteFile() expected %v, got %v", want, string(data))
		}
	})
}

// 测试用例：无效路径导致 OpenFile 失败
func TestWriteFile_OpenFileError(t *testing.T) {
	invalidPath := "/invalid/path/test.json"
	err := jsonx.WriteFile(invalidPath, map[string]string{"key": "value"})
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

// 测试用例：不可序列化数据导致 Encode 失败
func TestWriteFile_EncodeError(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.json")

	// channel 类型不能被 JSON 序列化
	err := jsonx.WriteFile(filePath, make(chan int))
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}
