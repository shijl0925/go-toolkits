package systemx_test

import (
	"github.com/shijl0925/go-toolkits/systemx"
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

func TestStartProcess(t *testing.T) {
	pid, err := systemx.StartProcess("sleep", "10")
	if err != nil {
		t.Errorf("StartProcess error: %v", err)
	}
	if pid <= 0 {
		t.Errorf("Invalid pid: %d", pid)
	}
}
