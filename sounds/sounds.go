package sounds

import (
	. "../constants"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
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
	fade chan bool
}

type Volume struct {
	Streamer beep.Streamer
	Base     float64
	Volume   float64
	Silent   bool
}

func NewSound(nr uint8) Sound {
	s := new(snd)
	s.nr = nr
	s.done = make(chan bool)
	s.quit = make(chan bool)
	s.fade = make(chan bool)
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
	case Alarm1:
		s.path = "soundeffects/The Essential Retro/General Sounds/Alarms/Alarms/sfx_alarm_loop1.ogg"
	case Alarm2:
		s.path = "soundeffects/The Essential Retro/General Sounds/Alarms/Alarms/sfx_alarm_loop2.ogg"
	case Alarm3:
		s.path = "soundeffects/The Essential Retro/General Sounds/Alarms/Alarms/sfx_alarm_loop3.ogg"
	case Alarm4:
		s.path = "soundeffects/The Essential Retro/General Sounds/Alarms/Alarms/sfx_alarm_loop4.ogg"
	case Alarm5:
		s.path = "soundeffects/The Essential Retro/General Sounds/Alarms/Alarms/sfx_alarm_loop5.ogg"
	case Alarm6:
		s.path = "soundeffects/The Essential Retro/General Sounds/Alarms/Alarms/sfx_alarm_loop6.ogg"
	case Alarm7:
		s.path = "soundeffects/The Essential Retro/General Sounds/Alarms/Alarms/sfx_alarm_loop7.ogg"
	case Alarm8:
		s.path = "soundeffects/The Essential Retro/General Sounds/Alarms/Alarms/sfx_alarm_loop8.ogg"
	}
	return s
}

func (s *snd) PlaySound() {

	f, err := os.Open(s.path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	streamer, format, err := vorbis.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	var ctrl *beep.Ctrl
	if s.nr < 100 {
		ctrl = &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
	} else {
		ctrl = &beep.Ctrl{Streamer: beep.Loop(1, streamer), Paused: false}
	}

	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}

	if s.nr < 100 {
		// Musik
		speaker.Play(volume)
	} else {
		// Soundeffekt
		if s.nr == Deathflash {
			speaker.Lock()
			volume.Volume = 2
			speaker.Unlock()
			speaker.Play(beep.Seq(volume, beep.Callback(func() {
				s.done <- true
			})))
		} else {
			// Merkwürdigerweise ist für die Falling-Sounds und Alarm-Sounds ein Resample notwendig
			speaker.Play(beep.Seq(beep.Resample(1, format.SampleRate, format.SampleRate*2, streamer), beep.Callback(func() {
				s.done <- true
			})))
		}
	}

A:
	for {
		select {
		case <-s.quit:
			speaker.Clear()
		case <-s.fade:
			for i := 0; i < 1000; i++ {
				time.Sleep(time.Millisecond * 2)
				speaker.Lock()
				volume.Volume -= 0.01
				speaker.Unlock()
			}
			speaker.Lock()
			volume.Silent = true
			speaker.Unlock()
			break A
		case <-s.done:
			break A
		}
	}
}

func (s *snd) StopSound() {
	s.quit <- true
}

func (s *snd) FadeOut() {
	s.fade <- true
}

func init() {
	sr := beep.SampleRate(44100)
	err := speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		return
	}
}
