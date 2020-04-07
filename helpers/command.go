package helpers

import (
	"os/exec"
	"strings"
)

func ExecuteCommand(binary string, command string) (string, error) {
	args := strings.Split(command, " ")
	out, err := exec.Command(binary, args...).CombinedOutput()

	return string(out), err
}

func ExecuteInShell(shell string, command string) (string, error) {
	out, err := exec.Command(shell, "-c", command).CombinedOutput()

	return strings.TrimSpace(string(out)), err
}