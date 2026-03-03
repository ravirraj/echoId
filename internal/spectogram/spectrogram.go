package spectogram

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/dsp/fourier"
)

func GenSpectogram(samples []float64) (int, error) {

	windowSize := 2048
	hopSize := 512
	// var result []float64
	spectogram := [][]float64{}

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

		for k := 0; k < n; k++ {
			realPart := real(output[k])
			imgPart := imag(output[k])

			mang := realPart*realPart + imgPart*imgPart
			mangitudes[k] = mang
			// result = append(result, mang)

		}
		spectogram = append(spectogram, mangitudes)

		frameCount++
		// fmt.Println(real(output[0]))

	}
	fmt.Println(spectogram[0][100])
	fmt.Println(len(spectogram[0]))


	// fmt.Println(result[0])
	// fmt.Println(real(output[0]))

	return frameCount, nil

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
