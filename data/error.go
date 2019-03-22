package data

import (
	"fmt"
	"os/exec"
)

func getErrorMessage(err error) string {

	exitErr, ok := err.(*exec.ExitError)
	message := err.Error()

	if ok {
		stderr := string(exitErr.Stderr)
		if len(stderr) != 0 {
			message = fmt.Sprintf("%.200s", stderr)
		}
	}

	return message
}
