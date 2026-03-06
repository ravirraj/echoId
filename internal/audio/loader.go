package audio

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		return nil, err
	}
	// fmt.Println(decoder.Length())
	buf := make([]byte, 4096)
	// fmt.Println("raw bytes:", buf[:20])

	samples := []float64{}

	for {

		n, err := decoder.Read(buf)
		// fmt.Println("raw bytes:", buf[:20])

		if n == 0 {
			break
		}

		// for i := 0; i+3 < n; i += 4 {

		// 	left := int16(binary.LittleEndian.Uint16(buf[i : i+2]))
		// 	right := int16(binary.LittleEndian.Uint16(buf[i+2 : i+4]))
		// 	mono := (float64(left) + float64(right)) / 2.0

		// 	samples = append(samples, mono/32768.0)
		// }
		for i := 0; i+1 < n; i += 2 {

			sample := int16(binary.LittleEndian.Uint16(buf[i : i+2]))

			samples = append(samples, float64(sample)/32768.0)
		}

		if err != nil {
			break
		}
	}
	// fmt.Println(decoder.SampleRate())
	// fmt.Println(len(samples))
	// fmt.Println("raw bytes:", buf[:20])

	// fmt.Println(samples[:4])
	fmt.Println(len(samples))
	fmt.Printf("%.20f\n", samples[0])
	fmt.Println(samples[0] * 1000000000000000000000)
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
