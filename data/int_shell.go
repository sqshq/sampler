package data

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

type BasicInteractiveShell struct {
	item      *Item
	variables []string
	stdoutCh  chan string
	stderrCh  chan string
	stdin     io.WriteCloser
	cmd       *exec.Cmd
	errCount  int
	timeout   time.Duration
}

func (s *BasicInteractiveShell) init() error {

	cmd := exec.Command("sh", "-c", *s.item.initScript)
	enrichEnvVariables(cmd, s.variables)
	cmd.Wait()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	stdoutCh := make(chan string)
	stderrCh := make(chan string)

	go func() {
		for stdoutScanner.Scan() {
			stdoutCh <- stdoutScanner.Text()
			stderrCh <- stderrScanner.Text()
		}
	}()

	s.stdoutCh = stdoutCh
	s.stderrCh = stderrCh
	s.stdin = stdin
	s.cmd = cmd

	err = cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (s *BasicInteractiveShell) execute() (string, error) {

	_, err := io.WriteString(s.stdin, s.item.sampleScript+"\n")
	if err != nil {
		s.errCount++
		if s.errCount > errorThreshold {
			_ = s.cmd.Wait()
			s.item.basicShell = nil // restart session
		}
		return "", errors.New(fmt.Sprintf("Failed to execute command: %s", err))
	}

	timeout := make(chan bool, 1)

	go func() {
		time.Sleep(s.timeout)
		timeout <- true
	}()

	var resultText strings.Builder
	var errorText strings.Builder

	for {
		select {
		case stdout := <-s.stdoutCh:
			if len(stdout) > 0 {
				resultText.WriteString(stdout)
				resultText.WriteString("\n")
			}
		case stderr := <-s.stderrCh:
			if len(stderr) > 0 {
				errorText.WriteString(stderr)
				errorText.WriteString("\n")
			}
		case <-timeout:
			if errorText.Len() > 0 {
				return "", errors.New(errorText.String())
			} else {
				return resultText.String(), nil
			}
		}
	}
}
