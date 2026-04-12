package filex

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
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

func IsSameStat(path1, path2 string) (bool, error) {
	if path1 == "" || path2 == "" {
		return false, fmt.Errorf("path is empty")
	}

	info1, err := os.Stat(path1)
	if err != nil {
		return false, fmt.Errorf("os.Stat(%s) failed: %v", path1, err)
	}

	info2, err := os.Stat(path2)
	if err != nil {
		return false, fmt.Errorf("os.Stat(%s) failed: %v", path2, err)
	}

	return os.SameFile(info1, info2), nil
}

// GetAbsPath returns the absolute path of the given path.
// If the path is relative, it will be resolved relative to the current working directory.
func GetAbsPath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("filepath.Abs(%s) failed: %v", path, err)
	}
	return absPath, nil
}

// GetBaseName returns the last element of the path, with leading directories removed.
// It ensures that even paths ending with a separator are handled correctly.
func GetBaseName(path string) string {
	return filepath.Base(filepath.Clean(path))
}

// GetDirName 返回给定路径的目录部分，清理路径后提取其父目录。
// 注意：若路径为空或无法提取目录，则返回 "." 或对应平台的等价路径。
func GetDirName(path string) string {
	return filepath.Dir(filepath.Clean(path))
}

// GetExtension returns the file extension of the given path.
// If the path has no extension or is empty, it returns an empty string.
func GetExtension(path string) string {
	return filepath.Ext(path)
}

// SplitText returns the file name(without the extension) and extension of the given path.
// 分割路径为文件名（不含扩展）和扩展名
func SplitText(path string) (string, string) {
	cleanedPath := filepath.Clean(path)
	ext := filepath.Ext(cleanedPath)
	if ext == "" {
		return cleanedPath, ""
	}
	fileName := cleanedPath[:len(cleanedPath)-len(ext)]
	return fileName, ext
}

// CreateFile create a file in path.
func CreateFile(path string) error {
	if path == "" {
		return fmt.Errorf("path is empty")
	}

	cleanedPath := filepath.Clean(path)

	info, err := os.Lstat(cleanedPath)
	if err == nil {
		// File exists; allow only if not a symlink
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing to create file because path is a symlink")
		}
		return fmt.Errorf("file already exists")
	} else if !os.IsNotExist(err) {
		return err
	}

	dirName := filepath.Dir(cleanedPath)
	err = CreateDir(dirName)
	if err != nil {
		return fmt.Errorf("failed to create parent directories: %w", err)
	}

	f, err := os.OpenFile(cleanedPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
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

	// 规范化路径
	cleanedPath := filepath.Clean(absPath)

	// os.MkdirAll 会处理路径已存在的情况（如果它是一个目录），不会返回错误。
	// 如果路径已存在但不是目录，它会返回错误。
	err := os.MkdirAll(cleanedPath, 0750)
	if err != nil {
		// 如果错误是因为路径已作为文件存在，os.MkdirAll 会返回一个错误。
		// 我们不需要显式检查 os.IsExist(err) 然后返回 nil，因为 MkdirAll 的行为已经覆盖了这一点。
		return fmt.Errorf("failed to create directory %q: %w", cleanedPath, err)
	}
	return nil
}

func CopyDir(srcPath string, dstPath string) error {
	if srcPath == "" || dstPath == "" {
		return fmt.Errorf("srcPath or dstPath is empty")
	}
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

// CopyFile copy src file to dist file.
func CopyFile(srcPath string, dstPath string) error {
	if srcPath == "" || dstPath == "" {
		return fmt.Errorf("srcPath or dstPath is empty")
	}

	cleanedSrcPath := filepath.Clean(srcPath)
	cleanedDstPath := filepath.Clean(dstPath)

	srcInfo, err := os.Lstat(cleanedSrcPath)
	if err != nil {
		return fmt.Errorf("failed to get source directory info: %w", err)
	}

	// Prevent copying symlinks: check if srcPath is a symbolic link
	if srcInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("source file %q is a symbolic link and will not be copied", srcPath)
	}

	if dstInfo, err := os.Lstat(cleanedDstPath); err == nil {
		if dstInfo.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("destination file %q is a symbolic link and will not be overwritten", dstPath)
		}
		if os.SameFile(srcInfo, dstInfo) {
			return fmt.Errorf("source and destination refer to the same file")
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to get destination file info: %w", err)
	}

	srcFile, err := os.Open(cleanedSrcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 确保目标目录存在
	dstDir := filepath.Dir(cleanedDstPath)
	err = CreateDir(dstDir)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	distFile, err := os.OpenFile(cleanedDstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer distFile.Close()

	// 使用缓冲拷贝
	buf := make([]byte, 64*1024) // 64KB buffer
	_, err = io.CopyBuffer(distFile, srcFile, buf)
	if err != nil {
		return fmt.Errorf("error copying file content: %w", err)
	}

	if err := distFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync copied file: %w", err)
	}

	// 设置目标文件权限与源文件一致
	if err := distFile.Chmod(srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to set destination file permissions: %w", err)
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
	path = filepath.Clean(path)
	if path == "" {
		return nil, fmt.Errorf("path is empty")
	}
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
	if path == "" {
		return fmt.Errorf("path is empty")
	}
	// 清理路径，防止路径穿越攻击
	cleanedPath := filepath.Clean(path)

	info, err := os.Lstat(cleanedPath)
	if err != nil {
		return err
	}

	// Prevent remove symlinks: check if path is a symbolic link
	if (info.Mode() & os.ModeSymlink) != 0 {
		return fmt.Errorf("path %s is a symlink; removing it is not allowed", cleanedPath)
	}

	if info.IsDir() {
		return fmt.Errorf("%s is a directory", cleanedPath)
	}
	return os.Remove(cleanedPath)
}

func RemoveDir(path string) error {
	if path == "" {
		return fmt.Errorf("path is empty")
	}
	// 清理路径，防止路径穿越攻击
	cleanedPath := filepath.Clean(path)

	info, err := os.Lstat(cleanedPath)
	if err != nil {
		return err
	}

	if (info.Mode() & os.ModeSymlink) != 0 {
		return fmt.Errorf("path %s is a symlink; removing it is not allowed", cleanedPath)
	}

	// Walk the directory and refuse to remove if symlink is found inside
	err = filepath.WalkDir(cleanedPath, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing to remove %s: contains symlink at %s", cleanedPath, p)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", cleanedPath)
	}

	return os.RemoveAll(cleanedPath)
}

func ReadLines(path string) ([]string, error) {
	if path == "" {
		return nil, fmt.Errorf("path is empty")
	}
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		// Wrap scanner error for more context
		return nil, fmt.Errorf("error scanning file %s: %w", filepath.Clean(path), err)
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
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}
	bytes, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", fmt.Errorf("reading file %s: %w", path, err)
	}
	return string(bytes), nil
}

// WriteStringToFile write string to target file.
func WriteStringToFile(filepath string, content string) error {
	if filepath == "" {
		return fmt.Errorf("filepath is empty")
	}
	err := os.WriteFile(filepath, []byte(content), 0600)
	return err
}

// IsLink checks if a file is symbol link or not.
func IsLink(path string) (bool, error) {
	if path == "" {
		return false, fmt.Errorf("path is empty")
	}
	fi, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return fi.Mode()&os.ModeSymlink != 0, nil
}

// FileMode return file's mode and permission.
func FileMode(path string) (fs.FileMode, error) {
	if path == "" {
		return 0, fmt.Errorf("path is empty")
	}
	fi, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}
	return fi.Mode(), nil
}

// FileSize returns file size in bytes.
func FileSize(path string) (int64, error) {
	if path == "" {
		return 0, fmt.Errorf("path is empty")
	}
	info, err := os.Lstat(path) // Use Lstat to get info about the link/file itself, not following symlinks
	if err != nil {
		return 0, fmt.Errorf("failed to get file info for %s: %w", path, err)
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return 0, fmt.Errorf("path '%s' is a symbolic link, operation not permitted on symlinks", path)
	}

	// For directories, Size() is platform-dependent and usually not the sum of contents.
	// This behavior is consistent with os.FileInfo.Size().
	return info.Size(), nil
}

// WalkResult 包含目录遍历结果
type WalkResult struct {
	Root  string   // 当前目录路径
	Dirs  []string // 子目录名列表（相对路径）
	Files []string // 文件名列表
	Err   error    // 遍历错误
}

func WalkV1(root string) []WalkResult {
	var results []WalkResult
	stack := []string{root} // 使用栈实现深度优先遍历

	for len(stack) > 0 {
		currentDir := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// 读取目录内容
		entries, err := os.ReadDir(currentDir)
		if err != nil {
			results = append(results, WalkResult{
				Root: currentDir,
				Err:  err,
			})
			continue
		}

		// 分离子目录和文件
		var dirs, files []string
		for _, entry := range entries {
			name := entry.Name()
			if entry.IsDir() {
				dirs = append(dirs, name)
			} else {
				files = append(files, name)
			}
		}

		// 按名称排序（与 Python 行为一致）
		sort.Strings(dirs)
		sort.Strings(files)

		// 记录结果
		results = append(results, WalkResult{
			Root:  currentDir,
			Dirs:  dirs,
			Files: files,
		})

		// 逆序压栈（保证顺序处理）
		for i := len(dirs) - 1; i >= 0; i-- {
			fullPath := filepath.Join(currentDir, dirs[i])
			stack = append(stack, fullPath)
		}
	}

	return results
}

func WalkV2(root string) ([]*WalkResult, error) {
	var results []*WalkResult
	resultMap := make(map[string]*WalkResult)

	err := filepath.Walk(root, func(filePath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if filePath == root {
			return nil
		}

		dirName := filepath.Dir(filePath)
		dirWr, ok := resultMap[dirName]

		if !ok {
			dirWr = &WalkResult{
				Root:  dirName,
				Dirs:  []string{},
				Files: []string{},
			}
			results = append(results, dirWr)
			resultMap[dirName] = dirWr
		}

		if info.IsDir() {
			dirWr.Dirs = append(dirWr.Dirs, info.Name())

			if _, ok := resultMap[filePath]; !ok {
				currentWr := &WalkResult{
					Root:  filePath,
					Dirs:  []string{},
					Files: []string{},
				}
				results = append(results, currentWr)
				resultMap[filePath] = currentWr
			}
		} else {
			dirWr.Files = append(dirWr.Files, info.Name())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}
