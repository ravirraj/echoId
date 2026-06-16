package spectrogram

import (
	"math"

	"gonum.org/v1/gonum/dsp/fourier"
)

func GenerateSpectrogram(samples []float64) [][]float64 {

	windowSize := 2048
	hopSize := 512

	spectrogram := [][]float64{}

	hann := GenerateHann(windowSize)
	fft := fourier.NewFFT(windowSize)

	for i := 0; i+windowSize <= len(samples); i += hopSize {

		rawFrame := samples[i : i+windowSize]

		frameCopy := make([]float64, windowSize)
		copy(frameCopy, rawFrame)

		for j := 0; j < windowSize; j++ {
			frameCopy[j] *= hann[j]
		}

		output := fft.Coefficients(nil, frameCopy)

		// Only use first half of FFT (remove symmetric part) and limit to useful frequency range
		// Use up to bin 512 (covers 0-12kHz at 48kHz sample rate, which is the most important range)
		maxBin := 512
		if len(output)/2 < maxBin {
			maxBin = len(output) / 2
		}

		magnitudes := make([]float64, maxBin)

		for k := 0; k < maxBin; k++ {
			realPart := real(output[k])
			imagPart := imag(output[k])
			magnitudes[k] = math.Sqrt(realPart*realPart + imagPart*imagPart)
		}

		spectrogram = append(spectrogram, magnitudes)
	}

	return spectrogram
}

func GenerateHann(N int) []float64 {

	result := make([]float64, N)

	for n := 0; n < N; n++ {
		progress := float64(n) / float64(N-1)
		angle := 2 * math.Pi * progress
		result[n] = 0.5 * (1 - math.Cos(angle))
	}

	return result
}
