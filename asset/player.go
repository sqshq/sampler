package asset

import (
	"log"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

type AudioPlayer struct {
	context *oto.Context
	beep    []byte
}

func NewAudioPlayer() *AudioPlayer {

	bytes, err := Asset("quindar-tone.mp3")
	if err != nil {
		log.Fatal("Failed to find audio file")
	}

	// Trigger audio is optional, so keep the application usable without a device.
	context, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
	})
	if err != nil {
		return nil
	}
	<-ready
	if context.Err() != nil {
		return nil
	}

	return &AudioPlayer{
		context: context,
		beep:    bytes,
	}
}

func (a *AudioPlayer) Beep() {

	decoder, err := mp3.NewDecoder(NewAssetFile(a.beep))
	if err != nil {
		panic(err)
	}

	player := a.context.NewPlayer(decoder)
	player.Play()
	for player.IsPlaying() {
		if err := player.Err(); err != nil {
			panic(err)
		}
		time.Sleep(time.Millisecond)
	}
	if err := player.Err(); err != nil {
		panic(err)
	}
}
