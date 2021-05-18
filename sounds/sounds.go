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

type snd struct {
	nr   uint8
	path string
	done chan bool
	quit chan bool
}

func NewSound(nr uint8) Sound {
	s := new(snd)
	s.nr = nr
	s.done = make(chan bool)
	s.quit = make(chan bool)
	switch nr {
	// Musik:
	case ThroughSpace:
		s.path = "soundeffects/through_space.ogg"
	case TheFieldOfDreams:
		s.path = "soundeffects/the_field_of_dreams.ogg"
	case OrbitalColossus:
		s.path = "soundeffects/orbital_colossus.ogg"
	case Fight:
		s.path = "soundeffects/fight.ogg"
	case JuhaniJunkalaTitle:
		s.path = "soundeffects/5 Action Chiptunes By Juhani Junkala/Juhani Junkala [Retro Game Music Pack] Title Screen.ogg"
	case JuhaniJunkalaLevel1:
		s.path = "soundeffects/5 Action Chiptunes By Juhani Junkala/Juhani Junkala [Retro Game Music Pack] Level 1.ogg"
	case JuhaniJunkalaLevel2:
		s.path = "soundeffects/5 Action Chiptunes By Juhani Junkala/Juhani Junkala [Retro Game Music Pack] Level 2.ogg"
	case JuhaniJunkalaLevel3:
		s.path = "soundeffects/5 Action Chiptunes By Juhani Junkala/Juhani Junkala [Retro Game Music Pack] Level 3.ogg"
	case JuhaniJunkalaEnd:
		s.path = "soundeffects/5 Action Chiptunes By Juhani Junkala/Juhani Junkala [Retro Game Music Pack] Ending.ogg"
	case ObservingTheStar:
		s.path = "soundeffects/ObservingTheStar/ObservingTheStar.ogg"
	case MyVeryOwnDeadShip:
		s.path = "soundeffects/projects/MyVeryOwnDeadShip.ogg"

	// Soundeffekte:
	case Deathflash:
		s.path = "soundeffects/deathflash.ogg"
	case Falling1:
		s.path = "soundeffects/The Essential Retro/Movement/Falling Sounds/sfx_sounds_falling1.ogg"
	case Falling2:
		s.path = "soundeffects/The Essential Retro/Movement/Falling Sounds/sfx_sounds_falling2.ogg"
	case Falling3:
		s.path = "soundeffects/The Essential Retro/Movement/Falling Sounds/sfx_sounds_falling3.ogg"
	case Falling4:
		s.path = "soundeffects/The Essential Retro/Movement/Falling Sounds/sfx_sounds_falling4.ogg"
	case Falling5:
		s.path = "soundeffects/The Essential Retro/Movement/Falling Sounds/sfx_sounds_falling5.ogg"
	case Falling6:
		s.path = "soundeffects/The Essential Retro/Movement/Falling Sounds/sfx_sounds_falling6.ogg"
	case Falling7:
		s.path = "soundeffects/The Essential Retro/Movement/Falling Sounds/sfx_sounds_falling7.ogg"
	case Falling8:
		s.path = "soundeffects/The Essential Retro/Movement/Falling Sounds/sfx_sounds_falling8.ogg"
	case Falling9:
		s.path = "soundeffects/The Essential Retro/Movement/Falling Sounds/sfx_sounds_falling9.ogg"
	case Falling10:
		s.path = "soundeffects/The Essential Retro/Movement/Falling Sounds/sfx_sounds_falling10.ogg"
	case Falling11:
		s.path = "soundeffects/The Essential Retro/Movement/Falling Sounds/sfx_sounds_falling11.ogg"
	case Falling12:
		s.path = "soundeffects/The Essential Retro/Movement/Falling Sounds/sfx_sounds_falling12.ogg"

	}
	return s
}

func (s *snd) PlaySound() {

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

func (s *snd) StopSound() {
	s.quit <- true
}

func init() {
	sr := beep.SampleRate(44100)
	err := speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		return
	}
}
