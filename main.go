package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

func makeCommand(cmdName string, cmdArgs []string) (*exec.Cmd, io.Reader, error) {
	cmd := exec.Command(cmdName, cmdArgs...) // #nosec
	cmdStdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	cmdStderrReader, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}
	cmdReader := io.MultiReader(cmdStdoutReader, cmdStderrReader)
	return cmd, cmdReader, nil
}

func runCommand(cmd *exec.Cmd) error {
	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	interval, err := time.ParseDuration(os.Args[1])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error parsing interval", err)
		os.Exit(1)
	}
	cmdName := os.Args[2]
	cmdArgs := os.Args[3:]

	lastOutput := time.Now()

	cmd, cmdReader, err := makeCommand(cmdName, cmdArgs)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error starting command", err)
		os.Exit(1)
	}

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

	runCommand(cmd)

	close(quit)
}
