package jsonx

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Encode encode data to json bytes. alias of json.Marshal
func Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Dumps encode data to json string, panic if error
// 将数据编码为JSON文本，如果发生错误则panic
//
// Example:
//
// testUser = Person{"Alice", "Smith", 18}
// res, _ := jsonx.Dumps(testUser) // res == `{"first_name":"Alice","last_name":"Smith","age":18}`
func Dumps(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// DumpsPretty encode data to json string with pretty format
// 将数据编码为带格式的JSON文本, 如果发生错误则panic
//
// Example:
//
// testUser = Person{"Alice", "Smith", 18}
// res, _ := jsonx.DumpsPretty(testUser, 4)
//
//	res == `{
//		   "first_name": "Alice",
//		   "last_name": "Smith",
//		   "age": 18
//		}`
func DumpsPretty(v any, indent int) string {
	if indent <= 0 {
		return Dumps(v)
	}
	indentStr := strings.Repeat(" ", indent)
	b, err := json.MarshalIndent(v, "", indentStr)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// EncodeToWriter encode data to json and write to writer.
// 将数据编码为JSON并写入writer
//
// Example:
//
// testUser = Person{"Alice", "Smith", 18}
// file, _ := os.OpenFile("test.json", os.O_CREATE|os.O_WRONLY, 0666)
// defer file.Close()
// err := jsonx.EncodeToWriter(testUser, file)
// if err != nil {...}
func EncodeToWriter(v any, w io.Writer) error {
	return json.NewEncoder(w).Encode(v)
}

// Decode decode json bytes to data ptr. alias of json.Unmarshal
//
// Example:
//
// var testUser Person
// data := []byte(`{"first_name":"Alice","last_name":"Smith","age":18}`)
// err := jsonx.Decode(data, &testUser)
// if err != nil {...}
// fmt.Println(testUser) // {18 Alice Smith}
func Decode(data []byte, ptr any) error {
	return json.Unmarshal(data, ptr)
}

// Loads decode json string to data ptr, panic if error
// 将JSON文本解码为数据指针，如果发生错误则恐慌
//
// Example:
//
// var testUser Person
// str := `{"first_name":"Alice","last_name":"Smith","age":18}`
// err := jsonx.Loads(str, &testUser)
// if err != nil {...}
// fmt.Println(testUser) // {18 Alice Smith}
func Loads(s string, ptr any) error {
	return json.Unmarshal([]byte(s), ptr)
}

// DecodeFromReader decode JSON from io reader.
// 从io.Reader中解码JSON文本到数据指针
//
// Example:
//
// var testUser Person
// file, _ := os.Open("test.json")
// err := jsonx.DecodeFromReader(file, &testUser)
// if err != nil {...}
// fmt.Println(testUser) // {18 Alice Smith}
func DecodeFromReader(r io.Reader, ptr any) error {
	return json.NewDecoder(r).Decode(ptr)
}

// WriteFile write data to JSON file
func WriteFile(filePath string, data any) error {
	// 使用 O_CREATE|O_WRONLY|O_TRUNC 确保文件被清空后再写入
	file, err := os.OpenFile(filepath.Clean(filePath), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	defer file.Close()
	return json.NewEncoder(file).Encode(data)
}

// ReadFile Read JSON file data
func ReadFile(filePath string, v any) error {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}

	defer file.Close()
	return json.NewDecoder(file).Decode(v)
}
