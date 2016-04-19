package main

import (
	"io"
	"os"
	"os/exec"
	"strconv"
)

const AppName = "ntimes"

func main() {
	cnt, err := strconv.Atoi(os.Args[1])
	cmdName := os.Args[2]
	cmdArgs := os.Args[3:]

	if err != nil {
		panic(err)
	}

	ntimes(cnt, cmdName, cmdArgs, os.Stdin, os.Stdout, os.Stderr)
}

func ntimes(cnt int, cmdName string, cmdArgs []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) {
	for i := 0; i < cnt; i++ {
		cmd := exec.Command(cmdName, cmdArgs...)
		cmd.Stdin = stdin
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		err := cmd.Run()

		if err != nil {
			panic(err)
		}
	}
}
