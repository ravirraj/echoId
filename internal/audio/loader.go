package audio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-audio/wav"
	// "github.com/hajimehoshi/go-mp3"
)

// []float64, int, error
// []float64z, int, error
func LoadWav(path string) ([]float64, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
		// fmt.Println(err)
	}

	decoder := wav.NewDecoder(file)
	if !decoder.IsValidFile() {
		return nil, fmt.Errorf("Invalid wav file")
		// fmt.Println(err)

	}

	buff, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, err
		// fmt.Println(err)

	}

	fmt.Println(decoder)

	// sampleRate := int(decoder.SampleRate)
	numChannel := buff.Format.NumChannels
	IntSamples := buff.AsIntBuffer().Data

	var samples []float64

	if numChannel == 1 {
		for _, s := range IntSamples {
			samples = append(samples, float64(s)/32768.0)
		}
	} else if numChannel == 2 {
		for i := 0; i < len(IntSamples); i += 2 {
			left := IntSamples[i]
			right := IntSamples[i+1]
			mono := (left + right) / 2
			samples = append(samples, float64(mono)/32768.0)

		}
	} else {
		return nil, fmt.Errorf("unsupported channel count: %d", numChannel)

	}

	return samples, nil

}

func LoadMp3(path string) ([]float64, error) {
	cmd := exec.Command(
		"ffmpeg",
		"-i",
		path,
		"converted.wav",
	)

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	mainDir := filepath.Join(dir, "temp", "converted")
	cmd.Dir = mainDir

	_, err = cmd.Output()

	if err != nil {
		return nil, err
	}

	inputPath := filepath.Join(mainDir, "converted.wav")

	samples, err := LoadWav(inputPath)
	if err != nil {
		return nil, err
	}

	err = os.Remove(inputPath)
	if err != nil {
		return nil, err
	}

	return samples, nil
}
func LoadAudio(path string) ([]float64, error) {
	ext := strings.ToLower(filepath.Ext(path))
	fmt.Println(ext)

	switch ext {
	case ".wav":
		return LoadWav(path)
	case ".mp3":
		return LoadMp3(path)
	default:
		return nil, fmt.Errorf("unsupported file , we support .wav,.mp3 only")
	}

}
