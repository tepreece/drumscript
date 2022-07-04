package main

import (
	"errors"
	"fmt"
	"unicode"
)

type PatternEvent struct {
	Instruments []int
}

type PatternBeat struct {
	Events []PatternEvent
}

type Pattern struct {
	Name          string
	Ready, Active bool
	ActiveEvent   int
	Beats         []PatternBeat
	Events        []PatternEvent
}

var (
	Patterns      []Pattern
	EventsPerBeat int
	NoOpEvent     PatternEvent
)

func startParsePattern(fields []string, state *parseState) error {
	if state.State != BaseState {
		return errors.New("Can't start a new pattern")
	}

	if len(fields) != 2 {
		return errors.New("Invalid pattern name")
	}

	name := fields[1]
	for _, r := range name {
		if !unicode.IsLetter(r) {
			return errors.New("Invalid pattern name")
		}
	}

	state.State = PatternState
	state.Pattern.Name = name

	return nil
}

func endParsePattern(state *parseState) {
	Patterns = append(Patterns, *state.Pattern)
	state.State = BaseState
	state.Pattern = new(Pattern)
}

func parsePatternLine(fields []string, state *parseState) error {
	var b PatternBeat
	for _, f := range fields {
		var e PatternEvent

		if f == "." {
			b.Events = append(b.Events, NoOpEvent)
			continue
		}

		for _, r := range f {
			instrument, err := getInstrumentIndex(r)
			if err != nil {
				return err
			}

			e.Instruments = append(e.Instruments, instrument)
		}

		b.Events = append(b.Events, e)
	}
	state.Pattern.Beats = append(state.Pattern.Beats, b)
	return nil
}

func getPatternIndex(name string) (int, error) {
	for i := range Patterns {
		if Patterns[i].Name == name {
			return i, nil
		}
	}

	return -1, errors.New(fmt.Sprintf("Pattern %s not defined", name))
}

func convertPatterns() {
	// for every defined beat, find the number of events in that beat
	eventsPerBeat := make(map[int]bool)
	for i := range Patterns {
		for j := range Patterns[i].Beats {
			eventsThisBeat := len(Patterns[i].Beats[j].Events)
			eventsPerBeat[eventsThisBeat] = true
		}
	}

	// find the lowest common multiple of all the events per beat
	var eventCounts []int
	for i := range eventsPerBeat {
		if eventsPerBeat[i] == true {
			eventCounts = append(eventCounts, i)
		}
	}

	EventsPerBeat = LeastCommonMultiple(eventCounts)

	// now add the relevant number of NoOpEvents to every beat
	for i := range Patterns {
		for j := range Patterns[i].Beats {
			noOpsPerEvent := (EventsPerBeat / len(Patterns[i].Beats[j].Events)) - 1
			for _, e := range Patterns[i].Beats[j].Events {
				Patterns[i].Events = append(Patterns[i].Events, e)
				for k := 0; k < noOpsPerEvent; k++ {
					Patterns[i].Events = append(Patterns[i].Events, NoOpEvent)
				}
			}
		}

		Patterns[i].Ready = true
	}
}
