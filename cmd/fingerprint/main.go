package main

import (
	"fmt"
	"os"

	"github.com/ravirraj/echoid/internal/audio"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: echoid <file.wav>")
		return
	}

	filePath := os.Args[1]

	samples, sampleRate, err := audio.LoadWav(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	duration := float64(len(samples)) / float64(sampleRate)

	fmt.Println("Sample Rate:", sampleRate)
	fmt.Println("Total Samples (mono):", len(samples))
	fmt.Printf("Duration: %.2f seconds\n", duration)

	fmt.Println("First 10 samples:")
	for i := 0; i < 10 && i < len(samples); i++ {
		fmt.Printf("%.5f\n", samples[i])
		// fmt.Println(samples[i])
	}

}
