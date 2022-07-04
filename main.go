package main

import (
	"flag"
	"fmt"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer midi.CloseDriver()

	portPtr := flag.String("p", "", "Override the port set in the script")
	listPtr := flag.Bool("l", false, "List the available MIDI ports")
	flag.Parse()

	if *listPtr {
		listPorts()
		return
	}

	OverridePortName = *portPtr
	if OverridePortName != "" {
		setPort([]string{})
	}

	fname := ""
	songname := ""

	args := flag.Args()
	switch len(args) {
	case 0:
		fmt.Printf("Usage: %s [-p port] filename [song]\n", os.Args[0])
	case 1:
		fname = args[0]
	case 2:
		fname = args[0]
		songname = args[1]
	}

	err := parseScript(fname)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !PortOpen {
		fmt.Println("No MIDI port opened")
		return
	}

	songIndex, err := getSongIndex(songname)
	if err != nil {
		if songname == "" {
			fmt.Println("No default song found")
		} else {
			fmt.Printf("Song %s not found\n", songname)
		}
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
