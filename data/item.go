package data

import (
	"bufio"
	"errors"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/kr/pty"
	"github.com/lunixbochs/vtclean"
	"github.com/sqshq/sampler/config"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	interactiveShellStartupTimeout  = 100 * time.Millisecond
	interactiveShellMinAwaitTimeout = 100 * time.Millisecond
	interactiveShellMaxAwaitTimeout = 1 * time.Second
	interactiveShellErrorThreshold  = 10
)

type Item struct {
	label            string
	sampleScript     string
	initScript       *string
	transformScript  *string
	color            *ui.Color
	rateMs           int
	errorsCount      int
	interactiveShell *InteractiveShell
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
			label:           *i.Label,
			sampleScript:    *i.SampleScript,
			initScript:      i.InitScript,
			transformScript: i.TransformScript,
			color:           i.Color,
			rateMs:          rateMs,
		}
		items = append(items, item)
	}

	return items
}

func (i *Item) nextValue(variables []string) (string, error) {

	if i.initScript != nil && i.interactiveShell == nil {
		err := i.initInteractiveShell(variables)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Failed to init interactive shell: %s", err))
		}
	}

	if i.initScript != nil {
		return i.executeInteractiveShellCmd(variables)
	} else {
		return i.executeCmd(variables, i.sampleScript)
	}
}

func (i *Item) executeCmd(variables []string, script string) (string, error) {

	cmd := exec.Command("sh", "-c", script)
	enrichEnvVariables(cmd, variables)

	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	result := vtclean.Clean(string(output), false)

	return result, nil
}

func (i *Item) initInteractiveShell(variables []string) error {

	cmd := exec.Command("sh", "-c", *i.initScript)
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

	i.interactiveShell = &InteractiveShell{
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

	_, err := io.WriteString(i.interactiveShell.File, fmt.Sprintf(" %s\n", i.sampleScript))
	if err != nil {
		i.errorsCount++
		if i.errorsCount > interactiveShellErrorThreshold {
			i.interactiveShell = nil // restart session
			i.errorsCount = 0
		}
		return "", errors.New(fmt.Sprintf("Failed to execute interactive shell cmd: %s", err))
	}

	softTimeout := make(chan bool, 1)
	hardTimeout := make(chan bool, 1)

	go func() {
		time.Sleep(i.getAwaitTimeout() / 4)
		softTimeout <- true
		time.Sleep(i.getAwaitTimeout() * 100)
		hardTimeout <- true
	}()

	var builder strings.Builder
	softTimeoutElapsed := false

await:
	for {
		select {
		case output := <-i.interactiveShell.Channel:
			o := vtclean.Clean(output, false)
			if len(o) > 0 && !strings.Contains(o, i.sampleScript) {
				builder.WriteString(o)
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

	return i.transformInteractiveShellCmd(sample)
}

func (i *Item) transformInteractiveShellCmd(sample string) (string, error) {

	if i.transformScript != nil && len(sample) > 0 {
		return i.executeCmd([]string{"sample=" + sample}, *i.transformScript)
	}

	return sample, nil
}

func (i *Item) getAwaitTimeout() time.Duration {

	timeout := time.Duration(i.rateMs) * time.Millisecond

	if timeout > interactiveShellMaxAwaitTimeout {
		return interactiveShellMaxAwaitTimeout
	} else if timeout < interactiveShellMinAwaitTimeout {
		return interactiveShellMinAwaitTimeout
	}

	return timeout
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
