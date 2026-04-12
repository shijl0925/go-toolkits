package filex_test

import (
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

// Test_IsSameStat_EmptyPath 测试空路径输入
func Test_IsSameStat_EmptyPath(t *testing.T) {
	_, err := filex.IsSameStat("", "any")
	if err == nil {
		t.Errorf("IsSameStat() expected error, got nil")
	}

	_, err = filex.IsSameStat("any", "")
	if err == nil {
		t.Errorf("IsSameStat() expected error, got nil")
	}
}

// Test_IsSameStat_StatFailFirst 测试第一个路径 Stat 失败
func Test_IsSameStat_StatFailFirst(t *testing.T) {
	_, err := filex.IsSameStat("nonexistent_file1", "nonexistent_file2")
	if err == nil {
		t.Errorf("IsSameStat() expected error, got nil")
	}
}

// Test_IsSameStat_StatFailSecond 测试第二个路径 Stat 失败
func Test_IsSameStat_StatFailSecond(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Errorf("os.CreateTemp() failed: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = filex.IsSameStat(tmpFile.Name(), "nonexistent_file2")
	if err == nil {
		t.Errorf("IsSameStat() expected error, got nil")
	}
}

// Test_IsSameStat_SameFile 测试两个路径指向同一文件
func Test_IsSameStat_SameFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Errorf("os.CreateTemp() failed: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	isSame, err := filex.IsSameStat(tmpFile.Name(), tmpFile.Name())
	if err != nil || !isSame {
		t.Errorf("IsSameStat() expect %v, get %v", true, isSame)
	}
}

// Test_IsSameStat_DifferentFiles 测试两个路径指向不同文件
func Test_IsSameStat_DifferentFiles(t *testing.T) {
	tmpFile1, err := os.CreateTemp("", "testfile1")
	if err != nil {
		t.Errorf("os.CreateTemp() failed: %v", err)
	}
	defer os.Remove(tmpFile1.Name())
	defer tmpFile1.Close()

	tmpFile2, err := os.CreateTemp("", "testfile2")
	if err != nil {
		t.Errorf("os.CreateTemp() failed: %v", err)
	}
	defer os.Remove(tmpFile2.Name())
	defer tmpFile2.Close()

	isSame, err := filex.IsSameStat(tmpFile1.Name(), tmpFile2.Name())
	if err != nil || isSame {
		t.Errorf("IsSameStat() expect %v, get %v", false, isSame)
	}
}

// TestGetAbsPath_EmptyPath 检查当输入为空字符串时是否返回正确错误
func TestGetAbsPath_EmptyPath(t *testing.T) {
	_, err := filex.GetAbsPath("")
	if err == nil {
		t.Errorf("GetAbsPath() expected error, got nil")
	}
}

// TestGetAbsPath_ValidRelativePath 检查相对路径是否能正确转换为绝对路径
func TestGetAbsPath_ValidRelativePath(t *testing.T) {
	cwd, _ := os.Getwd()
	expectedPath := filepath.Join(cwd, "test")

	absPath, err := filex.GetAbsPath("test")
	if err != nil {
		t.Errorf("GetAbsPath() failed: %v", err)
	}
	if absPath != expectedPath {
		t.Errorf("GetAbsPath() expected %s, got %s", expectedPath, absPath)
	}
}

//// TestGetAbsPath_InvalidPath 模拟非法路径导致 filepath.Abs 失败的情况
//func TestGetAbsPath_InvalidPath(t *testing.T) {
//	// 使用一个不可能存在的路径来触发错误（例如 Windows 不允许的路径）
//	invalidPath := `<invalid>:path?*"|`
//	_, err := filex.GetAbsPath(invalidPath)
//
//	if err == nil {
//		t.Errorf("GetAbsPath() expected error, got nil")
//	}
//}

// TestGetBaseName tests various cases for GetBaseName function.
func TestGetBaseName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal file path",
			input:    "/home/user/file.txt",
			expected: "file.txt",
		},
		{
			name:     "Directory path with trailing slash",
			input:    "/home/user/",
			expected: "user",
		},
		{
			name:     "Current directory",
			input:    ".",
			expected: ".",
		},
		{
			name:     "Parent directory",
			input:    "..",
			expected: "..",
		},
		{
			name:     "Relative file path",
			input:    "file.txt",
			expected: "file.txt",
		},
		//{
		//	name:     "Windows-style path on Unix",
		//	input:    `C:\Users\file.txt`,
		//	expected: "file.txt",
		//},
		{
			name:     "Multi-level path",
			input:    "a/b/c",
			expected: "c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filex.GetBaseName(tt.input)
			if result != tt.expected {
				t.Errorf("GetBaseName(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetDirName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal file path",
			input:    "/home/user/file.txt",
			expected: "/home/user",
		},
		{
			name:     "Directory path with trailing slash",
			input:    "/home/user/",
			expected: "/home",
		},
		{
			name:     "Current directory",
			input:    ".",
			expected: ".",
		},
		{
			name:     "Relative file path",
			input:    "file.txt",
			expected: ".",
		},
		//{
		//	name:     "Windows-style path on Unix",
		//	input:    `C:\Users\file.txt`,
		//	expected: "C:\Users",
		//},
		{
			name:     "Multi-level path",
			input:    "a/b/c",
			expected: "a/b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filex.GetDirName(tt.input)
			if result != tt.expected {
				t.Errorf("GetDirName(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestGetExtension tests various cases for the GetExtension function.
func TestGetExtension(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Normal file with extension", "file.txt", ".txt"},
		{"File without extension", "file", ""},
		{"Multiple dots in filename", "image.version.jpg", ".jpg"},
		{"Filename ends with a dot", "file.", "."},
		{"Path with directories and multiple extensions", "dir/file.tar.gz", ".gz"},
		{"Empty string input", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filex.GetExtension(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q but got %q for input %q", tt.expected, result, tt.input)
			}
		})
	}
}

func TestSplitText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected [2]string
	}{
		{"TC01 - 常规带扩展名", "/home/user/file.txt", [2]string{"/home/user/file", ".txt"}},
		{"TC02 - 多点扩展名", "file.tar.gz", [2]string{"file.tar", ".gz"}},
		{"TC03 - 简单扩展名", "image.png", [2]string{"image", ".png"}},
		{"TC04 - 无扩展名", "/tmp/no_extension", [2]string{"/tmp/no_extension", ""}},
		{"TC05 - 只有文件名", "README", [2]string{"README", ""}},
		//{"TC06 - 当前目录", ".", [2]string{"", "."}},
		//{"TC07 - 上级目录", "..", [2]string{".", "."}},
		{"TC08 - 结尾带斜杠", "/a/b/c/d/", [2]string{"/a/b/c/d", ""}},
		{"TC09 - 路径含扩展名但结尾是斜杠", "/a/b/c.txt/", [2]string{"/a/b/c", ".txt"}},
		{"TC10 - Windows路径", "C:\\Windows\\test.txt", [2]string{"C:\\Windows\\test", ".txt"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFileName, gotExt := filex.SplitText(tt.input)
			if tt.expected != [2]string{gotFileName, gotExt} {
				t.Errorf("SplitText(%q) = %q, %q; want %q, %q", tt.input, gotFileName, gotExt, tt.expected[0], tt.expected[1])
			}
			//assert.Equal(t, tt.expected[0], gotFileName)
			//assert.Equal(t, tt.expected[1], gotExt)
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
		{
			name:    "父文件夹不存在的情况",
			path:    filepath.Join(tempDir, "existing", "newfile.txt"),
			wantErr: false,
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
	err := os.MkdirAll(existFilePath, 0750)
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
		name string
		path string
	}{
		{
			name: "test01",
			path: "../LICENSE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := os.Lstat(tt.path)
			if err != nil {
				t.Fatalf("os.Lstat() error = %v", err)
			}

			result, err := filex.FileMode(tt.path)
			if err != nil {
				t.Errorf("FileMode() error = %v", err)
			}

			if result != info.Mode() {
				t.Errorf("FileMode() expected %v, got %v", info.Mode(), result)
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
		err := os.MkdirAll(existFilePath01, 0750)
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
		err := os.MkdirAll(existFilePath01, 0750)
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

	err := os.MkdirAll("tmp", 0750)
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

func Test_CopyFile_PreservesMode(t *testing.T) {
	tempDir := t.TempDir()
	srcPath := filepath.Join(tempDir, "source.txt")
	dstPath := filepath.Join(tempDir, "nested", "dest.txt")

	if err := os.WriteFile(srcPath, []byte("copy me"), 0600); err != nil {
		t.Fatalf("os.WriteFile() failed: %v", err)
	}

	if err := os.Chmod(srcPath, 0640); err != nil {
		t.Fatalf("os.Chmod() failed: %v", err)
	}

	if err := filex.CopyFile(srcPath, dstPath); err != nil {
		t.Fatalf("CopyFile() error = %v", err)
	}

	srcInfo, err := os.Lstat(srcPath)
	if err != nil {
		t.Fatalf("os.Lstat(src) failed: %v", err)
	}

	dstInfo, err := os.Lstat(dstPath)
	if err != nil {
		t.Fatalf("os.Lstat(dst) failed: %v", err)
	}

	if dstInfo.Mode() != srcInfo.Mode() {
		t.Errorf("CopyFile() mode = %v; want %v", dstInfo.Mode(), srcInfo.Mode())
	}
}

func Test_CopyFile_RejectsSameFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "same.txt")
	originalContent := []byte("original content")

	if err := os.WriteFile(filePath, originalContent, 0600); err != nil {
		t.Fatalf("os.WriteFile() failed: %v", err)
	}

	err := filex.CopyFile(filePath, filePath)
	if err == nil {
		t.Fatal("CopyFile() expected error when source and destination are the same file")
	}

	content, readErr := os.ReadFile(filePath)
	if readErr != nil {
		t.Fatalf("os.ReadFile() failed: %v", readErr)
	}

	if string(content) != string(originalContent) {
		t.Errorf("CopyFile() modified source content: got %q; want %q", content, originalContent)
	}
}

func Test_CopyFile_RejectsHardLinkDestination(t *testing.T) {
	tempDir := t.TempDir()
	srcPath := filepath.Join(tempDir, "source.txt")
	hardLinkPath := filepath.Join(tempDir, "source-hardlink.txt")
	originalContent := []byte("hard link content")

	if err := os.WriteFile(srcPath, originalContent, 0600); err != nil {
		t.Fatalf("os.WriteFile(src) failed: %v", err)
	}

	if err := os.Link(srcPath, hardLinkPath); err != nil {
		t.Fatalf("os.Link() failed: %v", err)
	}

	err := filex.CopyFile(srcPath, hardLinkPath)
	if err == nil {
		t.Fatal("CopyFile() expected error when destination is a hard link to the source")
	}

	content, readErr := os.ReadFile(srcPath)
	if readErr != nil {
		t.Fatalf("os.ReadFile(src) failed: %v", readErr)
	}

	if string(content) != string(originalContent) {
		t.Errorf("CopyFile() modified source content through hard link destination: got %q; want %q", content, originalContent)
	}
}

func Test_CopyFile_RejectsDestinationSymlink(t *testing.T) {
	tempDir := t.TempDir()
	srcPath := filepath.Join(tempDir, "source.txt")
	targetPath := filepath.Join(tempDir, "target.txt")
	linkPath := filepath.Join(tempDir, "dest-link.txt")

	if err := os.WriteFile(srcPath, []byte("source content"), 0600); err != nil {
		t.Fatalf("os.WriteFile(src) failed: %v", err)
	}

	originalTarget := []byte("target content")
	if err := os.WriteFile(targetPath, originalTarget, 0600); err != nil {
		t.Fatalf("os.WriteFile(target) failed: %v", err)
	}

	if err := os.Symlink(targetPath, linkPath); err != nil {
		t.Fatalf("os.Symlink() failed: %v", err)
	}

	err := filex.CopyFile(srcPath, linkPath)
	if err == nil {
		t.Fatal("CopyFile() expected error when destination is a symbolic link")
	}

	content, readErr := os.ReadFile(targetPath)
	if readErr != nil {
		t.Fatalf("os.ReadFile(target) failed: %v", readErr)
	}

	if string(content) != string(originalTarget) {
		t.Errorf("CopyFile() unexpectedly modified symlink target: got %q; want %q", content, originalTarget)
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

	err := os.MkdirAll("tmp", 0750)
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
		err := os.MkdirAll(existFilePath01, 0750)
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
		err := os.MkdirAll(existFilePath01, 0750)
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
		err := os.MkdirAll(existFilePath01, 0750)
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
	os.MkdirAll(filepath.Join(testDir, "dir1", "subdir1"), 0750)
	os.MkdirAll(filepath.Join(testDir, "dir2"), 0750)
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
	os.MkdirAll(filepath.Join(testDir, "dir1"), 0750)
	os.MkdirAll(filepath.Join(testDir, "dir2"), 0750)

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
