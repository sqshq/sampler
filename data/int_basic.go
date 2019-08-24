package data

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/lunixbochs/vtclean"
	"io"
	"os/exec"
	"strings"
	"time"
)

// BasicInteractiveShell represents non-PTY interactive shell sampling metadata
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

	cmd := exec.Command("sh", "-c", s.item.initScripts[0])
	enrichEnvVariables(cmd, s.variables)

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

	for i := 1; i < len(s.item.initScripts); i++ {
		_, err := io.WriteString(s.stdin, fmt.Sprintf(" %s\n", s.item.initScripts[i]))
		if err != nil {
			return err
		}
		time.Sleep(startupTimeout) // TODO wait until cmd complete
	}

	return nil
}

func (s *BasicInteractiveShell) execute() (string, error) {

	if s.stdin == nil {
		return "", nil
	}

	_, err := io.WriteString(s.stdin, fmt.Sprintf(" %s\n", s.item.sampleScript))
	if err != nil {
		s.errCount++
		if s.errCount > errorThreshold {
			_ = s.cmd.Wait()
			s.item.basicShell = nil // restart session
		}
		return "", fmt.Errorf("failed to execute command: %s", err)
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
			}
			return s.item.transform(vtclean.Clean(resultText.String(), false))
		}
	}
}
