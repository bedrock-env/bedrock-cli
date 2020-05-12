package cmd

import (
	"io"
	"os"
	"strings"
	"testing"
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
