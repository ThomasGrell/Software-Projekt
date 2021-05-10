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
	path string
	done chan bool
	quit chan bool
}

func NewSound(nr uint8) *sound {
	s := new(sound)
	s.nr = nr
	s.done = make(chan bool)
	s.quit = make(chan bool)
	switch nr {
	// Musik:
	case ThroughSpace:
		s.path = "soundeffects/through_space.ogg"

	// Soundeffekte:
	case Deathflash:
		s.path = "soundeffects/deathflash.ogg"
	}
	return s
}

func (s *sound) PlaySound() {

	f, err := os.Open(s.path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	streamer, _, err := vorbis.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		s.done <- true
	})))

A:
	for {
		select {
		case <-s.quit:
			speaker.Clear()
		case <-s.done:
			// Falls es Musik ist dann wird diese in einer Endlosschleife wiederholt.
			if s.nr < 100 {
				streamer.Seek(0)
				speaker.Play(beep.Seq(streamer, beep.Callback(func() {
					s.done <- true
				})))
			} else {
				break A
			}
		}
	}
}

func (s *sound) StopSound() {
	s.quit <- true
}

func init() {
	sr := beep.SampleRate(44100)
	err := speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		return
	}
}
