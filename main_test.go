package main

import (
	"bytes"
	"io"
	"testing"
)

func TestNtimes(t *testing.T) {
	stdin := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	stdoutCh := make(chan io.ReadWriter)
	exitCh := make(chan bool)

	go printer(stdout, stdoutCh, exitCh)

	ntimes(3, "echo", []string{"foo", "bar", "baz"}, stdin, stdout, stderr, stdoutCh, 1)

	exitCh <- true

	expected := "foo bar baz\nfoo bar baz\nfoo bar baz\n"
	actual := stdout.String()

	if actual != expected {
		t.Error("The result does not match\nExpected:\n" + expected + "\nActual:\n" + actual)
	}
}
