package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type parseStateState int

const (
	BaseState parseStateState = iota
	PatternState
	SongState
)

type parseState struct {
	State   parseStateState
	Tempo   uint16
	Pattern *Pattern
	Song    *Song
}

func parseScript(fname string) error {
	state := new(parseState)
	state.Pattern = new(Pattern)
	state.Song = new(Song)

	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++

		// remove comments and leading/trailing whitespace
		line := strings.TrimSpace(strings.SplitN(scanner.Text(), "#", 2)[0])
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}

		err = parseFields(fields, state)
		if err != nil {
			return errors.New(fmt.Sprintf("%s on line %d", err.Error(), lineNumber))
		}
	}

	switch state.State {
	case PatternState:
		endParsePattern(state)
	case SongState:
		endParseSong(state)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	convertPatterns()
	if err := convertSongs(); err != nil {
		return err
	}

	return nil
}

func parseFields(fields []string, state *parseState) error {
	switch state.State {
	case BaseState:
		switch fields[0] {
		case "port":
			return setPort(fields)
		case "instrument":
			return createInstrument(fields)
		case "pattern":
			return startParsePattern(fields, state)
		case "song":
			return startParseSong(fields, state)
		default:
			return errors.New("Invalid command")
		}
	case PatternState:
		switch fields[0] {
		case "end":
			endParsePattern(state)
			return nil
		case "pattern":
			endParsePattern(state)
			return startParsePattern(fields, state)
		case "song":
			endParsePattern(state)
			return startParseSong(fields, state)
		default:
			return parsePatternLine(fields, state)
		}
	case SongState:
		switch fields[0] {
		case "end":
			endParseSong(state)
			return nil
		case "pattern":
			endParseSong(state)
			return startParsePattern(fields, state)
		case "song":
			endParseSong(state)
			return startParseSong(fields, state)
		case "tempo":
			return parseSongTempo(fields, state)
		case "repeat":
			return parseSongRepeat(fields, state)
		case "chain":
			return parseSongChain(fields, state)
		default:
			return parseSongLine(fields, state)
		}
	default:
		return errors.New("Invalid parse state")
	}
	return nil
}
