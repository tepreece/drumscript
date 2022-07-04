package main

import (
	"errors"
	"fmt"
	"gitlab.com/gomidi/midi/v2"
	"strconv"
	"unicode"
)

const (
	PercussionChannel uint8  = 9 // MIDI channel 10, but the midi library is 0-indexed
	DefaultVelocity   uint64 = 127
)

type Instrument struct {
	Letter   rune
	On, Off  midi.Message
	Sounding bool
}

var (
	Instruments []*Instrument
)

func createInstrument(fields []string) error {
	instrument := new(Instrument)

	if len(fields) < 3 {
		return errors.New("Incomplete instrument specification")
	}

	// check for a valid instrument letter
	if len(fields[1]) != 1 {
		return errors.New("Invalid instrument letter")
	}

	letter := ' '
	for _, c := range fields[1] {
		letter = c
		break
	}

	if !unicode.IsLetter(letter) {
		return errors.New("Invalid instrument letter")
	}

	for i := range Instruments {
		if Instruments[i].Letter == letter {
			return errors.New("Repeated instrument letter")
		}
	}

	instrument.Letter = letter

	// check for a valid instrument number
	number, err := strconv.ParseUint(fields[2], 10, 8)
	if err != nil {
		return errors.New("Invalid instrument number")
	}

	// check for a valid velocity
	velocity := DefaultVelocity
	if len(fields) >= 4 {
		velocity, err = strconv.ParseUint(fields[3], 10, 8)
		if err != nil {
			return errors.New("Invalid velocity")
		}
	}

	instrument.On = midi.NoteOn(PercussionChannel, uint8(number), uint8(velocity))
	instrument.Off = midi.NoteOff(PercussionChannel, uint8(number))

	Instruments = append(Instruments, instrument)

	return nil
}

func getInstrumentIndex(letter rune) (int, error) {
	for i := range Instruments {
		if Instruments[i].Letter == letter {
			return i, nil
		}
	}

	return -1, errors.New(fmt.Sprintf("Instrument %s not defined", string(letter)))
}

func Trigger(index int) {
	if !Instruments[index].Sounding {
		send(Instruments[index].On)
		Instruments[index].Sounding = true
	}
}
