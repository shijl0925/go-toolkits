package filex_test

import (
	"github.com/shijl0925/go-toolkits/filex"
	"os"
	"path/filepath"
	"testing"
)

func TestIsExist(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "Path exists",
			path:     "../jsonx/testdata/test.json",
			expected: true,
		},
		{
			name:     "Path does not exist",
			path:     "testdata/test.json",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filex.IsExist(tt.path)
			if result != tt.expected {
				t.Errorf("IsExist() expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCreateFile(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "成功创建新文件",
			path:    filepath.Join(tempDir, "newfile.txt"),
			wantErr: false,
		},
		{
			name:    "文件已存在",
			path:    filepath.Join(tempDir, "existfile.txt"),
			wantErr: true,
		},
		{
			name:    "无效路径",
			path:    "/invalid/path/testfile.txt",
			wantErr: true,
		},
		{
			name:    "空路径",
			path:    "",
			wantErr: true,
		},
	}

	// 创建一个已存在的文件用于测试
	existFilePath := filepath.Join(tempDir, "existfile.txt")
	f, err := os.Create(existFilePath)
	if err != nil {
		t.Fatalf("准备测试文件失败: %v", err)
	}
	f.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := filex.CreateFile(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// 检查文件是否真的创建成功
				if _, err := os.Stat(tt.path); os.IsNotExist(err) {
					t.Errorf("期望文件存在，但未找到: %s", tt.path)
				}
			}
		})
	}
}

func TestCreateDir(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "empty path",
			input:   "",
			wantErr: true,
		},
		{
			name:    "dir already exists via IsExist",
			input:   "/tmp/existing",
			wantErr: false,
		},
		{
			name:    "mkdir success",
			input:   "/tmp/newdir",
			wantErr: false,
		},
	}

	// 创建一个已存在的文件用于测试
	existFilePath := filepath.Join("tmp", "existing")
	err := os.MkdirAll(existFilePath, 0755)
	if err != nil {
		t.Fatalf("准备测试文件夹失败: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := filex.CreateDir(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
