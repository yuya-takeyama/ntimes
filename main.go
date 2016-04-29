package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"

	flags "github.com/jessevdk/go-flags"
)

const AppName = "ntimes"

type Options struct {
	ShowVersion bool `short:"v" long:"version" description:"Show version"`
}

var opts Options

func main() {
	parser := flags.NewParser(&opts, flags.Default^flags.PrintErrors)
	parser.Name = AppName
	parser.Usage = "[OPTIONS] -- COMMAND"

	args, err := parser.Parse()

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}

	if opts.ShowVersion {
		io.WriteString(os.Stdout, fmt.Sprintf("%s v%s, build %s\n", AppName, Version, GitCommit))
		return
	}

	cnt, err := strconv.Atoi(args[0])
	cmdName := args[1]
	cmdArgs := args[2:]

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
