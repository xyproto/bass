// cmd/create/main.go
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/xyproto/bass"
)

var (
	version     = "1.0.0"
	sampleRate  int
	duration    time.Duration
	baseFreq    float64
	showVersion bool
	showHelp    bool
)

func init() {
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showHelp, "help", false, "Show help information")
	flag.IntVar(&sampleRate, "samplerate", 44100, "Sample rate (in Hz)")
	flag.DurationVar(&duration, "duration", 10*time.Second, "Duration of the audio (e.g., 10s, 5m)")
	flag.Float64Var(&baseFreq, "freq", 55.0, "Base frequency for the bass sound (in Hz)")
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Printf("Bass Synth Generator, version %s\n", version)
		os.Exit(0)
	}

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	length := sampleRate * int(duration.Seconds())
	detune := []float64{-0.01, -0.005, 0.0, 0.005, 0.01}

	bassOscillators := bass.DetunedOscillators(baseFreq, detune, length, sampleRate)
	env := bass.ApplyEnvelope(bassOscillators, 0.1, 0.4, 0.6, 0.7, sampleRate)
	filtered := bass.LowPassFilter(env, 200, sampleRate)
	driven := bass.Drive(filtered, 1.2)
	limited := bass.Limiter(driven)

	if err := bass.SaveToWav("bass_output.wav", limited, sampleRate); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Successfully generated 'bass_output.wav'")
}
