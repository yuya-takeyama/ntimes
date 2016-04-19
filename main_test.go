package main

import (
	"bytes"
	"testing"
)

func TestNtimes(t *testing.T) {
	stdin := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	ntimes(3, "echo", []string{"foo", "bar", "baz"}, stdin, stdout, stderr)

	expected := "foo bar baz\nfoo bar baz\nfoo bar baz\n"
	actual := stdout.String()

	if actual != expected {
		t.Error("The result does not match\nExpected:\n" + expected + "\nActual:\n" + actual)
	}
}
