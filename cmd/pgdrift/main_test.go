package main

import (
	"os/exec"
	"strings"
	"testing"
)

// TestMainNoArgs verifies that the binary exits with a non-zero code and
// prints a usage hint when required flags are missing.
func TestMainNoArgs(t *testing.T) {
	cmd := exec.Command("go", "run", ".", )
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit code when no flags provided")
	}
	output := string(out)
	if !strings.Contains(output, "--source") && !strings.Contains(output, "-source") {
		t.Errorf("expected usage hint in output, got: %s", output)
	}
}

// TestMainMissingTarget verifies that omitting --target also produces an error.
func TestMainMissingTarget(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "--source", "postgres://localhost/db")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit code when --target is missing")
	}
	output := string(out)
	if !strings.Contains(output, "--source and --target are required") {
		t.Errorf("unexpected error output: %s", output)
	}
}
