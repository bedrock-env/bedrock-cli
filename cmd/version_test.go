package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

var out io.Writer = os.Stdout

func TestVersion(t *testing.T) {
	out, err := executeCommand(rootCmd, "version")

	if err != nil {
		t.Fatal(err)
	}

	result := strings.TrimSpace(string(out))
	expected := "Bedrock " + VERSION
	if result != expected {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, result)
	}
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}