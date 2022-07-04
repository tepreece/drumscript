package main

import (
	"errors"
	"gitlab.com/gomidi/midi/v2"
	"time"
)

var (
	send                        func(msg midi.Message) error
	PortOpen                    bool
	Tempo                       uint16
	ActiveSong, ActiveSongEvent int
	ActivePatterns              []int
	ticker                      *time.Ticker
	done                        = make(chan bool)
	stopped                     = make(chan bool)
	eventLoopRunning            bool
)

const (
	DefaultTempo uint16 = 120
)

func setPort(fields []string) error {
	if len(fields) != 2 {
		return errors.New("Invalid port name")
	}

	out, err := midi.FindOutPort(fields[1])
	if err != nil {
		return err
	}

	send, err = midi.SendTo(out)
	if err != nil {
		return err
	}

	PortOpen = true

	return nil
}

func eventLoop() {
	if eventLoopRunning {
		return
	}

	go func() {
		eventLoopRunning = true
		onOff := true
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if onOff {
					stopNotes()
				} else {
					nextEvent()
				}
				onOff = !onOff
			}
		}
		eventLoopRunning = false
	}()
}

func playSong(songIndex int) error {
	song := &Songs[songIndex]

	// add a tempo to the first event if necessary
	if Tempo == 0 {
		if len(song.Events) == 0 {
			return errors.New("Song has no events")
		} else {
			if song.Events[0].Tempo == 0 {
				song.Events[0].Tempo = DefaultTempo
				println("setting default tempo")
			}
		}
	}

	ActiveSong = songIndex
	ActiveSongEvent = -1
	ActivePatterns = make([]int, 0)
	nextEvent()

	return nil
}

func setTempo(tempo uint16) {
	Tempo = tempo
	tickDuration := time.Minute / time.Duration(int(tempo)*EventsPerBeat*2)
	ticker = time.NewTicker(tickDuration)
	eventLoop()
}

func nextEvent() {
	var newActivePatterns []int
	for _, i := range ActivePatterns {
		if Patterns[i].ActiveEvent < len(Patterns[i].Events) {
			for j := range Patterns[i].Events[Patterns[i].ActiveEvent].Instruments {
				Trigger(Patterns[i].Events[Patterns[i].ActiveEvent].Instruments[j])
			}
		}
		Patterns[i].ActiveEvent++
		if Patterns[i].ActiveEvent < len(Patterns[i].Events) {
			newActivePatterns = append(newActivePatterns, i)
		}
	}

	ActivePatterns = newActivePatterns

	if len(ActivePatterns) > 0 {
		return
	}

	ActiveSongEvent++
	if ActiveSongEvent >= len(Songs[ActiveSong].Events) {
		stopped <- true
		return
	}

	if Songs[ActiveSong].Events[ActiveSongEvent].Repeat {
		ActiveSongEvent = -1
		nextEvent()
		return
	}

	if Songs[ActiveSong].Events[ActiveSongEvent].Chain {
		ActiveSong = Songs[ActiveSong].Events[ActiveSongEvent].ChainIndex
		ActiveSongEvent = -1
		nextEvent()
		return
	}

	if Songs[ActiveSong].Events[ActiveSongEvent].Tempo != 0 {
		setTempo(Songs[ActiveSong].Events[ActiveSongEvent].Tempo)
	}

	for i := range Songs[ActiveSong].Events[ActiveSongEvent].Patterns {
		activatePattern(Songs[ActiveSong].Events[ActiveSongEvent].Patterns[i])
	}
}

func activatePattern(index int) {
	Patterns[index].ActiveEvent = 0
	ActivePatterns = append(ActivePatterns, index)
}

func stopNotes() {
	for i := range Instruments {
		if Instruments[i].Sounding {
			send(Instruments[i].Off)
			Instruments[i].Sounding = false
		}
	}
}
