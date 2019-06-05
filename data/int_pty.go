package data

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/kr/pty"
	"github.com/lunixbochs/vtclean"
	"io"
	"os/exec"
	"strings"
	"time"
)

const (
	startupTimeout  = 100 * time.Millisecond
	minAwaitTimeout = 100 * time.Millisecond
	maxAwaitTimeout = 1 * time.Second
)

/**
 * Experimental
 */
type PtyInteractiveShell struct {
	item      *Item
	variables []string
	cmd       *exec.Cmd
	File      io.WriteCloser
	ch        chan string
	errCount  int
}

func (s *PtyInteractiveShell) init() error {

	cmd := exec.Command("sh", "-c", *s.item.initScript)
	enrichEnvVariables(cmd, s.variables)

	file, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	channel := make(chan string)

	go func() {
		for scanner.Scan() {
			channel <- scanner.Text()
		}
	}()

	s.cmd = cmd
	s.File = file
	s.ch = channel

	_, err = file.Read(make([]byte, 4096))
	if err != nil {
		return err
	}

	time.Sleep(startupTimeout)

	return nil
}

func (s *PtyInteractiveShell) execute() (string, error) {

	_, err := io.WriteString(s.File, fmt.Sprintf(" %s\n", s.item.sampleScript))
	if err != nil {
		s.errCount++
		if s.errCount > errorThreshold {
			s.item.ptyShell = nil // restart session
		}
		return "", errors.New(fmt.Sprintf("Failed to execute command: %s", err))
	}

	softTimeout := make(chan bool, 1)
	hardTimeout := make(chan bool, 1)

	go func() {
		time.Sleep(s.getAwaitTimeout() / 2)
		softTimeout <- true
		time.Sleep(s.getAwaitTimeout() * 100)
		hardTimeout <- true
	}()

	var builder strings.Builder
	softTimeoutElapsed := false

await:
	for {
		select {
		case out := <-s.ch:
			cout := vtclean.Clean(out, false)
			if len(cout) > 0 && !strings.Contains(cout, s.item.sampleScript) {
				builder.WriteString(cout)
				builder.WriteString("\n")
				if softTimeoutElapsed {
					break await
				}
			}
		case <-softTimeout:
			if builder.Len() > 0 {
				break await
			} else {
				softTimeoutElapsed = true
			}
		case <-hardTimeout:
			break await
		}
	}

	sample := strings.TrimSpace(builder.String())

	return s.item.transform(sample)
}

func (s *PtyInteractiveShell) getAwaitTimeout() time.Duration {

	timeout := time.Duration(s.item.rateMs) * time.Millisecond

	if timeout > maxAwaitTimeout {
		return maxAwaitTimeout
	} else if timeout < minAwaitTimeout {
		return minAwaitTimeout
	}

	return timeout
}
