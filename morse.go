package main

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
	"math"
	"os"
	"strings"
	"time"
)

var morseAlphabet = map[string]string{
	"a": ".-",
	"b": "-...",
	"c": "-.-.",
	"d": "-..",
	"e": ".",
	"f": "..-.",
	"g": "--.",
	"h": "....",
	"i": "..",
	"j": ".---",
	"k": "-.-",
	"l": ".-..",
	"m": "--",
	"n": "-.",
	"o": "---",
	"p": ".--.",
	"q": "--.-",
	"r": ".-.",
	"s": "...",
	"t": "-",
	"u": "..-",
	"v": "...-",
	"w": ".--",
	"x": "-..-",
	"y": "-.--",
	"z": "--.",
	"1": ".----",
	"2": "..---",
	"3": "...--",
	"4": "....-",
	"5": ".....",
	"6": "-....",
	"7": "--...",
	"8": "---..",
	"9": "----.",
	"0": "-----",
	" ": "   ",
}

type stereoSine struct {
	*portaudio.Stream
	stepL, phaseL float64
	stepR, phaseR float64
}

func newStereoSine(freqL, freqR, samplerate float64) *stereoSine {
	s := &stereoSine{nil, freqL / samplerate, 0, freqR / samplerate, 0}
	var err error
	s.Stream, err = portaudio.OpenDefaultStream(0, 2, samplerate, 0, s.processAudio)
	chk(err)
	return s
}

func (g *stereoSine) processAudio(out [][]float32) {
	for i := range out[0] {
		out[0][i] = float32(math.Sin(2 * math.Pi * g.phaseL))
		_, g.phaseL = math.Modf(g.phaseL + g.stepL)
		out[1][i] = float32(math.Sin(2 * math.Pi * g.phaseR))
		_, g.phaseR = math.Modf(g.phaseR + g.stepR)
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func toMorse(str string, s *stereoSine) {
	str = strings.ToLower(str)
	var toMorse string
	for i, r := range str {
		c := string(r)
		toMorse += morseAlphabet[c]
		if i < len(str) && c != " " {
			toMorse += " "
		}
	}

	fmt.Println(toMorse)

	for _, r := range toMorse {
		c := string(r)
		if c == "." {
			chk(s.Start())
			time.Sleep(75 * time.Millisecond)
			chk(s.Stop())
			time.Sleep(50 * time.Millisecond)
		} else if c == "-" {
			chk(s.Start())
			time.Sleep(175 * time.Millisecond)
			chk(s.Stop())
			time.Sleep(75 * time.Millisecond)
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}

}

func main() {
	args := os.Args[1:]
	portaudio.Initialize()
	defer portaudio.Terminate()
	toMorse(args[0], newStereoSine(3000, 5000, 44100))

}
