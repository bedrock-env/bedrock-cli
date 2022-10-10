package helpers

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type execInShellContext = func(name string, arg ...string) *exec.Cmd

func ExecuteCommand(binary string, command string) (string, error) {
	return ExecuteCommandWithArgs(binary, strings.Split(command, " "))
}

func ExecuteCommandWithArgs(binary string, command []string) (string, error) {
	out, err := exec.Command(binary, command...).CombinedOutput()

	return string(out), err
}

func ExecuteCommandInShell(execCtx execInShellContext, shell string, command string) (string, error) {
	cmd := execCtx(shell, "-c", command)
	out, err := cmd.CombinedOutput()

	return strings.TrimSpace(string(out)), err
}

func PromptYN(message string) (result bool) {
	fmt.Printf("%s%s (y/n) %s", ColorYellow, message, ColorReset)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')

	return strings.TrimSpace(response) == "y"
}
