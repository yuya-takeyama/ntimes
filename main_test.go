package main

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"testing"
)

func TestNtimes(t *testing.T) {
	stdin := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	stdoutCh := make(chan io.ReadWriter)
	exitCh := make(chan bool)

	go printer(stdout, stdoutCh, exitCh)

	ntimes(3, "echo", []string{"foo", "bar", "baz"}, stdin, stderr, stdoutCh, 1)

	exitCh <- true

	expected := "foo bar baz\nfoo bar baz\nfoo bar baz\n"
	actual := stdout.String()

	if actual != expected {
		t.Error("The result does not match\nExpected:\n" + expected + "\nActual:\n" + actual)
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
