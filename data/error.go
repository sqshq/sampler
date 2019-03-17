package data

import (
	"fmt"
	"os/exec"
)

func getErrorMessage(err *exec.ExitError) string {
	stderr := string(err.Stderr)
	if len(stderr) == 0 {
		return err.Error()
	} else {
		return fmt.Sprintf("%.200s", stderr)
	}
}
