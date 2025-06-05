package cryptor

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

const maxFileSize = 1 << 30 // 1GB max file size (adjustable as needed)

// Base64StdEncode encode string with base64 encoding.
func Base64StdEncode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// Base64StdDecode decode a base64 encoded string.
func Base64StdDecode(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 string: %w", err)
	}
	return string(b), nil
}

// Md5Stream returns md5 hash of a string.
// 适合流式处理大文件
func Md5Stream(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Md5String returns the MD5 hash of the input string as a hexadecimal string.
func Md5String(s string) string {
	// Convert string to bytes
	data := []byte(s)
	// Compute MD5 hash
	h := md5.Sum(data)
	// Encode to hexadecimal string and return
	return hex.EncodeToString(h[:])
}

// Md5File returns the MD5 hash of a file as a hexadecimal string.
func Md5File(filePath string) (string, error) {
	if len(strings.TrimSpace(filePath)) == 0 {
		return "", fmt.Errorf("invalid file path: empty or whitespace only")
	}

	stat, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %w", err)
	}
	if stat.IsDir() {
		return "", fmt.Errorf("file is a directory")
	}
	if stat.Size() > maxFileSize {
		return "", fmt.Errorf("file size exceeds limit of %d bytes", maxFileSize)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q: %w", filePath, err)
	}
	defer file.Close()

	h := md5.New()

	// Optional: use a buffer for better performance on large files
	buf := make([]byte, 64*1024) // 64KB buffer
	if _, err := io.CopyBuffer(h, file, buf); err != nil {
		return "", fmt.Errorf("error reading file %q: %w", filePath, err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// Sha256String returns the SHA256 hash of the input string as a hexadecimal string.
func Sha256String(s string) string {
	// Convert string to bytes
	data := []byte(s)

	// Compute MD5 hash
	h := sha256.Sum256(data)

	return hex.EncodeToString(h[:])
}

// Sha256File returns the sha256 hash of a file as a hexadecimal string.
func Sha256File(filePath string) (string, error) {
	if len(strings.TrimSpace(filePath)) == 0 {
		return "", fmt.Errorf("invalid file path: empty or whitespace only")
	}

	stat, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %w", err)
	}
	if stat.IsDir() {
		return "", fmt.Errorf("file is a directory")
	}
	if stat.Size() > maxFileSize {
		return "", fmt.Errorf("file size exceeds limit of %d bytes", maxFileSize)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %q: %w", filePath, err)
	}
	defer file.Close()

	sh := sha256.New()

	// Optional: use a buffer for better performance on large files
	buf := make([]byte, 64*1024) // 64KB buffer
	if _, err := io.CopyBuffer(sh, file, buf); err != nil {
		return "", fmt.Errorf("error reading file %q: %w", filePath, err)
	}
  
	return hex.EncodeToString(sh.Sum(nil)), nil
}
