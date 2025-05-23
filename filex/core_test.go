package filex_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/shijl0925/go-toolkits/filex"
	"github.com/shijl0925/go-toolkits/slicex"
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

			if result == "" {
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

// TestWalk_NormalDirectoryStructure tests the Walk function with a normal directory structure
func TestWalk_NormalDirectoryStructure(t *testing.T) {
	// 创建临时测试目录结构
	testDir := t.TempDir()

	// 创建子目录和文件
	os.MkdirAll(filepath.Join(testDir, "dir1", "subdir1"), 0755)
	os.MkdirAll(filepath.Join(testDir, "dir2"), 0755)
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(testDir, "file2.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(testDir, "dir1", "subdir1", "subfile.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(testDir, "dir1", "file3.txt"), []byte("content"), 0644)

	// 执行Walk函数
	results, err := filex.WalkV2(testDir)
	if err != nil {
		t.Errorf("WalkV2() error = %v", err)
	}

	// 验证结果数量
	expectedCount := 4 // 根目录 + dir1 + subdir1 + dir2
	if len(results) != expectedCount {
		t.Errorf("Expected %d results, got %d", expectedCount, len(results))
	}

	// 验证根目录结果
	for _, result := range results {
		if result.Err != nil {
			t.Errorf("Unexpected error in directory %s: %v", result.Root, result.Err)
		}

		switch {
		case result.Root == testDir:
			expectedDirs := []string{"dir1", "dir2"}
			expectedFiles := []string{"file1.txt", "file2.txt"}
			if !slicex.EqualUnordered(result.Dirs, expectedDirs) || !slicex.EqualUnordered(result.Files, expectedFiles) {
				t.Errorf("Root directory mismatch: Dirs=%v, Files=%v", result.Dirs, result.Files)
			}
		case filepath.Base(result.Root) == "dir1":
			expectedDirs := []string{"subdir1"}
			expectedFiles := []string{"file3.txt"}
			if !slicex.EqualUnordered(result.Dirs, expectedDirs) || !slicex.EqualUnordered(result.Files, expectedFiles) {
				t.Errorf("Dir1 directory mismatch: Dirs=%v, Files=%v", result.Dirs, result.Files)
			}
		case filepath.Base(result.Root) == "subdir1":
			var expectedDirs []string

			expectedFiles := []string{"subfile.txt"}
			if !slicex.EqualUnordered(result.Dirs, expectedDirs) || !slicex.EqualUnordered(result.Files, expectedFiles) {
				t.Errorf("Subdir1 directory mismatch: Dirs=%v, Files=%v", result.Dirs, result.Files)
			}
		case filepath.Base(result.Root) == "dir2":
			var expectedDirs []string

			var expectedFiles []string

			if !slicex.EqualUnordered(result.Dirs, expectedDirs) || !slicex.EqualUnordered(result.Files, expectedFiles) {
				t.Errorf("Dir2 directory mismatch: Dirs=%v, Files=%v", result.Dirs, result.Files)
			}
		}
	}
}

// TestWalk_NonExistentDirectory tests Walk function with a non-existent directory
func TestWalk_NonExistentDirectory(t *testing.T) {
	nonExistentDir := filepath.Join("non", "existent", "dir")
	_, err := filex.WalkV2(nonExistentDir)
	if err == nil {
		t.Errorf("Walk error: %v", err)
	}
}

// TestWalk_EmptyDirectory tests Walk function with an empty directory
func TestWalk_EmptyDirectory(t *testing.T) {
	emptyDir := t.TempDir()
	results, err := filex.WalkV2(emptyDir)
	if err != nil {
		t.Errorf("Walk error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 result, got %d", len(results))
	}
}

// TestWalk_DirectoryOnlyWithSubdirectories tests Walk function with directories only
func TestWalk_DirectoryOnlyWithSubdirectories(t *testing.T) {
	testDir := t.TempDir()
	os.MkdirAll(filepath.Join(testDir, "dir1"), 0755)
	os.MkdirAll(filepath.Join(testDir, "dir2"), 0755)

	results, err := filex.WalkV2(testDir)
	if err != nil {
		t.Errorf("WalkV2() error = %v", err)
	}

	if len(results) != 3 { // 根目录 + dir1 + dir2
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	for _, result := range results {
		if result.Err != nil {
			t.Errorf("Unexpected error in directory %s: %v", result.Root, result.Err)
		}

		if len(result.Files) != 0 {
			t.Errorf("Expected no files in directory %s, got %v", result.Root, result.Files)
		}
	}

	// 验证根目录的子目录
	if !slicex.EqualUnordered(results[0].Dirs, []string{"dir1", "dir2"}) {
		t.Errorf("Expected Dirs=[dir1, dir2], got %v", results[0].Dirs)
	}
}

// TestWalk_DirectoryOnlyWithFiles tests Walk function with files only
func TestWalk_DirectoryOnlyWithFiles(t *testing.T) {
	testDir := t.TempDir()
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(testDir, "file2.txt"), []byte("content"), 0644)

	results, err := filex.WalkV2(testDir)
	if err != nil {
		t.Errorf("WalkV2() error = %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	if results[0].Err != nil {
		t.Errorf("Unexpected error: %v", results[0].Err)
	}

	if len(results[0].Dirs) != 0 {
		t.Errorf("Expected no subdirectories, got %v", results[0].Dirs)
	}

	if !slicex.EqualUnordered(results[0].Files, []string{"file1.txt", "file2.txt"}) {
		t.Errorf("Expected Files=[file1.txt, file2.txt], got %v", results[0].Files)
	}
}
