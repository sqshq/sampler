package console

import (
	"fmt"
	ui "github.com/sqshq/termui"
	"log"
	"time"
)

const (
	MaxRenderInterval = 1000 * time.Millisecond
	MinRenderInterval = 100 * time.Millisecond
	AppTitle          = "sampler"
	AppVersion        = "0.1.0"
)

type Console struct{}

func (self *Console) Init() {

	fmt.Printf("\033]0;%s\007", AppTitle)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
}

func (self *Console) Close() {
	ui.Close()
}
