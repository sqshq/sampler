package console

import (
	"fmt"
	ui "github.com/sqshq/termui"
	"log"
	"time"
)

const (
	RenderRate = 50 * time.Millisecond // TODO not a constant, should be dynamically chosen based on min X scale (per each chart? should be tested). if it is 1 sec, it should be 100 ms, if 2 - 200 ms, if 3 - 300, 4 - 400, 5 - 500 and 500 is max. smth like that.
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
