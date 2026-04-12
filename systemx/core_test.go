package systemx_test

import (
	"errors"
	"fmt"
	"github.com/shijl0925/go-toolkits/systemx"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestExecCommand(t *testing.T) {
	t.Helper()

	expectedDir := t.TempDir()
	stdout, stderr, err := systemx.ExecCommand(helperCommand(t, "pwd"), func(cmd *exec.Cmd) {
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
		cmd.Dir = expectedDir
	})
	if err != nil {
		t.Errorf("ExecCommand error: %v", err)
	}
	if strings.TrimSpace(stdout) != expectedDir {
		t.Fatalf("expected stdout %q, got %q", expectedDir, stdout)
	}
	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}

	// error command
	stdout, stderr, err = systemx.ExecCommand(helperCommand(t, "exit1"), helperEnvOption())
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "helper error") {
		t.Fatalf("expected helper stderr, got %q", stderr)
	}
}

func TestExecCommandQuotedArgs(t *testing.T) {
	stdout, stderr, err := systemx.ExecCommand(helperCommand(t, "printargs", "hello world"), helperEnvOption())
	if err != nil {
		t.Fatalf("ExecCommand error: %v", err)
	}

	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}

	if stdout != "hello world" {
		t.Fatalf("expected quoted argument output, got %q", stdout)
	}
}

func TestExecCommandEscapedQuotesAndSpaces(t *testing.T) {
	stdout, stderr, err := systemx.ExecCommand(helperCommand(t, "printargs", `hello "quoted" world @#$`), helperEnvOption())
	if err != nil {
		t.Fatalf("ExecCommand error: %v", err)
	}

	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}

	if stdout != `hello "quoted" world @#$` {
		t.Fatalf("expected escaped quote output, got %q", stdout)
	}
}

func TestExecCommandInvalidSyntax(t *testing.T) {
	_, _, err := systemx.ExecCommand(`printf "unterminated`)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
	if !strings.Contains(err.Error(), "unterminated quote") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStartProcess(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable error: %v", err)
	}
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")

	pid, err := systemx.StartProcess(executable, "-test.run=TestHelperProcess", "--", "sleep", "10s")
	if err != nil {
		t.Errorf("StartProcess error: %v", err)
	}
	if pid <= 0 {
		t.Errorf("Invalid pid: %d", pid)
	}

	if err := systemx.KillProcess(pid); err != nil && !errors.Is(err, os.ErrProcessDone) {
		t.Fatalf("KillProcess error: %v", err)
	}
}

func TestStartProcessRejectsRelativePath(t *testing.T) {
	relativeCommand := "." + string(os.PathSeparator) + "sleep"
	if _, err := systemx.StartProcess(relativeCommand, "1"); err == nil {
		t.Fatal("expected error for relative command path")
	} else if !strings.Contains(err.Error(), "relative command paths are not allowed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHelperProcess(t *testing.T) {
	t.Helper()

	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	separatorIndex := -1
	for i, arg := range os.Args {
		if arg == "--" {
			separatorIndex = i
			break
		}
	}

	if separatorIndex == -1 || separatorIndex+1 >= len(os.Args) {
		fmt.Fprint(os.Stderr, "missing helper args")
		os.Exit(2)
	}

	action := os.Args[separatorIndex+1]
	args := os.Args[separatorIndex+2:]

	switch action {
	case "pwd":
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Fprint(os.Stdout, wd)
	case "printargs":
		fmt.Fprint(os.Stdout, strings.Join(args, " "))
	case "exit1":
		fmt.Fprint(os.Stderr, "helper error")
		os.Exit(1)
	case "sleep":
		if len(args) != 1 {
			fmt.Fprint(os.Stderr, "missing duration")
			os.Exit(2)
		}
		duration, err := time.ParseDuration(args[0])
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(2)
		}
		time.Sleep(duration)
	default:
		fmt.Fprintf(os.Stderr, "unknown helper action: %s", action)
		os.Exit(2)
	}

	os.Exit(0)
}

func helperCommand(t *testing.T, args ...string) string {
	t.Helper()

	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable error: %v", err)
	}

	commandArgs := append([]string{executable, "-test.run=TestHelperProcess", "--"}, args...)
	command := make([]string, 0, len(commandArgs))
	for _, arg := range commandArgs {
		command = append(command, quoteCommandArg(arg))
	}

	return strings.Join(command, " ")
}

func quoteCommandArg(arg string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `"`, `\"`)
	return `"` + replacer.Replace(arg) + `"`
}

func helperEnvOption() func(cmd *exec.Cmd) {
	return func(cmd *exec.Cmd) {
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	}
}
