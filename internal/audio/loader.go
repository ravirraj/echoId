package audio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-audio/wav"
)

const (
	normalizeSampleRate = 44100
	normalizeChannels   = 1
	normalizeBitDepth   = 16
)

func LoadWav(path string) ([]float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open wav: %w", err)
	}
	defer file.Close()

	decoder := wav.NewDecoder(file)
	if !decoder.IsValidFile() {
		return nil, fmt.Errorf("invalid wav file")
	}

	buff, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, fmt.Errorf("decode wav: %w", err)
	}

	numChannel := buff.Format.NumChannels
	IntSamples := buff.AsIntBuffer().Data

	var samples []float64

	if numChannel == 1 {
		for _, s := range IntSamples {
			samples = append(samples, float64(s)/32768.0)
		}
	} else if numChannel == 2 {
		for i := 0; i < len(IntSamples); i += 2 {
			mono := (IntSamples[i] + IntSamples[i+1]) / 2
			samples = append(samples, float64(mono)/32768.0)
		}
	} else {
		return nil, fmt.Errorf("unsupported channel count: %d", numChannel)
	}

	return samples, nil
}

func LoadAudio(path string) ([]float64, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	tempDir := filepath.Join(dir, "temp", "normalized")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}

	outputPath := filepath.Join(tempDir, "normalized.wav")

	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", path,
		"-ac", "1",
		"-ar", "44100",
		"-sample_fmt", "s16",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg: %w\n%s", err, string(output))
	}

	samples, err := LoadWav(outputPath)
	if err != nil {
		return nil, err
	}

	if err := os.Remove(outputPath); err != nil {
		return nil, fmt.Errorf("cleanup temp: %w", err)
	}

	return samples, nil
}
