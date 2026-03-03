package spectogram

import (
	"math"
)

func GenSpectogram(samples []float64) (int, error) {

	windowSize := 2048
	hopSize := 512

	frameCount := 0
	// GenerateHann()
	hann := GenerateHann(windowSize)
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

		frameCount++
	}

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
