package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// IntervalSec the maximum number of seconds between output.
var IntervalSec = 9

func main() {
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	lastOutput := time.Now()

	cmd := exec.Command(cmdName, cmdArgs...) // #nosec
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			lastOutput = time.Now()
			fmt.Println(scanner.Text())
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				secondsSinceLastOutput := time.Since(lastOutput) / time.Second
				if int(secondsSinceLastOutput) >= IntervalSec {
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
