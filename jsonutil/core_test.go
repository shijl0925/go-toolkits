package jsonutil_test

import (
	"bytes"
	"fmt"
	"github.com/shijl0925/go-toolkits/jsonutil"
	"io/ioutil"
	"os"
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
		got, err := jsonutil.Encode(testUser)
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
		if got := jsonutil.Dumps(testUser); got != want {
			t.Errorf("Dumps() expected %v, got %v", want, got)
		}
	})
}

func Test_DumpsPretty(t *testing.T) {
	want := `{
    "first_name": "Alice",
    "last_name": "Smith",
    "age": 18
}`
	t.Run("test1", func(t *testing.T) {
		if got := jsonutil.DumpsPretty(testUser, 4); got != want {
			t.Errorf("DumpsPretty() expected %v, got %v", want, got)
		}
	})
}

func Test_EncodeToWriter(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		buf := &bytes.Buffer{}
		want := `{"first_name":"Alice","last_name":"Smith","age":18}
`
		err := jsonutil.EncodeToWriter(testUser, buf)
		if err != nil {
			t.Errorf("EncodeToWriter() error: %v", err)
		}
		if buf.String() != want {
			t.Errorf("EncodeToWriter() expected %v, got %v", want, buf.String())
		}
	})

	t.Run("test2", func(t *testing.T) {
		tmpfile, _ := ioutil.TempFile("", "example.json")
		err := jsonutil.EncodeToWriter(testUser, tmpfile)
		want := `{"first_name":"Alice","last_name":"Smith","age":18}
`

		if err != nil {
			t.Errorf("EncodeToWriter() error: %v", err)
		}

		buf, err := ioutil.ReadFile(tmpfile.Name())
		if string(buf) != want {
			t.Errorf("EncodeToWriter() expected %v, got %v", want, string(buf))
		}
	})
}

func Test_Decode(t *testing.T) {
	var user = &Person{}
	str := `{"first_name":"Alice","last_name":"Smith","age":18}`
	t.Run("test1", func(t *testing.T) {
		err := jsonutil.Decode([]byte(str), user)
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
		err := jsonutil.Loads(str, user)
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

		err := jsonutil.DecodeFromReader(buf, user1)
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
		err := jsonutil.DecodeFromReader(file, user2)
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
	t.Run("test", func(t *testing.T) {
		var user2 = &Person{}
		err := jsonutil.ReadFile("testdata/test.json", user2)
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

func Test_WriteFile(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		filePath := "testdata/test.json"
		err := jsonutil.WriteFile(filePath, testUser)
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
