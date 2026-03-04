package main

import (
	"fmt"
	"os"

	"github.com/ravirraj/echoid/internal/audio"
	"github.com/ravirraj/echoid/internal/fingerprint"
	"github.com/ravirraj/echoid/internal/peak"
	spectrogram "github.com/ravirraj/echoid/internal/spectogram"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: echoid <file.wav>")
		return
	}

	filePath := os.Args[1]

	samples, _, err := audio.LoadWav(filePath)

	if err != nil {
		fmt.Println(err)
		return
	}
	spectrogram := spectrogram.GenerateSpectrogram(samples)

	peaks := peak.DetectPeaks(spectrogram)

	fingerprints := fingerprint.GenerateFingerprints(peaks)
	fmt.Println(fingerprints)

}
