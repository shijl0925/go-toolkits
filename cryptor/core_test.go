package cryptor_test

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"github.com/shijl0925/go-toolkits/cryptor"
	"os"
	"testing"
)

func TestBase64StdEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "TC01 - 常规英文字符串",
			input:    "hello",
			expected: "aGVsbG8=",
		},
		{
			name:     "TC02 - 空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "TC03 - 数字字符串",
			input:    "123456",
			expected: "MTIzNDU2",
		},
		{
			name:     "TC04 - 中文及特殊字符",
			input:    "特殊字符!@#",
			expected: "54m55q6K5a2X56ymIUAj",
		},
		{
			name:     "TC05 - 包含等号的字符串",
			input:    "abc==",
			expected: "YWJjPT0=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cryptor.Base64StdEncode(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Test cases for Base64StdDecode function
func TestBase64StdDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Valid Base64 String",
			input:    "SGVsbG8gd29ybGQh", // "Hello world!"
			expected: "Hello world!",
			wantErr:  false,
		},
		{
			name:     "Invalid Base64 String",
			input:    "invalid_base64",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Empty String",
			input:    "",
			expected: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cryptor.Base64StdDecode(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result)
				}
			}
		})
	}
}

func TestMd5String(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "hello",
			expected: "5d41402abc4b2a76b9719d911017c592",
		},
		{
			input:    "",
			expected: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			input:    "你好",
			expected: "7eca689f0d3389d9dea66ae112e5cfd7",
		},
		{
			input:    "123456",
			expected: "e10adc3949ba59abbe56e057f20f883e",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := cryptor.Md5String(tt.input)
			if result != tt.expected {
				t.Errorf("Md5String(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestFileMD5_ValidFile tests valid file with known content.
func TestMd5File(t *testing.T) {
	content := []byte("hello world")
	tmpDir := t.TempDir()
	filePath := tmpDir + "/testfile.txt"

	err := os.WriteFile(filePath, content, 0600)
	if err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf("%x", md5.Sum(content))
	md5Str, err := cryptor.Md5File(filePath)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if md5Str != expected {
		t.Errorf("expected MD5 %q, got %q", expected, md5Str)
	}
}

func TestSha256String(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "hello",
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			input:    "GoLang",
			expected: "1fd260c169fb539f91b41aacb6fdbe831e899e3789187029665f96922dae2c60",
		},
		{
			input:    "!@#$%^&*()_+{}[];:'\"\\|,.<>/?`~ ",
			expected: "f277c2e6a97fd9167be6be3dd568c23c2f16702ccd3ec4af2d480e9903b3a5e8",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Input: %q", tt.input), func(t *testing.T) {
			result := cryptor.Sha256String(tt.input)
			if result != tt.expected {
				t.Errorf("Sha256String(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSha256File(t *testing.T) {
	content := []byte("hello world")
	tmpDir := t.TempDir()
	filePath := tmpDir + "/testfile.txt"

	err := os.WriteFile(filePath, content, 0600)
	if err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf("%x", sha256.Sum256(content))
	shaStr, err := cryptor.Sha256File(filePath)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if shaStr != expected {
		t.Errorf("expected Sha256 %q, got %q", expected, shaStr)
	}
}
