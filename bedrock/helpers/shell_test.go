package helpers

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
)

const (
	testStdoutValue = "testing"
)

func TestShellProcessSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stdout, testStdoutValue)
	os.Exit(0)
}

func fakeExecCommandSuccess(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestShellProcessSuccess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_TEST_PROCESS=1"}
	return cmd
}

func TestNewExecuteCommandInShell(t *testing.T) {
	out, _ := ExecuteCommandInShell(fakeExecCommandSuccess, "echo", "'testing'")
	if out != "testing" {
		log.Fatalf("Expected: %s\nActual: %s", "testing", out)
	}
}
