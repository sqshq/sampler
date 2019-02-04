package asset

import (
	"fmt"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"io"
	"log"
)

func beep() error {

	bytes, err := Asset("quindar-tone")
	if err != nil {
		log.Fatal("Can't find asset file")
	}

	d, err := mp3.NewDecoder(NewAssetFile(bytes))
	if err != nil {
		return err
	}
	defer d.Close()

	p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		return err
	}

	fmt.Print("\a")

	return nil
}
