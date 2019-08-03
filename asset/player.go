package asset

import (
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"io"
	"log"
)

type AudioPlayer struct {
	player *oto.Player
	beep   []byte
}

func NewAudioPlayer() *AudioPlayer {

	bytes, err := Asset("quindar-tone.mp3")
	if err != nil {
		log.Fatal("Failed to find audio file")
	}

	player, err := oto.NewPlayer(44100, 2, 2, 8192)
	if err != nil {
		// it is expected to fail when some of the system
		// libraries are not available (e.g. libasound2)
		// it is not the main functionality of the application,
		// so we allow startup in no-sound mode
		return nil
	}

	return &AudioPlayer{
		player: player,
		beep:   bytes,
	}
}

func (a *AudioPlayer) Beep() {

	decoder, err := mp3.NewDecoder(NewAssetFile(a.beep))
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(a.player, decoder); err != nil {
		panic(err)
	}
}

func (a *AudioPlayer) Close() {
	_ = a.player.Close()
}
