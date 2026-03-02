package audio

import (
	"fmt"
	"os"

	"github.com/go-audio/wav"
)

// []float64, int, error
// []float64z, int, error
func LoadWav(path string) ([]float64, int, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, 0, err
		// fmt.Println(err)
	}

	decoder := wav.NewDecoder(file)
	if !decoder.IsValidFile() {
		return nil, 0, fmt.Errorf("Invalid wav file ")
		// fmt.Println(err)

	}

	buff, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, 0, err
		// fmt.Println(err)

	}

	_ = buff

	fmt.Println(decoder)

	sampleRate := int(decoder.SampleRate)
	numChannel := buff.Format.NumChannels
	IntSamples := buff.AsIntBuffer().Data

	var samples []float64

	if numChannel == 1 {
		for _, s := range IntSamples {
			samples = append(samples, float64(s)/32768.0)
		}
	}

	if numChannel == 2 {
		for i := 0; i < len(IntSamples); i += 2 {
			left := i
			right := i + 2
			mono := (right + left) / 2
			samples = append(samples, float64(mono)/32768.0)

		}
	} else {
		return nil, 0, fmt.Errorf("unsupported channel count: %d", numChannel)

	}

	return samples, sampleRate, nil

}
