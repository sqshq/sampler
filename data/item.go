package data

import (
	"bufio"
	"errors"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/kr/pty"
	"github.com/sqshq/sampler/config"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

const interactiveShellStartupTimeout = 100 * time.Millisecond

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
	Channel chan string
	File    io.WriteCloser
	Cmd     *exec.Cmd
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
			return "", errors.New(fmt.Sprintf("Failed to init interactive shell: %s", err))
		}
	}

	if i.InitScript != nil {
		return i.executeInteractiveShellCmd(variables)
	} else {
		return i.executeCmd(variables, i.SampleScript)
	}
}

func (i *Item) executeCmd(variables []string, script string) (string, error) {

	cmd := exec.Command("sh", "-c", script)
	enrichEnvVariables(cmd, variables)

	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (i *Item) initInteractiveShell(variables []string) error {

	cmd := exec.Command("sh", "-c", *i.InitScript)
	enrichEnvVariables(cmd, variables)

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

	i.InteractiveShell = &InteractiveShell{
		Channel: channel,
		File:    file,
		Cmd:     cmd,
	}

	_, err = file.Read(make([]byte, 4096))
	if err != nil {
		return err
	}

	time.Sleep(interactiveShellStartupTimeout)

	return nil
}

func (i *Item) executeInteractiveShellCmd(variables []string) (string, error) {

	_, err := io.WriteString(i.InteractiveShell.File, fmt.Sprintf(" %s\n", i.SampleScript))
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to execute interactive shell cmd: %s", err))
	}

	timeout := make(chan bool, 1)

	go func() {
		time.Sleep(time.Duration(i.RateMs))
		timeout <- true
	}()

	var outputText strings.Builder

	for {
		select {
		case output := <-i.InteractiveShell.Channel:
			if !strings.Contains(output, i.SampleScript) && len(output) > 0 {
				outputText.WriteString(output)
				outputText.WriteString("\n")
			}
		case <-timeout:
			sample := cleanupOutput(outputText.String())
			return i.transformInteractiveShellCmd(sample)
		}
	}
}

func (i *Item) transformInteractiveShellCmd(sample string) (string, error) {

	if i.TransformScript != nil && len(sample) > 0 {
		return i.executeCmd([]string{"sample=" + sample}, *i.TransformScript)
	}

	return sample, nil
}

func enrichEnvVariables(cmd *exec.Cmd, variables []string) {
	cmd.Env = os.Environ()
	for _, variable := range variables {
		cmd.Env = append(cmd.Env, variable)
	}
}

func cleanupOutput(output string) string {
	s := strings.TrimSpace(output)
	if idx := strings.Index(s, "\r"); idx != -1 {
		return s[idx+1:]
	}
	return s
}
