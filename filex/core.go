package filex

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// IsExist checks if a file or directory exists.
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

// CreateFile create a file in path.
func CreateFile(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// CreateDir create directory in absolute path. param `absPath` like /a/, /a/b/.
func CreateDir(absPath string) error {
	if absPath == "" {
		return fmt.Errorf("absPath is empty")
	}
	if IsExist(absPath) {
		return nil
	}
	err := os.MkdirAll(absPath, 0755)
	if err != nil {
		// 忽略已存在
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func CopyDir(srcPath string, dstPath string) error {
	srcInfo, err := os.Lstat(srcPath)
	if err != nil {
		return fmt.Errorf("failed to get source directory info: %w", err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("source path %s is not a directory", srcPath)
	}
	// 防止符号链接导致的循环引用或意外复制
	if (srcInfo.Mode() & os.ModeSymlink) != 0 {
		return fmt.Errorf("source path %s is a symlink", srcPath)
	}

	err = CreateDir(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// 设置目标目录权限与源目录一致
	if err := os.Chmod(dstPath, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set destination directory permissions: %w", err)
	}

	entries, err := os.ReadDir(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		srcDir := filepath.Join(srcPath, entry.Name())
		dstDir := filepath.Join(dstPath, entry.Name())

		entryInfo, err := os.Lstat(srcDir)
		if err != nil {
			return fmt.Errorf("failed to get file info for %s: %w", srcDir, err)
		}

		// 跳过符号链接
		if (entryInfo.Mode() & os.ModeSymlink) != 0 {
			continue
		}

		if entry.IsDir() {
			err = CopyDir(srcDir, dstDir)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcDir, dstDir)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyFile copy src file to dest file.
func CopyFile(srcPath string, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	distFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer distFile.Close()

	_, err = io.Copy(distFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

func IsDir(path string) (bool, error) {
	if path == "" {
		return false, fmt.Errorf("path is empty")
	}
	info, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func IsFile(path string) (bool, error) {
	if path == "" {
		return false, fmt.Errorf("path is empty")
	}
	info, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return !info.IsDir(), nil
}

// ListDir returns all files in the directory.
// It will return an error if the path is not exist.
// It will return an error if the path is not a directory.
// return all files and dirs
func ListDir(path string) ([]string, error) {
	if !IsExist(path) {
		return nil, fmt.Errorf("path %s not exist", path)
	}

	if isDir, err := IsDir(path); err != nil {
		return nil, err
	} else if !isDir {
		return nil, fmt.Errorf("path %s is not a directory", path)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
}

func RemoveFile(path string) error {
	// 清理路径，防止路径穿越攻击
	cleanedPath := filepath.Clean(path)

	info, err := os.Lstat(cleanedPath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("%s is a directory", cleanedPath)
	}
	return os.Remove(cleanedPath)
}

func RemoveDir(path string) error {
	// 清理路径，防止路径穿越攻击
	cleanedPath := filepath.Clean(path)

	info, err := os.Lstat(cleanedPath)
	if err != nil {
		return err
	}

	if (info.Mode() & os.ModeSymlink) != 0 {
		return fmt.Errorf("path %s is a symlink; removing it is not allowed", cleanedPath)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", cleanedPath)
	}

	return os.RemoveAll(cleanedPath)
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 读取当前行内容
		line := scanner.Text()
		result = append(result, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

//func ReadLinesV1(path string) ([]string, error) {
//	file, err := os.Open(path)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	var result []string
//
//	buf := bufio.NewReader(file)
//
//	for {
//		line, _, err := buf.ReadLine()
//		l := string(line)
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			continue
//		}
//		result = append(result, l)
//	}
//
//	return result, nil
//}

// ReadFileToString return string of file content.
func ReadFileToString(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading file %s: %w", path, err)
	}
	return string(bytes), nil
}

// WriteStringToFile write string to target file.
func WriteStringToFile(filepath string, content string) error {
	err := os.WriteFile(filepath, []byte(content), 0600)
	return err
}

// IsLink checks if a file is symbol link or not.
func IsLink(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeSymlink != 0
}

// FileMode return file's mode and permission.
func FileMode(path string) (fs.FileMode, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}
	return fi.Mode(), nil
}

// FileSize returns file size in bytes.
func FileSize(path string) (int64, error) {
	f, err := os.Stat(path)
	if err != nil {
		// 可选：对特定错误进行封装或记录日志
		return 0, fmt.Errorf("get file state %s: %w", path, err)
	}

	if f.Mode()&os.ModeSymlink != 0 {
		return 0, fmt.Errorf("path %s is a symlink; removing it is not allowed", path)
	}
	return f.Size(), nil
}
