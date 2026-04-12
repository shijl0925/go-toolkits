package systemx_test

import (
	"github.com/shijl0925/go-toolkits/systemx"
	"os"
	"os/exec"
	"testing"
)

func TestExecCommand(t *testing.T) {
	// linux or mac
	stdout, stderr, err := systemx.ExecCommand("ls", func(cmd *exec.Cmd) {
		cmd.Dir = "/"
	})
	if err != nil {
		t.Errorf("ExecCommand error: %v", err)
	}
	t.Logf("stdout: %s", stdout)
	t.Logf("stderr: %s", stderr)

	// error command
	stdout, stderr, err = systemx.ExecCommand("abc")
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
	t.Logf("stdout: %s", stdout)
	t.Logf("stderr: %s", stderr)
}

func TestExecCommandQuotedArgs(t *testing.T) {
	stdout, stderr, err := systemx.ExecCommand(`printf "%s" "hello world"`)
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

func TestExecCommandInvalidSyntax(t *testing.T) {
	_, _, err := systemx.ExecCommand(`printf "unterminated`)
	if err == nil {
		t.Fatal("expected parse error, got nil")
	}
}

func TestStartProcess(t *testing.T) {
	pid, err := systemx.StartProcess("sleep", "10")
	if err != nil {
		t.Errorf("StartProcess error: %v", err)
	}
	if pid <= 0 {
		t.Errorf("Invalid pid: %d", pid)
	}

	if err := systemx.KillProcess(pid); err != nil {
		t.Fatalf("KillProcess error: %v", err)
	}
}

func TestStartProcessRejectsRelativePath(t *testing.T) {
	relativeCommand := "." + string(os.PathSeparator) + "sleep"
	if _, err := systemx.StartProcess(relativeCommand, "1"); err == nil {
		t.Fatal("expected error for relative command path")
	}
}
