package main

import (
	"fmt"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer midi.CloseDriver()

	fname := "test.drumscript"
	err := parseScript(fname)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !PortOpen {
		fmt.Println("No MIDI port opened")
		return
	}

	songIndex, err := getSongIndex("")
	if err != nil {
		fmt.Println("No default song specified")
		return
	}

	if err := playSong(songIndex); err != nil {
		fmt.Println(err.Error())
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	select {
	case <-sc:
	case <-stopped:
	}

	stopNotes()
	done <- true
}
