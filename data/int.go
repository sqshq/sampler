package data

import "time"

const (
	startupTimeout  = 200 * time.Millisecond
	minAwaitTimeout = 100 * time.Millisecond
	maxAwaitTimeout = 1 * time.Second
)

type InteractiveShell interface {
	init() error
	execute() (string, error)
}
