package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type songEvent struct {
	Repeat, Chain bool
	Tempo         uint16
	ChainName     string
	ChainIndex    int
	Patterns      []int
}

type Song struct {
	Name   string
	Events []songEvent
}

var Songs []Song

func startParseSong(fields []string, state *parseState) error {
	if state.State != BaseState {
		return errors.New("Can't start a new song")
	}

	name := ""
	switch len(fields) {
	case 1:
		// no name => default song - do nothing
	case 2:
		// song name specified
		name = fields[1]
		for _, r := range name {
			if !unicode.IsLetter(r) {
				return errors.New("Invalid song name")
			}
			break
		}
	default:
		return errors.New("Invalid song name")
	}

	// check to see if we already have a song with this name
	for i := range Songs {
		if Songs[i].Name == name {
			return errors.New("Repeated song name")
		}
	}

	state.State = SongState
	state.Song.Name = name

	return nil
}

func endParseSong(state *parseState) {
	Songs = append(Songs, *state.Song)
	state.State = BaseState
	state.Song = new(Song)
}

func parseSongTempo(fields []string, state *parseState) error {
	if len(fields) != 2 {
		return errors.New("Invalid tempo")
	}

	tempo, err := strconv.ParseUint(fields[1], 10, 16)
	if err != nil {
		return errors.New("Invalid tempo")
	}

	if tempo == 0 {
		return errors.New("Invalid tempo")
	}

	state.Tempo = uint16(tempo)
	return nil
}

func parseSongRepeat(fields []string, state *parseState) error {
	if len(fields) != 1 {
		return errors.New("Invalid repeat command")
	}

	var e songEvent
	if state.Tempo != 0 {
		e.Tempo = state.Tempo
		state.Tempo = 0
	}

	e.Repeat = true
	state.Song.Events = append(state.Song.Events, e)

	endParseSong(state)

	return nil
}

func parseSongChain(fields []string, state *parseState) error {
	if len(fields) != 2 {
		return errors.New("Invalid chain command")
	}

	name := fields[1]
	for _, r := range name {
		if !unicode.IsLetter(r) {
			return errors.New("Invalid song name in chain")
		}
	}

	var e songEvent
	if state.Tempo != 0 {
		e.Tempo = state.Tempo
		state.Tempo = 0
	}

	e.Chain = true
	e.ChainName = name
	state.Song.Events = append(state.Song.Events, e)
	endParseSong(state)

	return nil
}

func parseSongLine(fields []string, state *parseState) error {
	var (
		e   songEvent
		err error
	)

	if state.Tempo != 0 {
		e.Tempo = state.Tempo
		state.Tempo = 0
	}

	first := true
	count := 1
	for _, f := range fields {
		// for the first field, it count be a count
		if first {
			first = false
			count, err = strconv.Atoi(f)
			if err == nil {
				if count < 1 {
					return errors.New("Invalid count")
				}

				// we found a count - continue to the next field
				continue
			} else {
				count = 1
			}
		}

		// for all the remaining fields, it must be a pattern name
		patternIndex, err := getPatternIndex(f)
		if err != nil {
			return errors.New(fmt.Sprintf("Undefined pattern %s", f))
		}
		e.Patterns = append(e.Patterns, patternIndex)
	}

	if len(e.Patterns) == 0 {
		return errors.New("No patterns")
	}

	for i := 0; i < count; i++ {
		state.Song.Events = append(state.Song.Events, e)
	}

	return nil
}

func getSongIndex(name string) (int, error) {
	for i := range Songs {
		if Songs[i].Name == name {
			return i, nil
		}
	}

	return -1, errors.New(fmt.Sprintf("Pattern %s not defined", name))
}

func convertSongs() error {
	var err error
	for i := range Songs {
		for j := range Songs[i].Events {
			if Songs[i].Events[j].Chain {
				name := Songs[i].Events[j].ChainName
				Songs[i].Events[j].ChainIndex, err = getSongIndex(name)
				if err != nil {
					errors.New(fmt.Sprintf("Undefined song to chain %s", name))
				}
			}
		}
	}

	return nil
}
