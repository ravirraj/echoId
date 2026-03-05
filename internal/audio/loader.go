package audio

import (
	"fmt"
	"os"

	"github.com/go-audio/wav"
	"github.com/hajimehoshi/go-mp3"
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
	}

	if numChannel == 2 {
		for i := 0; i < len(IntSamples); i += 2 {
			left := i
			right := i + 2
			mono := (right + left) / 2
			samples = append(samples, float64(mono)/32768.0)

		}
	} else {
		return nil, fmt.Errorf("unsupported channel count: %d", numChannel)

	}

	return samples, nil

}

func LoadMp3(path string) ([]float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, decoder.Length())

	_, err = decoder.Read(buf)
	if err != nil {
		return nil, err
	}

	samples := []float64{}

	for i := 0; i < len(buf); i += 2 {
		sample := int16(buf[i]) | int16(buf[i+1])<<8

		samples = append(samples, float64(sample)/32768.0)
	}

	return samples, nil
}
