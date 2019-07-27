//+build !windows

package data

import (
	"bufio"
	"fmt"
	"github.com/kr/pty"
	"github.com/lunixbochs/vtclean"
	"io"
	"os/exec"
	"strings"
	"time"
)

// PtyInteractiveShell represents PTY interactive shell sampling metadata
type PtyInteractiveShell struct {
	item      *Item
	variables []string
	cmd       *exec.Cmd
	file      io.WriteCloser
	ch        chan string
	errCount  int
	timeout   time.Duration
}

func (s *PtyInteractiveShell) init() error {

	cmd := exec.Command("sh", "-c", s.item.initScripts[0])
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
	s.file = file
	s.ch = channel

	_, err = file.Read(make([]byte, 4096))
	if err != nil {
		return err
	}

	time.Sleep(startupTimeout)

	for i := 1; i < len(s.item.initScripts); i++ {
		_, err = io.WriteString(s.file, fmt.Sprintf(" %s\n", s.item.initScripts[i]))
		if err != nil {
			return err
		}
		time.Sleep(startupTimeout) // TODO wait until cmd complete
	}

	return nil
}

func (s *PtyInteractiveShell) execute() (string, error) {

	_, err := io.WriteString(s.file, fmt.Sprintf(" %s\n", s.item.sampleScript))
	if err != nil {
		s.errCount++
		if s.errCount > errorThreshold {
			_ = s.cmd.Wait()
			_ = s.file.Close()
			s.item.ptyShell = nil // restart session
		}
		return "", fmt.Errorf("failed to execute command: %s", err)
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

	if s.timeout > maxAwaitTimeout {
		return maxAwaitTimeout
	} else if s.timeout < minAwaitTimeout {
		return minAwaitTimeout
	}

	return s.timeout
}
