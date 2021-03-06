package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/bedrock-env/bedrock-cli/bedrock"
)

var out io.Writer = os.Stdout

func TestVersion(t *testing.T) {
	out, err := executeCommand(rootCmd, "version")

	if err != nil {
		t.Fatal(err)
	}

	expected := fmt.Sprintf("Bedrock CLI %s\nBedrock Core %s", bedrock.VERSION, bedrock.CoreVersion())
	result := strings.TrimSpace(out)

	if result != expected {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, result)
	}
}
