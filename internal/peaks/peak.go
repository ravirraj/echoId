package peak

import "sort"

// type Peak struct {
// 	TimeIndex int
// 	FreqIndex int
// 	Magnitude float64
// }

// type Bin struct {
// 	Index int
// 	Value float64
// }

func DetectPeaks(spectrogram [][]float64) []Peak {

	peaks := []Peak{}

	// Calculate global magnitude threshold (use median-based approach)
	var allMagnitudes []float64
	for _, magnitudes := range spectrogram {
		for _, mag := range magnitudes {
			allMagnitudes = append(allMagnitudes, mag)
		}
	}
	sort.Float64s(allMagnitudes)
	threshold := allMagnitudes[len(allMagnitudes)*75/100] // 75th percentile

	for frameIndex, magnitudes := range spectrogram {
		// Find local maxima with neighborhood check
		for freqIndex := 3; freqIndex < len(magnitudes)-3; freqIndex++ {
			mag := magnitudes[freqIndex]

			// Must be above threshold
			if mag < threshold {
				continue
			}

			// Must be local maximum in frequency domain (check ±3 bins)
			isLocalMax := true
			for offset := -3; offset <= 3; offset++ {
				if offset == 0 {
					continue
				}
				if magnitudes[freqIndex+offset] >= mag {
					isLocalMax = false
					break
				}
			}

			if isLocalMax {
				peaks = append(peaks, Peak{
					TimeIndex: frameIndex,
					FreqIndex: freqIndex,
					Magnitude: mag,
				})
			}
		}
	}

	return peaks
}
