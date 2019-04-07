package data

import (
	"bufio"
	"errors"
	ui "github.com/gizak/termui/v3"
	"github.com/sqshq/sampler/config"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Item struct {
	Label            string
	SampleScript     string
	InitScript       *string
	TransformScript  *string
	Color            *ui.Color
	RateMs           int
	InteractiveShell *InteractiveShell
}

type InteractiveShell struct {
	StdoutCh chan string
	StderrCh chan string
	Stdin    io.WriteCloser
	Cmd      *exec.Cmd
}

func NewItems(cfgs []config.Item, rateMs int) []*Item {

	items := make([]*Item, 0)

	for _, i := range cfgs {
		item := &Item{
			Label:           *i.Label,
			SampleScript:    *i.SampleScript,
			InitScript:      i.InitScript,
			TransformScript: i.TransformScript,
			Color:           i.Color,
			RateMs:          rateMs,
		}
		items = append(items, item)
	}

	return items
}

func (i *Item) nextValue(variables []string) (string, error) {

	if i.InitScript != nil && i.InteractiveShell == nil {
		err := i.initInteractiveShell(variables)
		if err != nil {
			return "", err
		}
	}

	if i.InitScript != nil {
		return i.executeInteractiveShellCmd(variables)
	} else {
		return i.executeCmd(variables)
	}
}

func (i *Item) executeCmd(variables []string) (string, error) {

	cmd := exec.Command("sh", "-c", i.SampleScript)
	enrichEnvVariables(cmd, variables)

	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func (i *Item) initInteractiveShell(variables []string) error {

	cmd := exec.Command("sh", "-c", *i.InitScript)
	enrichEnvVariables(cmd, variables)

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

	i.InteractiveShell = &InteractiveShell{
		StdoutCh: stdoutCh,
		StderrCh: stderrCh,
		Stdin:    stdin,
		Cmd:      cmd,
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (i *Item) executeInteractiveShellCmd(variables []string) (string, error) {

	_, err := io.WriteString(i.InteractiveShell.Stdin, i.SampleScript+"\n")
	if err != nil {
		return "", err
	}

	timeout := make(chan bool, 1)

	go func() {
		time.Sleep(time.Duration(i.RateMs))
		timeout <- true
	}()

	var resultText strings.Builder
	var errorText strings.Builder

	for {
		select {
		case stdout := <-i.InteractiveShell.StdoutCh:
			if len(stdout) > 0 {
				resultText.WriteString(stdout)
				resultText.WriteString("\n")
			}
		case stderr := <-i.InteractiveShell.StderrCh:
			if len(stderr) > 0 {
				errorText.WriteString(stderr)
				errorText.WriteString("\n")
			}
		case <-timeout:
			if errorText.Len() > 0 {
				return "", errors.New(errorText.String())
			} else {
				return i.transformInteractiveShellCmd(resultText.String())
			}
		}
	}
}

func (i *Item) transformInteractiveShellCmd(value string) (string, error) {
	// TODO use transform script, if any
	return strings.TrimSpace(value), nil
}

func enrichEnvVariables(cmd *exec.Cmd, variables []string) {
	cmd.Env = os.Environ()
	for _, variable := range variables {
		cmd.Env = append(cmd.Env, variable)
	}
}
