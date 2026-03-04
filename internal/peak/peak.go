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

	for frameIndex, magnitudes := range spectrogram {

		n := len(magnitudes)

		bins := make([]Bin, n)

		for i := 0; i < n; i++ {
			bins[i] = Bin{
				Index: i,
				Value: magnitudes[i],
			}
		}

		sort.Slice(bins, func(i, j int) bool {
			return bins[i].Value > bins[j].Value
		})

		for p := 0; p < 5; p++ {
			peaks = append(peaks, Peak{
				TimeIndex: frameIndex,
				FreqIndex: bins[p].Index,
				Magnitude: bins[p].Value,
			})
		}
	}

	return peaks
}
