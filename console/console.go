package console

import (
	"fmt"
	ui "github.com/sqshq/termui"
	"log"
	"time"
)

const (
	RenderRate = 25 * time.Millisecond
	Title      = "sampler"
)

type Console struct{}

func (self *Console) Init() {

	fmt.Printf("\033]0;%s\007", Title)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
}

func (self *Console) Close() {
	ui.Close()
}
