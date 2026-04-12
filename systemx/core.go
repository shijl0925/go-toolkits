package systemx

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type (
	Option func(*exec.Cmd)
)

// IsWindows check if current os is windows.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsLinux check if current os is linux.
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsMac check if current os is macos.
func IsMac() bool {
	return runtime.GOOS == "darwin"
}

// GetOsEnv gets the value of the environment variable named by the key.
func GetOsEnv(key string) string {
	return os.Getenv(key)
}

// SetOsEnv sets the value of the environment variable named by the key.
func SetOsEnv(key, value string) error {
	return os.Setenv(key, value)
}

// RemoveOsEnv remove a single environment variable.
func RemoveOsEnv(key string) error {
	return os.Unsetenv(key)
}

// ExecCommand execute command, return the stdout and stderr string of command, and error if error occur
// param `command` is a complete command string, like, ls -a (linux), dir(windows), ping 127.0.0.1
// in linux,  use /bin/bash -c to execute command
// in windows, use powershell.exe to execute command
func ExecCommand(command string, opts ...Option) (stdout, stderr string, err error) {
	cmdPath, args, err := resolveCommand(command)
	if err != nil {
		return "", "", err
	}

	var out bytes.Buffer
	var errOut bytes.Buffer

	cmd := exec.Command(cmdPath, args...) // #nosec G204 -- executable path is validated and resolved before execution

	for _, opt := range opts {
		if opt != nil {
			opt(cmd)
		}
	}
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err = cmd.Run()

	if err != nil {
		stderr = errOut.String()
		return "", stderr, err
	}

	stdout = out.String()

	return stdout, "", nil
}

func StartProcess(command string, args ...string) (int, error) {
	cmdPath, err := resolveExecutable(command)
	if err != nil {
		return 0, err
	}

	cmd := exec.Command(cmdPath, args...) // #nosec G204 -- executable path is validated and resolved before execution

	if err := cmd.Start(); err != nil {
		return 0, err
	}

	return cmd.Process.Pid, nil
}

// StopProcess stop a process by pid.
func StopProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return process.Signal(os.Kill)
}

// KillProcess kill a process by pid.
func KillProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return process.Kill()
}

func resolveCommand(command string) (string, []string, error) {
	parts, err := splitCommandLine(command)
	if err != nil {
		return "", nil, err
	}

	cmdPath, err := resolveExecutable(parts[0])
	if err != nil {
		return "", nil, err
	}

	return cmdPath, parts[1:], nil
}

func resolveExecutable(command string) (string, error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return "", fmt.Errorf("empty command provided")
	}

	if strings.ContainsRune(command, 0) {
		return "", fmt.Errorf("invalid command provided")
	}

	if filepath.IsAbs(command) {
		return filepath.Clean(command), nil
	}

	if strings.ContainsRune(command, filepath.Separator) {
		return "", fmt.Errorf("relative command paths are not allowed")
	}

	path, err := exec.LookPath(command)
	if err != nil {
		return "", err
	}

	return path, nil
}

func splitCommandLine(command string) ([]string, error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return nil, fmt.Errorf("empty command provided")
	}

	var (
		args     []string
		current  strings.Builder
		quote    rune
		escaped  bool
		flushArg = func() {
			if current.Len() == 0 {
				return
			}
			args = append(args, current.String())
			current.Reset()
		}
	)

	for _, r := range command {
		switch {
		case escaped:
			current.WriteRune(r)
			escaped = false
		case r == '\\':
			escaped = true
		case quote != 0:
			if r == quote {
				quote = 0
				continue
			}
			current.WriteRune(r)
		case r == '\'' || r == '"':
			quote = r
		case r == ' ' || r == '\t' || r == '\n' || r == '\r':
			flushArg()
		default:
			current.WriteRune(r)
		}
	}

	if escaped {
		return nil, fmt.Errorf("invalid command: trailing escape")
	}

	if quote != 0 {
		return nil, fmt.Errorf("invalid command: unterminated quote")
	}

	flushArg()
	if len(args) == 0 {
		return nil, fmt.Errorf("empty command provided")
	}

	return args, nil
}
