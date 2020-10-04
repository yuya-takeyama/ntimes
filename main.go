package main

import (
	"bytes"
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

	stdoutCh := make(chan io.ReadWriter)
	exitCh := make(chan bool)

	if err != nil {
		panic(err)
	}

	go printer(os.Stdout, stdoutCh, exitCh)

	ntimes(cnt, cmdName, cmdArgs, os.Stdin, os.Stderr, stdoutCh, opts.Parallels)

	exitCh <- true
}

func ntimes(cnt int, cmdName string, cmdArgs []string, stdin io.Reader, stderr io.Writer, stdoutCh chan io.ReadWriter, parallels int) {
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

			stdoutBuffer := new(bytes.Buffer)

			cmd := exec.Command(cmdName, cmdArgs...)
			cmd.Stdin = stdin
			cmd.Stdout = stdoutBuffer
			cmd.Stderr = stderr

			err := cmd.Run()
			if err != nil {
				panic(err)
			}

			stdoutCh <- stdoutBuffer
		}()
	}

	wg.Wait()
}

func printer(stdout io.Writer, stdoutCh chan io.ReadWriter, exitCh chan bool) {
	for {
		select {
		case r := <-stdoutCh:
			_, _ = io.Copy(stdout, r)
		case <-exitCh:
			return
		}
	}
}

func errorf(message string, args ...interface{}) {
	subMessage := fmt.Sprintf(message, args...)
	_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", appName, subMessage)
}
