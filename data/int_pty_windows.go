package data

import (
	"errors"
	"time"
)

type PtyInteractiveShell struct {
	item      *Item
	variables []string
	timeout   time.Duration
}

func (s *PtyInteractiveShell) init() error {
	return errors.New("PTY mode is not supported on Windows")
}

func (s *PtyInteractiveShell) execute() (string, error) {
	return "", errors.New("PTY mode is not supported on Windows")
}
