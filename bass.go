package bass

import (
	"fmt"
	"math"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

// SawtoothOscillator generates a sawtooth waveform at a specific frequency
func SawtoothOscillator(freq float64, length int, sampleRate int) []float64 {
	osc := make([]float64, length)
	for i := range osc {
		osc[i] = 2 * (float64(i)/float64(sampleRate)*freq - math.Floor(0.5+float64(i)/float64(sampleRate)*freq))
	}
	return osc
}

// ApplyEnvelope applies an ADSR envelope to the waveform
func ApplyEnvelope(samples []float64, attack, decay, sustain, release float64, sampleRate int) []float64 {
	adsr := make([]float64, len(samples))
	for i := range samples {
		t := float64(i) / float64(sampleRate)
		if t < attack {
			adsr[i] = samples[i] * (t / attack)
		} else if t < attack+decay {
			adsr[i] = samples[i] * (1 - (t-attack)/decay*(1-sustain))
		} else if t < float64(len(samples))/float64(sampleRate)-release {
			adsr[i] = samples[i] * sustain
		} else {
			adsr[i] = samples[i] * (1 - (t-(float64(len(samples))/float64(sampleRate)-release))/release*sustain)
		}
	}
	return adsr
}

// LowPassFilter applies a basic low-pass filter to the samples
func LowPassFilter(samples []float64, cutoff float64, sampleRate int) []float64 {
	filtered := make([]float64, len(samples))
	alpha := 2 * math.Pi * cutoff / float64(sampleRate)
	prev := 0.0
	for i, sample := range samples {
		filtered[i] = prev + alpha*(sample-prev)
		prev = filtered[i]
	}
	return filtered
}

// Drive applies a simple drive effect by scaling and clipping
func Drive(samples []float64, gain float64) []float64 {
	driven := make([]float64, len(samples))
	for i, sample := range samples {
		driven[i] = sample * gain
		if driven[i] > 1 {
			driven[i] = 1
		} else if driven[i] < -1 {
			driven[i] = -1
		}
	}
	return driven
}

// Limiter ensures the signal doesn't exceed [-1, 1] range
func Limiter(samples []float64) []float64 {
	limited := make([]float64, len(samples))
	for i, sample := range samples {
		if sample > 1 {
			limited[i] = 1
		} else if sample < -1 {
			limited[i] = -1
		} else {
			limited[i] = sample
		}
	}
	return limited
}

// DetunedOscillators generates multiple detuned sawtooth oscillators and combines them
func DetunedOscillators(freq float64, detune []float64, length int, sampleRate int) []float64 {
	numOsc := len(detune)
	combined := make([]float64, length)
	for _, d := range detune {
		osc := SawtoothOscillator(freq*(1+d), length, sampleRate)
		for i := range combined {
			combined[i] += osc[i] / float64(numOsc) // Average to avoid high amplitudes
		}
	}
	return combined
}

// SaveToWav saves the waveform to a wav file
func SaveToWav(filename string, samples []float64, sampleRate int) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating wav file: %v", err)
	}
	defer outFile.Close()

	enc := wav.NewEncoder(outFile, sampleRate, 16, 1, 1)

	buf := &audio.IntBuffer{
		Format: &audio.Format{SampleRate: sampleRate, NumChannels: 1},
		Data:   make([]int, len(samples)),
	}

	for i, sample := range samples {
		buf.Data[i] = int(sample * math.MaxInt16) // Convert to 16-bit PCM
	}

	if err := enc.Write(buf); err != nil {
		return fmt.Errorf("error writing wav file: %v", err)
	}

	return enc.Close()
}
