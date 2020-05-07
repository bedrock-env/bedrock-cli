package helpers

import (
	"os/exec"
	"strings"
)

func ExecuteCommand(binary string, command string) (string, error) {
	return ExecuteCommandWithArgs(binary, strings.Split(command, " "))
}

func ExecuteCommandWithArgs(binary string, command []string) (string, error) {
	out, err := exec.Command(binary, command...).CombinedOutput()

	return string(out), err
}

func ExecuteCommandInShell(shell string, command string) (string, error) {
	out, err := exec.Command(shell, "-c", command).CombinedOutput()

	return strings.TrimSpace(string(out)), err
}
