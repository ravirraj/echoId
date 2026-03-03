package spectogram

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/dsp/fourier"
)

type Peak struct {
	TimeIndex int
	FreqIndex int
	Magnitude float64
}
type Bin struct {
	Index int
	Value float64
}

type Fingerprint struct {
	Freq1      int
	Freq2      int
	DeltaTime  int
	AnchorTime int
}

func GenSpectogram(samples []float64) ([][]float64, error) {

	windowSize := 2048
	hopSize := 512
	// var result []float64
	spectogram := [][]float64{}

	peaks := []Peak{}
	fingerprints := []Fingerprint{}
	fanOut := 5
	maxDelta := 50 // frames

	frameCount := 0
	// GenerateHann()
	hann := GenerateHann(windowSize)
	fft := fourier.NewFFT(windowSize)
	for i := 0; i+windowSize <= len(samples); i += hopSize {
		// fmt.Println(samples[i])

		rawframe := samples[i : i+windowSize]
		frameCopy := make([]float64, windowSize)

		// frameCopy = append(frameCopy, rawframe...)
		copy(frameCopy, rawframe)

		for j := 0; j < windowSize; j++ {

			frameCopy[j] *= hann[j]

			// fmt.Println(frame)

		}

		output := fft.Coefficients(nil, frameCopy)

		n := len(output)
		mangitudes := make([]float64, n)
		bins := make([]Bin, n)
		for k := 0; k < n; k++ {
			realPart := real(output[k])
			imgPart := imag(output[k])

			mang := realPart*realPart + imgPart*imgPart
			mangitudes[k] = mang

			bins[k] = Bin{Index: k, Value: mangitudes[k]}
			// result = append(result, mang)

		}
		sort.Slice(bins, func(i, j int) bool { return bins[i].Value > bins[j].Value })

		for p := 0; p < 5; p++ {
			peaks = append(peaks, Peak{TimeIndex: frameCount, FreqIndex: bins[p].Index, Magnitude: bins[p].Value})
		}
		spectogram = append(spectogram, mangitudes)

		frameCount++
		// fmt.Println(real(output[0]))

	}

	for i := 0; i < len(peaks); i++ {
	innerLoop:
		for j := i + 1; j < len(peaks) && j < i+fanOut; j++ {
			deltatime := peaks[j].TimeIndex - peaks[i].TimeIndex
			if deltatime > maxDelta {
				break innerLoop
			}
			fingerprints = append(fingerprints, Fingerprint{
				Freq1:      peaks[i].FreqIndex,
				Freq2:      peaks[j].FreqIndex,
				DeltaTime:  deltatime,
				AnchorTime: peaks[i].TimeIndex,
			})
		}
	}
	fmt.Println(spectogram[0][100])
	fmt.Println(len(spectogram[0]))
	fmt.Println("total peaks ", len(peaks))
	fmt.Println("total fingerprints:", len(fingerprints))
	// fmt.Println(result[0])
	// fmt.Println(real(output[0]))

	return spectogram, nil

}

func GenerateHann(N int) []float64 {

	var result []float64

	// w[n] = 0.5 * (1 - cos(2πn/(N-1)))
	for n := 0; n < N; n++ {
		progress := float64(n) / (float64(N) - 1)
		angle := 2 * math.Pi * float64(progress)
		// math.Cos(angle)
		value := 0.5 * (1 - math.Cos(angle))
		result = append(result, value)

	}

	return result

}
