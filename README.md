# drumscript

A simple MIDI drum machine scripting language.

drumscript triggers MIDI events based on scripted drum patterns; you need a MIDI instrument, sound card or software synth (such as TiMidity++) to hear the sounds.

## Building and Installation

You need Go >= 1.16 and the ability to build CGo packages. You will also need the headers for libportmidi - on Ubuntu or similar you can install the package "libportmidi-dev".

Then as usual:

    go mod tidy
    go build

You can optionally copy the drumscript executable to be anywhere on your $PATH.

The language is documented in the comments of example.drumscript, which is designed to be run with TiMidity++. Use the following commands:

    timidity -iA -B2,8 &
    ./drumscript example.drumscript

