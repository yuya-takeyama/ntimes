package main

import (
	"bytes"
	"errors"
	"os/exec"
	"testing"
)

func TestSerial(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "2", "--", "sh", "-c", `echo "Hi!"; sleep 1; echo "Bye"`)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	expected := "Hi!\nBye\nHi!\nBye\n"
	if err := cmd.Run(); err != nil {
		t.Errorf("failed: %v", err)
	}

	if stdout.String() != expected {
		t.Errorf("stdout doesn't match\nExpected:\n%s\nActual:\n%s", expected, stdout.String())
	}
}
func TestParallel(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-p", "2", "2", "--", "sh", "-c", `echo "Hi!"; sleep 1; echo "Bye"`)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	expected := "Hi!\nHi!\nBye\nBye\n"
	if err := cmd.Run(); err != nil {
		t.Errorf("failed: %v", err)
	}

	if stdout.String() != expected {
		t.Errorf("stdout doesn't match\nExpected:\n%s\nActual:\n%s", expected, stdout.String())
	}
}

func TestVersion(t *testing.T) {
	cmd := exec.Command("go", "run", "-ldflags", "-X main.version=1.2.3 -X main.gitCommit=deadbeef", "main.go", "--version")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	expected := "ntimes v1.2.3, build deadbeef\n"
	if err := cmd.Run(); err != nil {
		t.Errorf("failed: %v", err)
	}

	if stdout.String() != expected {
		t.Errorf("stdout doesn't match\nExpected:\n%s\nActual:\n%s", expected, stdout.String())
	}
}

func TestHelp(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "--help")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	expected := `Usage:
  ntimes N [OPTIONS] -- COMMAND

Application Options:
  -p, --parallels= Parallel degree of execution (default: 1)
  -v, --version    Show version

Help Options:
  -h, --help       Show this help message
`
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Errorf("failed: %v", err)
		}
	}

	if stdout.String() != "" {
		t.Errorf("stdout should be empty")
	}

	if stderr.String() != expected {
		t.Errorf("stderr doesn't match\nExpected: \n%s\nActual:\n%s", expected, stderr.String())
	}
}

func TestUnknownOption(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "--foo")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	expected := "ntimes: flag parse error: unknown flag `foo'\nexit status 1\n"
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Errorf("failed: %v", err)
		}
	}

	if stdout.String() != "" {
		t.Errorf("stdout should be empty")
	}

	if stderr.String() != expected {
		t.Errorf("stderr doesn't match\nExpected: \n%s\nActual:\n%s", expected, stderr.String())
	}
}
