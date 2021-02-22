package sounds

import (
	. "../constants"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"log"
	"os"
	"time"
)

type sound struct {
	nr   uint8
	done chan bool
	quit chan bool
}

func NewSound(nr uint8) *sound {
	s := new(sound)
	s.done = make(chan bool)
	s.quit = make(chan bool)
	s.nr = nr
	return s
}

func (s *sound) PlaySound() {

	var path string

	switch s.nr {
	// Musik:
	case ThroughSpace:
		path = "soundeffects/through_space.ogg"

	// Soundeffekte:
	case Deathflash:
		path = "soundeffects/deathflash.ogg"
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := vorbis.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return
	}

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		s.done <- true
	})))

	select {
	case <-s.quit:
		speaker.Clear()
	case <-s.done:
		speaker.Clear()
	}
}

func (s *sound) StopMusic() {
	s.quit <- true
}
