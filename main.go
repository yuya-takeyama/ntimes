package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"

	flags "github.com/jessevdk/go-flags"
)

const appName = "ntimes"

var (
	version   = ""
	gitCommit = ""
)

type options struct {
	Parallels   int  `short:"p" long:"parallels" description:"Parallel degree of execution" default:"1"`
	ShowVersion bool `short:"v" long:"version" description:"Show version"`
}

func main() {
	var opts options
	parser := flags.NewParser(&opts, flags.Default^flags.PrintErrors)
	parser.Name = appName
	parser.Usage = "N [OPTIONS] -- COMMAND"

	args, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			if flagsErr.Type == flags.ErrHelp {
				parser.WriteHelp(os.Stderr)

				return
			}
		}

		errorf("flag parse error: %s", err)
		os.Exit(1)
	}

	if opts.ShowVersion {
		_, _ = io.WriteString(os.Stdout, fmt.Sprintf("%s v%s, build %s\n", appName, version, gitCommit))

		return
	}

	cnt, err := strconv.Atoi(args[0])
	cmdName := args[1]
	cmdArgs := args[2:]

	if err != nil {
		panic(err)
	}

	ntimes(cnt, cmdName, cmdArgs, os.Stdin, os.Stdout, os.Stderr, opts.Parallels)
}

func ntimes(cnt int, cmdName string, cmdArgs []string, stdin io.Reader, stdout io.Writer, stderr io.Writer, parallels int) {
	var wg sync.WaitGroup

	sema := make(chan bool, parallels)

	for i := 0; i < cnt; i++ {
		wg.Add(1)

		go func() {
			sema <- true

			defer func() {
				wg.Done()
				<-sema
			}()

			cmd := exec.Command(cmdName, cmdArgs...)
			cmd.Stdin = stdin
			cmd.Stdout = stdout
			cmd.Stderr = stderr

			err := cmd.Run()
			if err != nil {
				panic(err)
			}
		}()
	}

	wg.Wait()
	close(sema)
}

func errorf(message string, args ...interface{}) {
	subMessage := fmt.Sprintf(message, args...)
	_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", appName, subMessage)
}
