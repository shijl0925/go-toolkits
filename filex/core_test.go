package filex_test

import (
	"github.com/shijl0925/go-toolkits/filex"
	"io/fs"
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

func TestReadLines(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected int
	}{
		{
			name:     "test01",
			path:     "../LICENSE",
			expected: 21,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := filex.ReadLines(tt.path)
			if err != nil {
				t.Errorf("ReadLines() error = %v", err)
			}

			if len(result) != tt.expected {
				t.Errorf("ReadLines() expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestReadFileToString(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected int
	}{
		{
			name:     "test01",
			path:     "../LICENSE",
			expected: 21,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := filex.ReadFileToString(tt.path)
			if err != nil {
				t.Errorf("ReadFileToString() error = %v", err)
			}

			if len(result) == 0 {
				t.Errorf("ReadFileToString() expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func Test_WriteStringToFile(t *testing.T) {
	t.Run("TC01 - Valid file", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test.txt")
		content := `{"first_name":"Alice","last_name":"Smith","age":18}`
		err := filex.WriteStringToFile(filePath, content)
		if err != nil {
			t.Errorf("WriteStringToFile() error: %v", err)
		}

		// 验证文件内容
		data, _ := os.ReadFile(filePath)
		if content != string(data) {
			t.Errorf("WriteStringToFile() expected %v, got %v", content, data)
		}
	})
}

func TestFileMode(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected fs.FileMode
	}{
		{
			name:     "test01",
			path:     "../LICENSE",
			expected: fs.FileMode(0644),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := filex.FileMode(tt.path)
			if err != nil {
				t.Errorf("FileMode() error = %v", err)
			}

			if result != tt.expected {
				t.Errorf("FileMode() expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFileSize(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected int64
	}{
		{
			name:     "test01",
			path:     "../LICENSE",
			expected: 1071,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := filex.FileSize(tt.path)
			if err != nil {
				t.Errorf("FileSize() error = %v", err)
			}

			if result != tt.expected {
				t.Errorf("FileSize() expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func Test_IsDir(t *testing.T) {
	t.Run("test01", func(t *testing.T) {
		// 创建一个已存在的文件用于测试
		existFilePath01 := filepath.Join("tmp", "existing")
		err := os.MkdirAll(existFilePath01, 0755)
		if err != nil {
			t.Fatalf("准备测试文件夹失败: %v", err)
		}
		defer os.RemoveAll(existFilePath01)

		existFilePath02 := filepath.Join("tmp", "existfile.txt")
		_, err = os.Create(existFilePath02)
		if err != nil {
			t.Fatalf("准备测试文件失败: %v", err)
		}
		defer os.Remove(existFilePath02)

		result01, err := filex.IsDir("tmp/existing")
		if err != nil {
			t.Errorf("IsDir() error = %v", err)
		}
		if !result01 {
			t.Errorf("IsDir() expected %v, got %v", true, result01)
		}
		result02, err := filex.IsDir("tmp/existfile.txt")
		if err != nil {
			t.Errorf("IsDir() error = %v", err)
		}
		if result02 {
			t.Errorf("IsDir() expected %v, got %v", false, result02)
		}
	})
}

func Test_IsFile(t *testing.T) {
	t.Run("test01", func(t *testing.T) {
		// 创建一个已存在的文件用于测试
		existFilePath01 := filepath.Join("tmp", "existing")
		err := os.MkdirAll(existFilePath01, 0755)
		if err != nil {
			t.Fatalf("准备测试文件夹失败: %v", err)
		}
		defer os.RemoveAll(existFilePath01)

		existFilePath02 := filepath.Join("tmp", "existfile.txt")
		_, err = os.Create(existFilePath02)
		if err != nil {
			t.Fatalf("准备测试文件失败: %v", err)
		}
		defer os.Remove(existFilePath02)

		result01, err := filex.IsFile("tmp/existing")
		if err != nil {
			t.Errorf("IsFile() error = %v", err)
		}
		if result01 {
			t.Errorf("IsFile() expected %v, got %v", true, result01)
		}
		result02, err := filex.IsFile("tmp/existfile.txt")
		if err != nil {
			t.Errorf("IsFile() error = %v", err)
		}
		if !result02 {
			t.Errorf("IsFile() expected %v, got %v", false, result02)
		}
	})
}

func Test_CopyFile(t *testing.T) {
	tests := []struct {
		name     string
		srcPath  string
		dstPath  string
		expected int
	}{
		{
			name:    "test01",
			srcPath: "../LICENSE",
			dstPath: "tmp/LICENSE",
		},
	}

	err := os.MkdirAll("tmp", 0755)
	if err != nil {
		t.Fatalf("准备测试文件夹失败: %v", err)
	}
	defer os.Remove("tmp/LICENSE")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := filex.CopyFile(tt.srcPath, tt.dstPath)
			if err != nil {
				t.Errorf("CopyFile() error = %v", err)
			}
		})
	}
}

func Test_CopyDir(t *testing.T) {
	tests := []struct {
		name     string
		srcPath  string
		dstPath  string
		expected int
	}{
		{
			name:    "test01",
			srcPath: "../.github",
			dstPath: "tmp/.github",
		},
	}

	err := os.MkdirAll("tmp", 0755)
	if err != nil {
		t.Fatalf("准备测试文件夹失败: %v", err)
	}
	defer os.RemoveAll("tmp/.github")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := filex.CopyDir(tt.srcPath, tt.dstPath)
			if err != nil {
				t.Errorf("CopyDir() error = %v", err)
			}
		})
	}
}

func Test_ListDir(t *testing.T) {
	t.Run("test01", func(t *testing.T) {
		// 创建一个已存在的文件用于测试
		existFilePath01 := filepath.Join("tmp", "existing")
		err := os.MkdirAll(existFilePath01, 0755)
		if err != nil {
			t.Fatalf("准备测试文件夹失败: %v", err)
		}
		defer os.RemoveAll(existFilePath01)

		existFilePath02 := filepath.Join("tmp", "existfile.txt")
		_, err = os.Create(existFilePath02)
		if err != nil {
			t.Fatalf("准备测试文件失败: %v", err)
		}
		defer os.Remove(existFilePath02)

		result, err := filex.ListDir("tmp")
		if err != nil {
			t.Errorf("ListDir() error = %v", err)
		}
		if len(result) != 2 {
			t.Errorf("ListDir() expected %v, got %v", 2, len(result))
		}
	})
}

func Test_RemoveDir(t *testing.T) {
	t.Run("test01", func(t *testing.T) {
		// 创建一个已存在的文件用于测试
		existFilePath01 := filepath.Join("tmp", "existing")
		err := os.MkdirAll(existFilePath01, 0755)
		if err != nil {
			t.Fatalf("准备测试文件夹失败: %v", err)
		}
		defer os.RemoveAll(existFilePath01)

		existFilePath02 := filepath.Join("tmp", "existfile.txt")
		_, err = os.Create(existFilePath02)
		if err != nil {
			t.Fatalf("准备测试文件失败: %v", err)
		}
		defer os.Remove(existFilePath02)

		err = filex.RemoveDir("tmp/existing")
		if err != nil {
			t.Errorf("RemoveDir() error = %v", err)
		}

		err = filex.RemoveDir("tmp/existfile.txt")
		if err == nil {
			t.Errorf("RemoveDir() error = %v", err)
		}
	})
}

func Test_RemoveFile(t *testing.T) {
	t.Run("test01", func(t *testing.T) {
		// 创建一个已存在的文件用于测试
		existFilePath01 := filepath.Join("tmp", "existing")
		err := os.MkdirAll(existFilePath01, 0755)
		if err != nil {
			t.Fatalf("准备测试文件夹失败: %v", err)
		}
		defer os.RemoveAll(existFilePath01)

		existFilePath02 := filepath.Join("tmp", "existfile.txt")
		_, err = os.Create(existFilePath02)
		if err != nil {
			t.Fatalf("准备测试文件失败: %v", err)
		}
		defer os.Remove(existFilePath02)

		err = filex.RemoveFile("tmp/existing")
		if err == nil {
			t.Errorf("RemoveFile() error = %v", err)
		}

		err = filex.RemoveFile("tmp/existfile.txt")
		if err != nil {
			t.Errorf("RemoveFile() error = %v", err)
		}
	})
}
