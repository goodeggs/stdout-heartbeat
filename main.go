package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func main() {
	interval, err := time.ParseDuration(os.Args[1])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error parsing interval", err)
		os.Exit(1)
	}
	cmdName := os.Args[2]
	cmdArgs := os.Args[3:]

	lastOutput := time.Now()

	cmd := exec.Command(cmdName, cmdArgs...) // #nosec
	cmdStdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}
	cmdStderrReader, err := cmd.StderrPipe()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error creating StderrPipe for Cmd", err)
		os.Exit(1)
	}
	cmdReader := io.MultiReader(cmdStdoutReader, cmdStderrReader)
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			lastOutput = time.Now()
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if time.Since(lastOutput) >= interval {
					fmt.Println("♥")
					lastOutput = time.Now()
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}

	close(quit)
}
